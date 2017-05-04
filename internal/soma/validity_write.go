/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/Sirupsen/logrus"
)

// ValidityWrite handles write requests for validity definitions
type ValidityWrite struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
	appLog     *logrus.Logger
	reqLog     *logrus.Logger
	errLog     *logrus.Logger
}

// newValidityWrite returns a new ValidityWrite handler
func newValidityWrite() (w ValidityWrite) {
	w = ValidityWrite{}
	w.Input = make(chan msg.Request, 64)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *ValidityWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for ValidityWrite
func (w *ValidityWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ValidityAdd: w.stmtAdd,
		stmt.ValidityDel: w.stmtRemove,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`validity`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-w.Shutdown:
			break runloop
		case req := <-w.Input:
			w.process(&req)
		}
	}
}

// process is the request dispatcher
func (w *ValidityWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case `add`:
		w.add(q, &result)
	case `remove`:
		w.remove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add inserts a new validity definition
func (w *ValidityWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtAdd.Exec(
		q.Validity.SystemProperty,
		q.Validity.ObjectType,
		q.Validity.Inherited,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Validity = append(mr.Validity, q.Validity)
	}
}

// remove deletes a validity definition
func (w *ValidityWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Validity.SystemProperty,
		q.Validity.ObjectType,
		q.Validity.Inherited,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Validity = append(mr.Validity, q.Validity)
	}
}

// shutdown signals the handler to shut down
func (w *ValidityWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
