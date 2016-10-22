package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaModeRequest struct {
	action string
	Mode   proto.Mode
	reply  chan somaResult
}

type somaModeResult struct {
	ResultError error
	Mode        proto.Mode
}

func (a *somaModeResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Modes = append(r.Modes,
			somaModeResult{ResultError: err})
	}
}

func (a *somaModeResult) SomaAppendResult(r *somaResult) {
	r.Modes = append(r.Modes, *a)
}

/* Read Access
 */
type somaModeReadHandler struct {
	input     chan somaModeRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaModeReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ModeList: r.list_stmt,
		stmt.ModeShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`mode`, err, statement)
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

func (r *somaModeReadHandler) process(q *somaModeRequest) {
	var (
		mode string
		rows *sql.Rows
		err  error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: modes/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&mode)
			result.Append(err, &somaModeResult{
				Mode: proto.Mode{
					Mode: mode,
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: mode/show for %s", q.Mode.Mode)
		err = r.show_stmt.QueryRow(q.Mode.Mode).Scan(
			&mode,
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

		result.Append(err, &somaModeResult{
			Mode: proto.Mode{
				Mode: mode,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaModeWriteHandler struct {
	input    chan somaModeRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaModeWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ModeAdd: w.add_stmt,
		stmt.ModeDel: w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`mode`, err, statement)
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

func (w *somaModeWriteHandler) process(q *somaModeRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: modes/add for %s", q.Mode.Mode)
		res, err = w.add_stmt.Exec(
			q.Mode.Mode,
		)
	case "delete":
		w.reqLog.Printf("R: modes/del for %s", q.Mode.Mode)
		res, err = w.del_stmt.Exec(
			q.Mode.Mode,
		)
	default:
		w.reqLog.Printf("R: unimplemented modes/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaModeResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaModeResult{})
	default:
		result.Append(nil, &somaModeResult{
			Mode: q.Mode,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaModeReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaModeWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
