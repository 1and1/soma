package main

import (
	"database/sql"
	"errors"

	log "github.com/Sirupsen/logrus"
)

// Message structs
type somaEnvironmentRequest struct {
	action      string
	environment string
	rename      string
	reply       chan []somaEnvironmentResult
}

type somaEnvironmentResult struct {
	err         error
	environment string
}

/*  Read Access
 *
 */
type somaEnvironmentReadHandler struct {
	input     chan somaEnvironmentRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaEnvironmentReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare("SELECT environment FROM soma.environments;")
	if err != nil {
		log.Fatal(err)
	}
	r.show_stmt, err = r.conn.Prepare("SELECT environment FROM soma.environments WHERE environment = $1;")
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

func (r *somaEnvironmentReadHandler) process(q *somaEnvironmentRequest) {
	var environment string
	var rows *sql.Rows
	var err error
	result := make([]somaEnvironmentResult, 0)

	switch q.action {
	case "list":
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if err != nil {
			result = append(result, somaEnvironmentResult{
				err:         err,
				environment: q.environment,
			})
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&environment)
			if err != nil {
				result = append(result, somaEnvironmentResult{
					err:         err,
					environment: q.environment,
				})
				err = nil
				continue
			}
			result = append(result, somaEnvironmentResult{
				err:         nil,
				environment: environment,
			})
		}
	case "show":
		err = r.show_stmt.QueryRow(q.environment).Scan(&environment)
		if err != nil {
			result = append(result, somaEnvironmentResult{
				err:         err,
				environment: ``,
			})
			q.reply <- result
			return
		}

		result = append(result, somaEnvironmentResult{
			err:         nil,
			environment: environment,
		})
	default:
		result = append(result, somaEnvironmentResult{
			err:         errors.New("not implemented"),
			environment: "",
		})
	}
	q.reply <- result
}

/*
 * Write Access
 */

type somaEnvironmentWriteHandler struct {
	input    chan somaEnvironmentRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaEnvironmentWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
  INSERT INTO soma.environments (environment)
  SELECT $1 WHERE NOT EXISTS (
    SELECT environment FROM soma.environments WHERE environment = $2
  );
  `)
	if err != nil {
		log.Fatal(err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
  DELETE FROM soma.environments
  WHERE environment = $1;
  `)
	if err != nil {
		log.Fatal(err)
	}
	defer w.del_stmt.Close()

	w.ren_stmt, err = w.conn.Prepare(`
  UPDATE soma.environments SET environment = $1
  WHERE environment = $2;
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

func (w *somaEnvironmentWriteHandler) process(q *somaEnvironmentRequest) {
	var res sql.Result
	var err error

	result := make([]somaEnvironmentResult, 0)
	switch q.action {
	case "add":
		res, err = w.add_stmt.Exec(q.environment, q.environment)
	case "delete":
		res, err = w.del_stmt.Exec(q.environment)
	case "rename":
		res, err = w.ren_stmt.Exec(q.rename, q.environment)
	default:
		result = append(result, somaEnvironmentResult{
			err:         errors.New("not implemented"),
			environment: "",
		})
		q.reply <- result
		return
	}
	if err != nil {
		result = append(result, somaEnvironmentResult{
			err:         err,
			environment: q.environment,
		})
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	if rowCnt == 0 {
		result = append(result, somaEnvironmentResult{
			err:         errors.New("No rows affected"),
			environment: q.environment,
		})
	} else if rowCnt > 1 {
		result = append(result, somaEnvironmentResult{
			err:         errors.New("Too many rows affected"),
			environment: q.environment,
		})
	} else {
		result = append(result, somaEnvironmentResult{
			err:         nil,
			environment: q.environment,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaEnvironmentReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaEnvironmentWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
