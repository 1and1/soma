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

// TeamWrite handles write requests for views
type TeamWrite struct {
	Input      chan msg.Request
	Shutdown   chan struct{}
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
	stmtUpdate *sql.Stmt
	appLog     *logrus.Logger
	reqLog     *logrus.Logger
	errLog     *logrus.Logger
	soma       *Soma
}

// newTeamWrite return a new TeamWrite handler with input buffer of
// length
func newTeamWrite(length int, s *Soma) (w *TeamWrite) {
	w = &TeamWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	w.soma = s
	return
}

// register initializes resources provided by the Soma app
func (w *TeamWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for TeamWrite
func (w *TeamWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.TeamAdd:    w.stmtAdd,
		stmt.TeamUpdate: w.stmtUpdate,
		stmt.TeamDel:    w.stmtRemove,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`team`, err, stmt.Name(statement))
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
func (w *TeamWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	// supervisor must be notified of user change
	super := w.soma.handlerMap.Get(`supervisor`).(*Supervisor)
	notify := msg.Request{
		Section: `map`,
		Action:  q.Action,
		Super: &msg.Supervisor{
			Object: `team`,
			Team:   q.Team,
		},
	}

	switch q.Action {
	case `add`:
		w.add(q, &result)
	case `remove`:
		w.remove(q, &result)
	case `update`:
		w.update(q, &result)
	default:
		result.UnknownRequest(q)
	}

	// send supervisor notify
	if result.IsOK() {
		super.Input <- notify
	}
	q.Reply <- result
}

// add inserts a new team
func (w *TeamWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	q.Team.Id = uuid.NewV4().String()
	if res, err = w.stmtAdd.Exec(
		q.Team.Id,
		q.Team.Name,
		q.Team.LdapId,
		q.Team.IsSystem,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Team = append(mr.Team, q.Team)
	}
}

// remove deletes a team
func (w *TeamWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Team.Id,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Team = append(mr.Team, q.Team)
	}
}

// update refreshes a team's information
func (w *TeamWrite) update(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtUpdate.Exec(
		q.Team.Name,
		q.Team.LdapId,
		q.Team.IsSystem,
		q.Team.Id,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Team = append(mr.Team, q.Team)
	}
}

// shutdown signals the handler to shut down
func (w *TeamWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
