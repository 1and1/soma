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
	"github.com/1and1/soma/lib/proto"
	"github.com/Sirupsen/logrus"
)

// ValidityRead handles read requests for validity definitions
type ValidityRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// newValidityRead returns a new ValidityRead handler with input buffer
// of length
func newValidityRead(length int) (r *ValidityRead) {
	r = &ValidityRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *ValidityRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for ValidityRead
func (r *ValidityRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ValidityList: r.stmtList,
		stmt.ValidityShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`validity`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-r.Shutdown:
			break runloop
		case req := <-r.Input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (r *ValidityRead) process(q *msg.Request) {
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

// list returns all validity definitions
func (r *ValidityRead) list(q *msg.Request, mr *msg.Result) {
	var (
		systemProperty, objectType string
		rows                       *sql.Rows
		err                        error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&systemProperty,
			&objectType,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Validity = append(mr.Validity, proto.Validity{
			SystemProperty: systemProperty,
			ObjectType:     objectType,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns all validity definitions for a specific system property
func (r *ValidityRead) show(q *msg.Request, mr *msg.Result) {
	var (
		systemProperty, objectType string
		isInherited                bool
		rows                       *sql.Rows
		err                        error
	)

	if rows, err = r.stmtShow.Query(
		q.Validity.SystemProperty,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&systemProperty,
			&objectType,
			&isInherited,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Validity = append(mr.Validity, proto.Validity{
			SystemProperty: systemProperty,
			ObjectType:     objectType,
			Direct:         !isInherited,
			Inherited:      isInherited,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *ValidityRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
