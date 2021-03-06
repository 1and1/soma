package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

// Message structs
type somaViewRequest struct {
	action string
	name   string
	View   proto.View
	reply  chan somaResult
}

type somaViewResult struct {
	ResultError error
	View        proto.View
}

func (a *somaViewResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Views = append(r.Views, somaViewResult{ResultError: err})
	}
}

func (a *somaViewResult) SomaAppendResult(r *somaResult) {
	r.Views = append(r.Views, *a)
}

/*  Read Access
 */
type somaViewReadHandler struct {
	input     chan somaViewRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaViewReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ViewList: r.list_stmt,
		stmt.ViewShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`view`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	for {
		select {
		case <-r.shutdown:
			break
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaViewReadHandler) process(q *somaViewRequest) {
	var (
		view string
		rows *sql.Rows
		err  error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.appLog.Printf("R: view/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&view)
			result.Append(err, &somaViewResult{
				View: proto.View{
					Name: view,
				},
			})
		}
		if err = rows.Err(); err != nil {
			result.Append(err, &somaViewResult{})
			err = nil
		}
	case "show":
		r.appLog.Printf("R: view/show for %s", q.View.Name)
		err = r.show_stmt.QueryRow(q.View.Name).Scan(
			&view,
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

		result.Append(err, &somaViewResult{
			View: proto.View{
				Name: view,
			},
		})
	default:
		r.errLog.Printf("R: unimplemented levels/%s", q.action)
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */

type somaViewWriteHandler struct {
	input    chan somaViewRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaViewWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ViewAdd:    w.add_stmt,
		stmt.ViewDel:    w.del_stmt,
		stmt.ViewRename: w.ren_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`view`, err, stmt.Name(statement))
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

func (w *somaViewWriteHandler) process(q *somaViewRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.appLog.Printf("R: view/add for %s", q.View.Name)
		res, err = w.add_stmt.Exec(
			q.View.Name,
		)
	case "delete":
		w.appLog.Printf("R: view/delete for %s", q.View.Name)
		res, err = w.del_stmt.Exec(
			q.View.Name,
		)
	case "rename":
		w.appLog.Printf("R: view/rename for %s", q.name)
		res, err = w.ren_stmt.Exec(
			q.View.Name,
			q.name,
		)
	default:
		w.errLog.Printf("R: unimplemented levels/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaViewResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaViewResult{})
	default:
		result.Append(nil, &somaViewResult{
			View: q.View,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaViewReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaViewWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
