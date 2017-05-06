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
	"github.com/1and1/soma/lib/proto"
	"github.com/Sirupsen/logrus"
)

// TeamRead handles read requests for teams
type TeamRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	stmtSync *sql.Stmt
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// newTeamRead return a new TeamRead handler with input buffer of length
func newTeamRead(length int) (r *TeamRead) {
	r = &TeamRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *TeamRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for TeamRead
func (r *TeamRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListTeams: r.stmtList,
		stmt.ShowTeams: r.stmtShow,
		stmt.SyncTeams: r.stmtSync,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`team`, err, stmt.Name(statement))
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
func (r *TeamRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case `show`:
		r.show(q, &result)
	case `sync`:
		r.sync(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all teams
func (r *TeamRead) list(q *msg.Request, mr *msg.Result) {
	var (
		teamID, teamName string
		rows             *sql.Rows
		err              error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&teamID,
			&teamName,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.Team = append(mr.Team, proto.Team{
			Id:   teamID,
			Name: teamName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

// show returns the details for a specific team
func (r *TeamRead) show(q *msg.Request, mr *msg.Result) {
	var (
		teamID, teamName string
		ldapID           int
		systemFlag       bool
		err              error
	)

	if err = r.stmtShow.QueryRow(
		q.Team.Id,
	).Scan(
		&teamID,
		&teamName,
		&ldapID,
		&systemFlag,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		mr.Clear(q.Section)
		return
	} else if err != nil {
		mr.ServerError(err)
		mr.Clear(q.Section)
		return
	}
	mr.Team = append(mr.Team, proto.Team{
		Id:       teamID,
		Name:     teamName,
		LdapId:   strconv.Itoa(ldapID),
		IsSystem: systemFlag,
	})
	mr.OK()
}

// sync returns all teams in a format suitable for sync processing
func (r *TeamRead) sync(q *msg.Request, mr *msg.Result) {
	var (
		teamID, teamName string
		ldapID           int
		systemFlag       bool
		rows             *sql.Rows
		err              error
	)

	if rows, err = r.stmtSync.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&teamID,
			&teamName,
			&ldapID,
			&systemFlag,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.Team = append(mr.Team, proto.Team{
			Id:       teamID,
			Name:     teamName,
			LdapId:   strconv.Itoa(ldapID),
			IsSystem: systemFlag,
		})
	}
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *TeamRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
