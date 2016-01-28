package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaStatusRequest struct {
	action string
	status somaproto.ProtoStatus
	reply  chan []somaStatusResult
}

type somaStatusResult struct {
	rErr   error
	lErr   error
	status somaproto.ProtoStatus
}

/* Read Access
 *
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

	r.list_stmt, err = r.conn.Prepare(`
SELECT status
FROM soma.check_instance_status; `)
	if err != nil {
		log.Fatal(err)
	}

	r.show_stmt, err = r.conn.Prepare(`
SELECT status
FROM soma.check_instance_status
WHERE status = $1;`)
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

func (r *somaStatusReadHandler) process(q *somaStatusRequest) {
	var status string
	var rows *sql.Rows
	var err error
	result := make([]somaStatusResult, 0)

	switch q.action {
	case "list":
		log.Printf("R: status/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaStatusResult{
				rErr: err,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&status)
			if err != nil {
				result = append(result, somaStatusResult{
					lErr: err,
				})
				err = nil
				continue
			}
			result = append(result, somaStatusResult{
				status: somaproto.ProtoStatus{
					Status: status,
				},
			})
		}
	case "show":
		log.Printf("R: status/show for %s", q.status.Status)
		err = r.show_stmt.QueryRow(q.status.Status).Scan(&status)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result = append(result, somaStatusResult{
					rErr: err,
				})
			}
			q.reply <- result
			return
		}

		result = append(result, somaStatusResult{
			status: somaproto.ProtoStatus{
				Status: status,
			},
		})
	default:
		result = append(result, somaStatusResult{
			rErr: errors.New("not implemented"),
		})
	}
	q.reply <- result
}

/* Write Access
 *
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

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.check_instance_status (
	status)
SELECT $1 WHERE NOT EXISTS (
	SELECT status
	FROM soma.check_instance_status
	WHERE status = $2);`)
	if err != nil {
		log.Fatal(err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.check_instance_status
WHERE status = $1;`)
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

func (w *somaStatusWriteHandler) process(q *somaStatusRequest) {
	var res sql.Result
	var err error
	result := make([]somaStatusResult, 0)

	switch q.action {
	case "add":
		log.Printf("R: status/add for %s", q.status.Status)
		res, err = w.add_stmt.Exec(
			q.status.Status,
			q.status.Status,
		)
	case "delete":
		log.Printf("R: statuss/del for %s", q.status.Status)
		res, err = w.del_stmt.Exec(
			q.status.Status,
		)
	default:
		log.Printf("R: unimplemented statuss/%s", q.action)
		result = append(result, somaStatusResult{
			rErr: errors.New("not implemented"),
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaStatusResult{
			rErr: err,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result = append(result, somaStatusResult{
			lErr: errors.New("No rows affected"),
		})
	case rowCnt > 1:
		result = append(result, somaStatusResult{
			lErr: fmt.Errorf("Too many rows affected: %d", rowCnt),
		})
	default:
		result = append(result, somaStatusResult{
			status: q.status,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix