package main

import (
	"database/sql"
	"errors"

	log "github.com/Sirupsen/logrus"
)

// Message structs
type somaObjectStateRequest struct {
	action string
	state  string
	rename string
	reply  chan []somaObjectStateResult
}

type somaObjectStateResult struct {
	err   error
	state string
}

/*  Read Access
 *
 */
type somaObjectStateReadHandler struct {
	input     chan somaObjectStateRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaObjectStateReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare("SELECT object_state FROM soma.object_states;")
	if err != nil {
		r.errLog.Fatal(err)
	}
	r.show_stmt, err = r.conn.Prepare("SELECT object_state FROM soma.object_states WHERE object_state = $1;")
	if err != nil {
		r.errLog.Fatal(err)
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

func (r *somaObjectStateReadHandler) process(q *somaObjectStateRequest) {
	var state string
	var rows *sql.Rows
	var err error
	result := make([]somaObjectStateResult, 0)

	switch q.action {
	case "list":
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaObjectStateResult{
				err:   err,
				state: q.state,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&state)
			if err != nil {
				result = append(result, somaObjectStateResult{
					err:   err,
					state: q.state,
				})
				err = nil
				continue
			}
			result = append(result, somaObjectStateResult{
				err:   nil,
				state: state,
			})
		}
	case "show":
		err = r.show_stmt.QueryRow(q.state).Scan(&state)
		if err != nil {
			result = append(result, somaObjectStateResult{
				err:   err,
				state: q.state,
			})
			q.reply <- result
			return
		}

		result = append(result, somaObjectStateResult{
			err:   nil,
			state: state,
		})
	default:
		result = append(result, somaObjectStateResult{
			err:   errors.New("not implemented"),
			state: "",
		})
	}
	q.reply <- result
}

/*
 * Write Access
 */

type somaObjectStateWriteHandler struct {
	input    chan somaObjectStateRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaObjectStateWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
  INSERT INTO soma.object_states (object_state)
  SELECT $1 WHERE NOT EXISTS (
    SELECT object_state FROM soma.object_states WHERE object_state = $2
  );
  `)
	if err != nil {
		w.errLog.Fatal(err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
  DELETE FROM soma.object_states
  WHERE object_state = $1;
  `)
	if err != nil {
		w.errLog.Fatal(err)
	}
	defer w.del_stmt.Close()

	w.ren_stmt, err = w.conn.Prepare(`
  UPDATE soma.object_states SET object_state = $1
  WHERE object_state = $2;
  `)
	if err != nil {
		w.errLog.Fatal(err)
	}
	defer w.ren_stmt.Close()

	for {
		select {
		case <-w.shutdown:
			break
		case req := <-w.input:
			w.process(&req)
		}
	}
}

func (w *somaObjectStateWriteHandler) process(q *somaObjectStateRequest) {
	var res sql.Result
	var err error

	result := make([]somaObjectStateResult, 0)
	switch q.action {
	case "add":
		res, err = w.add_stmt.Exec(q.state, q.state)
	case "delete":
		res, err = w.del_stmt.Exec(q.state)
	case "rename":
		res, err = w.ren_stmt.Exec(q.rename, q.state)
	default:
		result = append(result, somaObjectStateResult{
			err:   errors.New("not implemented"),
			state: "",
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaObjectStateResult{
			err:   err,
			state: q.state,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	if rowCnt == 0 {
		result = append(result, somaObjectStateResult{
			err:   errors.New("No rows affected"),
			state: q.state,
		})
	} else if rowCnt > 1 {
		result = append(result, somaObjectStateResult{
			err:   errors.New("Too many rows affected"),
			state: q.state,
		})
	} else {
		result = append(result, somaObjectStateResult{
			err:   nil,
			state: q.state,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaObjectStateReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaObjectStateWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
