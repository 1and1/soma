package main

import (
	"database/sql"
	"errors"
	"log"
)

// Message structs
type somaViewRequest struct {
	action string
	view   string
	rename string
	reply  chan []somaViewResult
}

type somaViewResult struct {
	err  error
	view string
}

/*  Read Access
 *
 */
type somaViewReadHandler struct {
	input     chan somaViewRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaViewReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare("SELECT view FROM soma.views;")
	if err != nil {
		log.Fatal(err)
	}
	r.show_stmt, err = r.conn.Prepare("SELECT view FROM soma.views WHERE view = $1;")
	if err != nil {
		log.Fatal(err)
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
	var view string
	var rows *sql.Rows
	var err error
	result := make([]somaViewResult, 0)

	switch q.action {
	case "list":
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaViewResult{
				err:  err,
				view: q.view,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&view)
			if err != nil {
				result = append(result, somaViewResult{
					err:  err,
					view: q.view,
				})
				err = nil
				continue
			}
			result = append(result, somaViewResult{
				err:  nil,
				view: view,
			})
		}
	case "show":
		err = r.show_stmt.QueryRow(q.view).Scan(&view)
		if err != nil {
			result = append(result, somaViewResult{
				err:  err,
				view: q.view,
			})
			q.reply <- result
			return
		}

		result = append(result, somaViewResult{
			err:  nil,
			view: view,
		})
	default:
		result = append(result, somaViewResult{
			err:  errors.New("not implemented"),
			view: "",
		})
	}
	q.reply <- result
}

/*
 * Write Access
 */

type somaViewWriteHandler struct {
	input    chan somaViewRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
}

func (w *somaViewWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
  INSERT INTO soma.views (view)
  SELECT $1 WHERE NOT EXISTS (
    SELECT view FROM soma.views WHERE view = $2
  );
  `)
	if err != nil {
		log.Fatal(err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
  DELETE FROM soma.views
  WHERE view = $1;
  `)
	if err != nil {
		log.Fatal(err)
	}
	defer w.del_stmt.Close()

	w.ren_stmt, err = w.conn.Prepare(`
  UPDATE soma.views SET view = $1
  WHERE view = $2;
  `)
	if err != nil {
		log.Fatal(err)
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

func (w *somaViewWriteHandler) process(q *somaViewRequest) {
	var res sql.Result
	var err error

	result := make([]somaViewResult, 0)
	switch q.action {
	case "add":
		res, err = w.add_stmt.Exec(q.view, q.view)
	case "delete":
		res, err = w.del_stmt.Exec(q.view)
	case "rename":
		res, err = w.ren_stmt.Exec(q.rename, q.view)
	default:
		result = append(result, somaViewResult{
			err:  errors.New("not implemented"),
			view: "",
		})
	}
	if err != nil {
		result = append(result, somaViewResult{
			err:  err,
			view: q.view,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	if rowCnt == 0 {
		result = append(result, somaViewResult{
			err:  errors.New("No rows affected"),
			view: q.view,
		})
	} else if rowCnt > 1 {
		result = append(result, somaViewResult{
			err:  errors.New("Too many rows affected"),
			view: q.view,
		})
	} else {
		result = append(result, somaViewResult{
			err:  nil,
			view: q.view,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
