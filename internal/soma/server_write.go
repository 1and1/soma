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
	uuid "github.com/satori/go.uuid"
)

// ServerWrite handles write requests for servers
type ServerWrite struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
	stmtPurge  *sql.Stmt
	stmtUpdate *sql.Stmt
	appLog     *logrus.Logger
	reqLog     *logrus.Logger
	errLog     *logrus.Logger
}

// newServerWrite return a new ServerWrite handler with input buffer of
// length
func newServerWrite(length int) (w *ServerWrite) {
	w = &ServerWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *ServerWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for ServerWrite
func (w *ServerWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AddServers:    w.stmtAdd,
		stmt.DeleteServers: w.stmtRemove,
		stmt.PurgeServers:  w.stmtPurge,
		stmt.UpdateServers: w.stmtUpdate,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`server`, err, stmt.Name(statement))
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
func (w *ServerWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case `add`:
		w.add(q, &result)
	case `remove`:
		w.remove(q, &result)
	case `purge`:
		w.purge(q, &result)
	case `update`:
		w.update(q, &result)
	case `insert-null`:
		w.insertNull(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add inserts a new server
func (w *ServerWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	q.Server.Id = uuid.NewV4().String()
	if res, err = w.stmtAdd.Exec(
		q.Server.Id,
		q.Server.AssetId,
		q.Server.Datacenter,
		q.Server.Location,
		q.Server.Name,
		q.Server.IsOnline,
		false,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Server = append(mr.Server, q.Server)
	}
}

// insertNull inserts the default server for nodes that have none
func (w *ServerWrite) insertNull(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	q.Server.Id = `00000000-0000-0000-0000-000000000000`
	q.Server.AssetId = 0
	q.Server.Location = `none`
	q.Server.Name = `soma-null-server`
	q.Server.IsOnline = true
	q.Server.IsDeleted = false
	if res, err = w.stmtAdd.Exec(
		q.Server.Id,
		q.Server.AssetId,
		q.Server.Datacenter,
		q.Server.Location,
		q.Server.Name,
		q.Server.IsOnline,
		q.Server.IsDeleted,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Server = append(mr.Server, q.Server)
	}
}

// remove marks a server as deleted
func (w *ServerWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtRemove.Exec(
		q.Server.Id,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Server = append(mr.Server, q.Server)
	}
}

// purge deletes servers marked as deleted from the database
func (w *ServerWrite) purge(q *msg.Request, mr *msg.Result) {
	var err error

	if _, err = w.stmtPurge.Exec(
		q.Server.Id,
	); err != nil {
		mr.ServerError(err)
		return
	}
	mr.Server = append(mr.Server, q.Server)
}

// update refreshes a servers details
func (w *ServerWrite) update(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtUpdate.Exec(
		q.Server.Id,
		q.Server.AssetId,
		q.Server.Datacenter,
		q.Server.Location,
		q.Server.Name,
		q.Server.IsOnline,
		q.Server.IsDeleted,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Server = append(mr.Server, q.Server)
	}
}

// shutdown signals the handler to shut down
func (w *ServerWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
