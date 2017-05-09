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

// UserRead handles read requests for users
type UserRead struct {
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

// newUserRead return a new UserRead handler with input buffer of length
func newUserRead(length int) (r *UserRead) {
	r = &UserRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *UserRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for UserRead
func (r *UserRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListUsers: r.stmtList,
		stmt.ShowUsers: r.stmtShow,
		stmt.SyncUsers: r.stmtSync,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`user`, err, stmt.Name(statement))
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
func (r *UserRead) process(q *msg.Request) {
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

// list returns all users
func (r *UserRead) list(q *msg.Request, mr *msg.Result) {
	var (
		userID, userName string
		rows             *sql.Rows
		err              error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&userID,
			&userName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.User = append(mr.User, proto.User{
			Id:       userID,
			UserName: userName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific user
func (r *UserRead) show(q *msg.Request, mr *msg.Result) {
	var (
		userID, userName              string
		firstName, lastName           string
		mailAddr, team                string
		employeeNr                    int
		isActive, isSystem, isDeleted bool
		err                           error
	)

	if err = r.stmtShow.QueryRow(
		q.User.Id,
	).Scan(
		&userID,
		&userName,
		&firstName,
		&lastName,
		&employeeNr,
		&mailAddr,
		&isActive,
		&isSystem,
		&isDeleted,
		&team,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.User = append(mr.User, proto.User{
		Id:             userID,
		UserName:       userName,
		FirstName:      firstName,
		LastName:       lastName,
		EmployeeNumber: strconv.Itoa(employeeNr),
		MailAddress:    mailAddr,
		IsActive:       isActive,
		IsSystem:       isSystem,
		IsDeleted:      isDeleted,
		TeamId:         team,
	})
	mr.OK()
}

// sync returns all user records suitable for sync update calculation
func (r *UserRead) sync(q *msg.Request, mr *msg.Result) {
	var (
		userID, userName    string
		firstName, lastName string
		mailAddr, team      string
		employeeNr          int
		isDeleted           bool
		err                 error
		rows                *sql.Rows
	)

	if rows, err = r.stmtSync.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&userID,
			&userName,
			&firstName,
			&lastName,
			&employeeNr,
			&mailAddr,
			&isDeleted,
			&team,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}

		mr.User = append(mr.User, proto.User{
			Id:             userID,
			UserName:       userName,
			FirstName:      firstName,
			LastName:       lastName,
			EmployeeNumber: strconv.Itoa(employeeNr),
			MailAddress:    mailAddr,
			IsDeleted:      isDeleted,
			TeamId:         team,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *UserRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
