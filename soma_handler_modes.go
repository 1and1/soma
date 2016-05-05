package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaModeRequest struct {
	action string
	Mode   somaproto.ProtoMode
	reply  chan somaResult
}

type somaModeResult struct {
	ResultError error
	Mode        somaproto.ProtoMode
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
}

func (r *somaModeReadHandler) run() {
	var err error

	log.Println("Prepare: mode/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT monitoring_system_mode
FROM   soma.monitoring_system_modes; `)
	if err != nil {
		log.Fatal("mode/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: mode/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT monitoring_system_mode
FROM   soma.monitoring_system_modes
WHERE  monitoring_system_mode = $1;`)
	if err != nil {
		log.Fatal("mode/show: ", err)
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

func (r *somaModeReadHandler) process(q *somaModeRequest) {
	var (
		mode string
		rows *sql.Rows
		err  error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: modes/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&mode)
			result.Append(err, &somaModeResult{
				Mode: somaproto.ProtoMode{
					Mode: mode,
				},
			})
		}
	case "show":
		log.Printf("R: mode/show for %s", q.Mode.Mode)
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
			Mode: somaproto.ProtoMode{
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
}

func (w *somaModeWriteHandler) run() {
	var err error

	log.Println("Prepare: mode/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.monitoring_system_modes (
	monitoring_system_mode)
SELECT $1::varchar WHERE NOT EXISTS (
	SELECT monitoring_system_mode
	FROM   soma.monitoring_system_modes
	WHERE  monitoring_system_mode = $1::varchar);`)
	if err != nil {
		log.Fatal("mode/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: mode/delete")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.monitoring_system_modes
WHERE  monitoring_system_mode = $1;`)
	if err != nil {
		log.Fatal("mode/delete: ", err)
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

func (w *somaModeWriteHandler) process(q *somaModeRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: modes/add for %s", q.Mode.Mode)
		res, err = w.add_stmt.Exec(
			q.Mode.Mode,
		)
	case "delete":
		log.Printf("R: modes/del for %s", q.Mode.Mode)
		res, err = w.del_stmt.Exec(
			q.Mode.Mode,
		)
	default:
		log.Printf("R: unimplemented modes/%s", q.action)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
