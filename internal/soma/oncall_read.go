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

// OncallRead handles read requests for oncall
type OncallRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// newOncallRead return a new OncallRead handler with input buffer of length
func newOncallRead(length int) (r *OncallRead) {
	r = &OncallRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *OncallRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for OncallRead
func (r *OncallRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.OncallList: r.stmtList,
		stmt.OncallShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`oncall`, err, stmt.Name(statement))
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
func (r *OncallRead) process(q *msg.Request) {
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

// list returns all oncall duties
func (r *OncallRead) list(q *msg.Request, mr *msg.Result) {
	var (
		oncallID, oncallName string
		rows                 *sql.Rows
		err                  error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&oncallID, &oncallName); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Oncall = append(mr.Oncall, proto.Oncall{
			Id:   oncallID,
			Name: oncallName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific oncall duty
func (r *OncallRead) show(q *msg.Request, mr *msg.Result) {
	var (
		oncallID, oncallName string
		oncallNumber         int
		err                  error
	)

	if err = r.stmtShow.QueryRow(
		q.Oncall.Id,
	).Scan(
		&oncallID,
		&oncallName,
		&oncallNumber,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Oncall = append(mr.Oncall, proto.Oncall{
		Id:     oncallID,
		Name:   oncallName,
		Number: strconv.Itoa(oncallNumber),
	})
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *OncallRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
