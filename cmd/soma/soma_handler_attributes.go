package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaAttributeRequest struct {
	action    string
	Attribute proto.Attribute
	reply     chan somaResult
}

type somaAttributeResult struct {
	ResultError error
	Attribute   proto.Attribute
}

func (a *somaAttributeResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Attributes = append(r.Attributes,
			somaAttributeResult{ResultError: err})
	}
}

func (a *somaAttributeResult) SomaAppendResult(r *somaResult) {
	r.Attributes = append(r.Attributes, *a)
}

/* Read Access
 */
type somaAttributeReadHandler struct {
	input     chan somaAttributeRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaAttributeReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AttributeList: r.list_stmt,
		stmt.AttributeShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(``, err, statement)
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

func (r *somaAttributeReadHandler) process(q *somaAttributeRequest) {
	var (
		attribute, cardinality string
		rows                   *sql.Rows
		err                    error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.appLog.Printf("R: attributes/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&attribute, &cardinality)
			result.Append(err, &somaAttributeResult{
				Attribute: proto.Attribute{
					Name:        attribute,
					Cardinality: cardinality,
				},
			})
		}
	case "show":
		r.appLog.Printf("R: attribute/show for %s", q.Attribute.Name)
		err = r.show_stmt.QueryRow(q.Attribute.Name).Scan(
			&attribute,
			&cardinality,
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

		result.Append(err, &somaAttributeResult{
			Attribute: proto.Attribute{
				Name:        attribute,
				Cardinality: cardinality,
			},
		})
	default:
		r.errLog.Printf("R: unimplemented attributes/%s", q.action)
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaAttributeWriteHandler struct {
	input    chan somaAttributeRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaAttributeWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AttributeAdd:    w.add_stmt,
		stmt.AttributeDelete: w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`attribute`, err, statement)
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

func (w *somaAttributeWriteHandler) process(q *somaAttributeRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.appLog.Printf("R: attributes/add for %s", q.Attribute.Name)
		res, err = w.add_stmt.Exec(
			q.Attribute.Name,
			q.Attribute.Cardinality,
		)
	case "delete":
		w.appLog.Printf("R: attributes/del for %s", q.Attribute.Name)
		res, err = w.del_stmt.Exec(
			q.Attribute.Name,
		)
	default:
		w.errLog.Printf("R: unimplemented attributes/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaAttributeResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaAttributeResult{})
	default:
		result.Append(nil, &somaAttributeResult{
			Attribute: q.Attribute,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaAttributeReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaAttributeWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
