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
	"strconv"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

// OncallWrite handles write requests for oncall
type OncallWrite struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	conn        *sql.DB
	stmtCreate  *sql.Stmt
	stmtUpdate  *sql.Stmt
	stmtDestroy *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newOncallWrite return a new OncallWrite handler with input buffer of
// length
func newOncallWrite(length int) (w *OncallWrite) {
	w = &OncallWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *OncallWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for OncallWrite
func (w *OncallWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.OncallAdd:    w.stmtCreate,
		stmt.OncallUpdate: w.stmtUpdate,
		stmt.OncallDel:    w.stmtDestroy,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`oncall`, err, stmt.Name(statement))
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
func (w *OncallWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case `create`:
		w.create(q, &result)
	case `destroy`:
		w.destroy(q, &result)
	case `update`:
		w.update(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// create inserts a new oncall
func (w *OncallWrite) create(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	q.Oncall.Id = uuid.NewV4().String()
	if res, err = w.stmtCreate.Exec(
		q.Oncall.Id,
		q.Oncall.Name,
		q.Oncall.Number,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Oncall = append(mr.Oncall, q.Oncall)
	}
}

// destroy removes an oncall entry
func (w *OncallWrite) destroy(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtDestroy.Exec(
		q.Oncall.Id,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Oncall = append(mr.Oncall, q.Oncall)
	}
}

// update refreshes an oncall entry
func (w *OncallWrite) update(q *msg.Request, mr *msg.Result) {
	var (
		name   sql.NullString
		number sql.NullInt64
		res    sql.Result
		n      int // ensure err not redeclared in if block
		err    error
	)

	// our update statement uses NULL to check which of the values
	// should be updated - can be both
	if q.Oncall.Name != "" {
		name = sql.NullString{String: q.Oncall.Name, Valid: true}
	} else {
		name = sql.NullString{String: "", Valid: false}
	}

	if q.Oncall.Number != "" {
		if n, err = strconv.Atoi(q.Oncall.Number); err != nil {
			mr.ServerError(err, q.Section)
			return
		}
		number = sql.NullInt64{Int64: int64(n), Valid: true}
	} else {
		number = sql.NullInt64{Int64: 0, Valid: false}
	}
	if res, err = w.stmtUpdate.Exec(
		name,
		number,
		q.Oncall.Id,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Oncall = append(mr.Oncall, q.Oncall)
	}
}

// shutdown signals the handler to shut down
func (w *OncallWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
