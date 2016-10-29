package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaStatusRequest struct {
	action string
	Status proto.Status
	reply  chan somaResult
}

type somaStatusResult struct {
	ResultError error
	Status      proto.Status
}

func (a *somaStatusResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Status = append(r.Status, somaStatusResult{ResultError: err})
	}
}

func (a *somaStatusResult) SomaAppendResult(r *somaResult) {
	r.Status = append(r.Status, *a)
}

/* Read Access
 */
type somaStatusReadHandler struct {
	input     chan somaStatusRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaStatusReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.StatusList: r.list_stmt,
		stmt.StatusShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`status`, err, stmt.Name(statement))
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

func (r *somaStatusReadHandler) process(q *somaStatusRequest) {
	var (
		status string
		rows   *sql.Rows
		err    error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: status/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&status)
			result.Append(err, &somaStatusResult{
				Status: proto.Status{
					Name: status,
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: status/show for %s", q.Status.Name)
		err = r.show_stmt.QueryRow(q.Status.Name).Scan(
			&status,
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

		result.Append(err, &somaStatusResult{
			Status: proto.Status{
				Name: status,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaStatusWriteHandler struct {
	input    chan somaStatusRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaStatusWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.StatusAdd: w.add_stmt,
		stmt.StatusDel: w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`status`, err, stmt.Name(statement))
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

func (w *somaStatusWriteHandler) process(q *somaStatusRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: status/add for %s", q.Status.Name)
		res, err = w.add_stmt.Exec(
			q.Status.Name,
		)
	case "delete":
		w.reqLog.Printf("R: status/del for %s", q.Status.Name)
		res, err = w.del_stmt.Exec(
			q.Status.Name,
		)
	default:
		w.reqLog.Printf("R: unimplemented status/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaStatusResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaStatusResult{})
	default:
		result.Append(nil, &somaStatusResult{
			Status: q.Status,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaStatusReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaStatusWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
