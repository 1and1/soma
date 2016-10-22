package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaLevelRequest struct {
	action string
	Level  proto.Level
	reply  chan somaResult
}

type somaLevelResult struct {
	ResultError error
	Level       proto.Level
}

func (a *somaLevelResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Levels = append(r.Levels, somaLevelResult{ResultError: err})
	}
}

func (a *somaLevelResult) SomaAppendResult(r *somaResult) {
	r.Levels = append(r.Levels, *a)
}

/* Read Access
 */
type somaLevelReadHandler struct {
	input     chan somaLevelRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaLevelReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.LevelList: r.list_stmt,
		stmt.LevelShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`level`, err, statement)
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-r.shutdown:
			break runloop
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaLevelReadHandler) process(q *somaLevelRequest) {
	var (
		level, short string
		numeric      uint16
		rows         *sql.Rows
		err          error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: levels/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&level, &short)
			result.Append(err, &somaLevelResult{
				Level: proto.Level{
					Name:      level,
					ShortName: short,
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: levels/show for %s", q.Level.Name)
		err = r.show_stmt.QueryRow(q.Level.Name).Scan(
			&level,
			&short,
			&numeric,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaLevelResult{
			Level: proto.Level{
				Name:      level,
				ShortName: short,
				Numeric:   numeric,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaLevelWriteHandler struct {
	input    chan somaLevelRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaLevelWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.LevelAdd: w.add_stmt,
		stmt.LevelDel: w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`level`, err, statement)
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-w.shutdown:
			break runloop
		case req := <-w.input:
			w.process(&req)
		}
	}
}

func (w *somaLevelWriteHandler) process(q *somaLevelRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: levels/add for %s", q.Level.Name)
		res, err = w.add_stmt.Exec(
			q.Level.Name,
			q.Level.ShortName,
			q.Level.Numeric,
		)
	case "delete":
		w.reqLog.Printf("R: levels/del for %s", q.Level.Name)
		res, err = w.del_stmt.Exec(
			q.Level.Name,
		)
	default:
		w.reqLog.Printf("R: unimplemented levels/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaLevelResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaLevelResult{})
	default:
		result.Append(nil, &somaLevelResult{
			Level: q.Level,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaLevelReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaLevelWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
