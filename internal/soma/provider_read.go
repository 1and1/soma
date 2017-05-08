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

// ProviderRead handles read requests for providers
type ProviderRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// newProviderRead return a new ProviderRead handler with input buffer of length
func newProviderRead(length int) (r *ProviderRead) {
	r = &ProviderRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *ProviderRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for ProviderRead
func (r *ProviderRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ProviderList: r.stmtList,
		stmt.ProviderShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`provider`, err, stmt.Name(statement))
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
func (r *ProviderRead) process(q *msg.Request) {
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

// list returns all providers
func (r *ProviderRead) list(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		rows     *sql.Rows
		provider string
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&provider); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Provider = append(mr.Provider, proto.Provider{
			Name: provider,
		})
	}
	mr.OK()
}

// show returns details about a specific provider
func (r *ProviderRead) show(q *msg.Request, mr *msg.Result) {
	var (
		provider string
		err      error
	)

	if err = r.stmtShow.QueryRow(
		q.Provider.Name,
	).Scan(
		&provider,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Provider = append(mr.Provider, proto.Provider{
		Name: provider,
	})
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *ProviderRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
