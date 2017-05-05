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

// NodeWrite handles write requests for nodes
type NodeWrite struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtPurge  *sql.Stmt
	stmtRemove *sql.Stmt
	stmtUpdate *sql.Stmt
	appLog     *logrus.Logger
	reqLog     *logrus.Logger
	errLog     *logrus.Logger
}

// newNodeWrite return a new NodeWrite handler with input buffer of
// length
func newNodeWrite(length int) (w *NodeWrite) {
	w = &NodeWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *NodeWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.reqLog = l[2]
}

// run is the event loop for NodeWrite
func (w *NodeWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.NodeAdd:    w.stmtAdd,
		stmt.NodeUpdate: w.stmtUpdate,
		stmt.NodeRemove: w.stmtRemove,
		stmt.NodePurge:  w.stmtPurge,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`node`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-w.Shutdown:
			break runloop
		case req := <-w.Input:
			go func() {
				w.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (w *NodeWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case `add`:
		w.add(q, &result)
	case `remove`:
		w.remove(q, &result)
	case `update`:
		w.update(q, &result)
	case `purge`:
		w.purge(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add inserts a new node
func (w *NodeWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	q.Node.Id = uuid.NewV4().String()
	if q.Node.ServerId == `` {
		q.Node.ServerId = `00000000-0000-0000-0000-000000000000`
	}
	if res, err = w.stmtAdd.Exec(
		q.Node.Id,
		q.Node.AssetId,
		q.Node.Name,
		q.Node.TeamId,
		q.Node.ServerId,
		q.Node.State,
		q.Node.IsOnline,
		false,
		q.User,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Node = append(mr.Node, q.Node)
	}
}

// remove delete a node
func (w *NodeWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Node.Id,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Node = append(mr.Node, q.Node)
	}
}

// update refreshes a node
func (w *NodeWrite) update(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtUpdate.Exec(
		q.Node.AssetId,
		q.Node.Name,
		q.Node.TeamId,
		q.Node.ServerId,
		q.Node.IsOnline,
		q.Node.IsDeleted,
		q.Node.Id,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Node = append(mr.Node, q.Node)
	}
}

// purge removes a node flagged as deleted
func (w *NodeWrite) purge(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtPurge.Exec(
		q.Node.Id,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Node = append(mr.Node, q.Node)
	}
}

// shutdown signals the handler to shut down
func (w *NodeWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
