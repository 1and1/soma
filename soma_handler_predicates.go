package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaPredicateRequest struct {
	action    string
	predicate somaproto.ProtoPredicate
	reply     chan []somaPredicateResult
}

type somaPredicateResult struct {
	rErr      error
	lErr      error
	predicate somaproto.ProtoPredicate
}

/* Read Access
 *
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

	r.list_stmt, err = r.conn.Prepare(`
SELECT predicate
FROM soma.configuration_predicates; `)
	if err != nil {
		log.Fatal(err)
	}

	r.show_stmt, err = r.conn.Prepare(`
SELECT predicate
FROM soma.configuration_predicates
WHERE predicate = $1;`)
	if err != nil {
		log.Fatal(err)
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

func (r *somaPredicateReadHandler) process(q *somaPredicateRequest) {
	var predicate string
	var rows *sql.Rows
	var err error
	result := make([]somaPredicateResult, 0)

	switch q.action {
	case "list":
		log.Printf("R: predicates/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaPredicateResult{
				rErr: err,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&predicate)
			if err != nil {
				result = append(result, somaPredicateResult{
					lErr: err,
				})
				err = nil
				continue
			}
			result = append(result, somaPredicateResult{
				predicate: somaproto.ProtoPredicate{
					Predicate: predicate,
				},
			})
		}
	case "show":
		log.Printf("R: predicates/show for %s", q.predicate.Predicate)
		err = r.show_stmt.QueryRow(q.predicate.Predicate).Scan(&predicate)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result = append(result, somaPredicateResult{
					rErr: err,
				})
			}
			q.reply <- result
			return
		}

		result = append(result, somaPredicateResult{
			predicate: somaproto.ProtoPredicate{
				Predicate: predicate,
			},
		})
	default:
		result = append(result, somaPredicateResult{
			rErr: errors.New("not implemented"),
		})
	}
	q.reply <- result
}

/* Write Access
 *
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

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.configuration_predicates (
	predicate)
SELECT $1 WHERE NOT EXISTS (
	SELECT predicate
	FROM soma.configuration_predicates
	WHERE predicate = $2);`)
	if err != nil {
		log.Fatal(err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.configuration_predicates
WHERE predicate = $1;`)
	if err != nil {
		log.Fatal(err)
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
	var res sql.Result
	var err error
	result := make([]somaPredicateResult, 0)

	switch q.action {
	case "add":
		log.Printf("R: predicates/add for %s", q.predicate.Predicate)
		res, err = w.add_stmt.Exec(
			q.predicate.Predicate,
			q.predicate.Predicate,
		)
	case "delete":
		log.Printf("R: predicates/del for %s", q.predicate.Predicate)
		res, err = w.del_stmt.Exec(
			q.predicate.Predicate,
		)
	default:
		log.Printf("R: unimplemented predicates/%s", q.action)
		result = append(result, somaPredicateResult{
			rErr: errors.New("not implemented"),
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaPredicateResult{
			rErr: err,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result = append(result, somaPredicateResult{
			lErr: errors.New("No rows affected"),
		})
	case rowCnt > 1:
		result = append(result, somaPredicateResult{
			lErr: fmt.Errorf("Too many rows affected: %d", rowCnt),
		})
	default:
		result = append(result, somaPredicateResult{
			predicate: q.predicate,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
