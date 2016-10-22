package main

import (
	"database/sql"
	"errors"

	"github.com/1and1/soma/internal/stmt"
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

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.EnvironmentList: r.list_stmt,
		stmt.EnvironmentShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`environment`, err, statement)
		}
		defer prepStmt.Close()
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

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.EnvironmentAdd:    w.add_stmt,
		stmt.EnvironmentDel:    w.del_stmt,
		stmt.EnvironmentRename: w.ren_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`environment`, err, statement)
		}
		defer prepStmt.Close()
	}

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
