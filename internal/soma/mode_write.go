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

// ModeWrite handles write requests for modes
type ModeWrite struct {
	Input       chan msg.Request
	Shutdown    chan struct{}
	conn        *sql.DB
	stmtCreate  *sql.Stmt
	stmtDestroy *sql.Stmt
	appLog      *logrus.Logger
	reqLog      *logrus.Logger
	errLog      *logrus.Logger
}

// newModeWrite return a new ModeWrite handler with input buffer of
// length
func newModeWrite(length int) (w *ModeWrite) {
	w = &ModeWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *ModeWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for ModeWrite
func (w *ModeWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ModeAdd: w.stmtCreate,
		stmt.ModeDel: w.stmtDestroy,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`mode`, err, stmt.Name(statement))
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
func (w *ModeWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case `create`:
		w.create(q, &result)
	case `destroy`:
		w.destroy(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// create inserts a new mode
func (w *ModeWrite) create(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtCreate.Exec(
		q.Mode.Mode,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Mode = append(mr.Mode, q.Mode)
	}
}

// destroy deletes a mode
func (w *ModeWrite) destroy(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtDestroy.Exec(
		q.Mode.Mode,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Mode = append(mr.Mode, q.Mode)
	}
}

// shutdown signals the handler to shut down
func (w *ModeWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
