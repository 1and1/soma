package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaAttributeRequest struct {
	action    string
	Attribute somaproto.Attribute
	reply     chan somaResult
}

type somaAttributeResult struct {
	ResultError error
	Attribute   somaproto.Attribute
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
}

func (r *somaAttributeReadHandler) run() {
	var err error

	log.Println("Prepare: attribute/list")
	if r.list_stmt, err = r.conn.Prepare(stmtAttributeList); err != nil {
		log.Fatal("attribute/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: attribute/show")
	if r.show_stmt, err = r.conn.Prepare(stmtAttributeShow); err != nil {
		log.Fatal("attribute/show: ", err)
	}
	defer r.show_stmt.Close()

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
		log.Printf("R: attributes/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&attribute, &cardinality)
			result.Append(err, &somaAttributeResult{
				Attribute: somaproto.Attribute{
					Attribute:   attribute,
					Cardinality: cardinality,
				},
			})
		}
	case "show":
		log.Printf("R: attribute/show for %s", q.Attribute.Attribute)
		err = r.show_stmt.QueryRow(q.Attribute.Attribute).Scan(
			&attribute,
			&cardinality,
		)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaAttributeResult{
			Attribute: somaproto.Attribute{
				Attribute:   attribute,
				Cardinality: cardinality,
			},
		})
	default:
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
}

func (w *somaAttributeWriteHandler) run() {
	var err error

	log.Println("Prepare: attribute/add")
	if w.add_stmt, err = w.conn.Prepare(stmtAttributeAdd); err != nil {
		log.Fatal("attribute/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: attribute/delete")
	if w.del_stmt, err = w.conn.Prepare(stmtAttributeDelete); err != nil {
		log.Fatal("attribute/delete: ", err)
	}
	defer w.del_stmt.Close()

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
		log.Printf("R: attributes/add for %s", q.Attribute.Attribute)
		res, err = w.add_stmt.Exec(
			q.Attribute.Attribute,
			q.Attribute.Cardinality,
		)
	case "delete":
		log.Printf("R: attributes/del for %s", q.Attribute.Attribute)
		res, err = w.del_stmt.Exec(
			q.Attribute.Attribute,
		)
	default:
		log.Printf("R: unimplemented attributes/%s", q.action)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
