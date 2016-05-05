package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaPredicateRequest struct {
	action    string
	Predicate somaproto.ProtoPredicate
	reply     chan somaResult
}

type somaPredicateResult struct {
	ResultError error
	Predicate   somaproto.ProtoPredicate
}

func (a *somaPredicateResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Predicates = append(r.Predicates,
			somaPredicateResult{ResultError: err})
	}
}

func (a *somaPredicateResult) SomaAppendResult(r *somaResult) {
	r.Predicates = append(r.Predicates, *a)
}

/* Read Access
 */
type somaPredicateReadHandler struct {
	input     chan somaPredicateRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaPredicateReadHandler) run() {
	var err error

	log.Println("Prepare: predicate/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT predicate
FROM   soma.configuration_predicates; `)
	if err != nil {
		log.Fatal("predicate/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: predicate/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT predicate
FROM   soma.configuration_predicates
WHERE  predicate = $1;`)
	if err != nil {
		log.Fatal("predicate/show: ", err)
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

func (r *somaPredicateReadHandler) process(q *somaPredicateRequest) {
	var (
		predicate string
		rows      *sql.Rows
		err       error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: predicates/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&predicate)
			result.Append(err, &somaPredicateResult{
				Predicate: somaproto.ProtoPredicate{
					Predicate: predicate,
				},
			})
		}
	case "show":
		log.Printf("R: predicate/show for %s", q.Predicate.Predicate)
		err = r.show_stmt.QueryRow(q.Predicate.Predicate).Scan(
			&predicate,
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

		result.Append(err, &somaPredicateResult{
			Predicate: somaproto.ProtoPredicate{
				Predicate: predicate,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaPredicateWriteHandler struct {
	input    chan somaPredicateRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaPredicateWriteHandler) run() {
	var err error

	log.Println("Prepare: predicate/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.configuration_predicates (
	predicate)
SELECT $1::varchar WHERE NOT EXISTS (
	SELECT predicate
	FROM   soma.configuration_predicates
	WHERE  predicate = $1::varchar);`)
	if err != nil {
		log.Fatal("predicate/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: predicate/delete")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.configuration_predicates
WHERE  predicate = $1;`)
	if err != nil {
		log.Fatal("predicate/delete: ", err)
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

func (w *somaPredicateWriteHandler) process(q *somaPredicateRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: predicates/add for %s", q.Predicate.Predicate)
		res, err = w.add_stmt.Exec(
			q.Predicate.Predicate,
		)
	case "delete":
		log.Printf("R: predicates/del for %s", q.Predicate.Predicate)
		res, err = w.del_stmt.Exec(
			q.Predicate.Predicate,
		)
	default:
		log.Printf("R: unimplemented predicates/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaPredicateResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaPredicateResult{})
	default:
		result.Append(nil, &somaPredicateResult{
			Predicate: q.Predicate,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
