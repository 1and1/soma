package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaUnitRequest struct {
	action string
	Unit   proto.Unit
	reply  chan somaResult
}

type somaUnitResult struct {
	ResultError error
	Unit        proto.Unit
}

func (a *somaUnitResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Units = append(r.Units, somaUnitResult{ResultError: err})
	}
}

func (a *somaUnitResult) SomaAppendResult(r *somaResult) {
	r.Units = append(r.Units, *a)
}

/* Read Access
 */
type somaUnitReadHandler struct {
	input     chan somaUnitRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaUnitReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.UnitList: r.list_stmt,
		stmt.UnitShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`unit`, err, statement)
		}
		defer prepStmt.Close()
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

func (r *somaUnitReadHandler) process(q *somaUnitRequest) {
	var (
		unit, name string
		rows       *sql.Rows
		err        error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.appLog.Printf("R: units/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&unit)
			result.Append(err, &somaUnitResult{
				Unit: proto.Unit{
					Unit: unit,
				},
			})
		}
		if err = rows.Err(); err != nil {
			result.Append(err, &somaUnitResult{})
			err = nil
		}
	case "show":
		r.appLog.Printf("R: units/show for %s", q.Unit.Unit)
		err = r.show_stmt.QueryRow(q.Unit.Unit).Scan(
			&unit,
			&name,
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

		result.Append(err, &somaUnitResult{
			Unit: proto.Unit{
				Unit: unit,
				Name: name,
			},
		})
	default:
		r.errLog.Printf("R: unimplemented units/%s", q.action)
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaUnitWriteHandler struct {
	input    chan somaUnitRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaUnitWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.UnitAdd: w.add_stmt,
		stmt.UnitDel: w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`unit`, err, statement)
		}
		defer prepStmt.Close()
	}

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

func (w *somaUnitWriteHandler) process(q *somaUnitRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.appLog.Printf("R: units/add for %s", q.Unit.Unit)
		res, err = w.add_stmt.Exec(
			q.Unit.Unit,
			q.Unit.Name,
		)
	case "delete":
		w.appLog.Printf("R: units/del for %s", q.Unit.Unit)
		res, err = w.del_stmt.Exec(
			q.Unit.Unit,
		)
	default:
		w.errLog.Printf("R: unimplemented units/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaUnitResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaUnitResult{})
	default:
		result.Append(nil, &somaUnitResult{
			Unit: q.Unit,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaUnitReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaUnitWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
