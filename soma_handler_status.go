package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaStatusRequest struct {
	action string
	Status somaproto.ProtoStatus
	reply  chan somaResult
}

type somaStatusResult struct {
	ResultError error
	Status      somaproto.ProtoStatus
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
}

func (r *somaStatusReadHandler) run() {
	var err error

	log.Println("Prepare: status/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT status
FROM   soma.check_instance_status;`)
	if err != nil {
		log.Fatal("status/list: ", err)
	}
	defer r.list_stmt.Close()

	r.show_stmt, err = r.conn.Prepare(`
SELECT status
FROM   soma.check_instance_status
WHERE  status = $1;`)
	if err != nil {
		log.Fatal("status/show: ", err)
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

func (r *somaStatusReadHandler) process(q *somaStatusRequest) {
	var (
		status string
		rows   *sql.Rows
		err    error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: status/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&status)
			result.Append(err, &somaStatusResult{
				Status: somaproto.ProtoStatus{
					Status: status,
				},
			})
		}
	case "show":
		log.Printf("R: status/show for %s", q.Status.Status)
		err = r.show_stmt.QueryRow(q.Status.Status).Scan(
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
			Status: somaproto.ProtoStatus{
				Status: status,
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
}

func (w *somaStatusWriteHandler) run() {
	var err error

	log.Println("Prepare: status/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.check_instance_status (
	status)
SELECT $1::varchar WHERE NOT EXISTS (
	SELECT status
	FROM   soma.check_instance_status
	WHERE  status = $1::varchar);`)
	if err != nil {
		log.Fatal("status/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: status/delete")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.check_instance_status
WHERE  status = $1;`)
	if err != nil {
		log.Fatal("status/delete: ", err)
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

func (w *somaStatusWriteHandler) process(q *somaStatusRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: status/add for %s", q.Status.Status)
		res, err = w.add_stmt.Exec(
			q.Status.Status,
		)
	case "delete":
		log.Printf("R: status/del for %s", q.Status.Status)
		res, err = w.del_stmt.Exec(
			q.Status.Status,
		)
	default:
		log.Printf("R: unimplemented status/%s", q.action)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
