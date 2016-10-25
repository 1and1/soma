package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type somaOncallRequest struct {
	action string
	Oncall proto.Oncall
	reply  chan somaResult
}

type somaOncallResult struct {
	ResultError error
	Oncall      proto.Oncall
}

func (a *somaOncallResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Oncall = append(r.Oncall, somaOncallResult{ResultError: err})
	}
}

func (a *somaOncallResult) SomaAppendResult(r *somaResult) {
	r.Oncall = append(r.Oncall, *a)
}

/* Read Access
 *
 */
type somaOncallReadHandler struct {
	input     chan somaOncallRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaOncallReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.OncallList: r.list_stmt,
		stmt.OncallShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`oncall`, err, statement)
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

func (r *somaOncallReadHandler) process(q *somaOncallRequest) {
	var (
		oncallId, oncallName string
		oncallNumber         int
		rows                 *sql.Rows
		err                  error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: oncall/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&oncallId, &oncallName)
			result.Append(err, &somaOncallResult{
				Oncall: proto.Oncall{
					Id:   oncallId,
					Name: oncallName,
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: oncall/show for %s", q.Oncall.Id)
		err = r.show_stmt.QueryRow(q.Oncall.Id).Scan(
			&oncallId,
			&oncallName,
			&oncallNumber,
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

		result.Append(err, &somaOncallResult{
			Oncall: proto.Oncall{
				Id:     oncallId,
				Name:   oncallName,
				Number: strconv.Itoa(oncallNumber),
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaOncallWriteHandler struct {
	input    chan somaOncallRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	upd_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaOncallWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.OncallAdd:    w.add_stmt,
		stmt.OncallUpdate: w.upd_stmt,
		stmt.OncallDel:    w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`oncall`, err, statement)
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

func (w *somaOncallWriteHandler) process(q *somaOncallRequest) {
	var (
		name   sql.NullString
		number sql.NullInt64
		res    sql.Result
		n      int // ensure err not redeclared in if block
		err    error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: oncall/add for %s", q.Oncall.Name)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Oncall.Name,
			q.Oncall.Number,
		)
		q.Oncall.Id = id.String()
	case "update":
		w.reqLog.Printf("R: oncall/update for %s", q.Oncall.Id)
		// our update statement uses NULL to check which of the values
		// should be updated
		if q.Oncall.Name == "" {
			name = sql.NullString{String: "", Valid: false}
		} else {
			name = sql.NullString{String: q.Oncall.Name, Valid: true}
		}

		if q.Oncall.Number != "" {
			n, err = strconv.Atoi(q.Oncall.Number)
			if err != nil {
				break
			}
			number = sql.NullInt64{Int64: int64(n), Valid: true}
		} else {
			number = sql.NullInt64{Int64: 0, Valid: false}
		}
		res, err = w.upd_stmt.Exec(
			name,
			number,
			q.Oncall.Id,
		)
	case "delete":
		w.reqLog.Printf("R: oncall/del for %s", q.Oncall.Id)
		res, err = w.del_stmt.Exec(
			q.Oncall.Id,
		)
	default:
		w.reqLog.Printf("R: unimplemented oncall/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaOncallResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaOncallResult{})
	default:
		result.Append(nil, &somaOncallResult{
			Oncall: q.Oncall,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaOncallReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaOncallWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
