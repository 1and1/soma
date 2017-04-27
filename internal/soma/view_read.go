/*-
 * Copyright (c) 2015-2017, Jörg Pernfuß
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

// ViewRead handles read requests for views
type ViewRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// run is the event loop for ViewRead
func (r *ViewRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ViewList: r.stmtList,
		stmt.ViewShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`view`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	for {
		select {
		case <-r.Shutdown:
			break
		case req := <-r.Input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (r *ViewRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case `show`:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all views
func (r *ViewRead) list(q *msg.Request, mr *msg.Result) {
}

// show returns the details of a specific view
func (r *ViewRead) show(q *msg.Request, mr *msg.Result) {
}

// shutdown signals the handler to shut down
func (r *ViewRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
