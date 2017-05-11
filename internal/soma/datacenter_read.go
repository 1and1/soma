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
	"github.com/1and1/soma/lib/proto"
	"github.com/Sirupsen/logrus"
)

// DatacenterRead handles read requests for datacenters
type DatacenterRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// newDatacenterRead return a new DatacenterRead handler with input
// buffer of length
func newDatacenterRead(length int) (r *DatacenterRead) {
	r = &DatacenterRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *DatacenterRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for DatacenterRead
func (r *DatacenterRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.DatacenterList: r.stmtList,
		stmt.DatacenterShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`datacenter`, err, stmt.Name(statement))
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
func (r *DatacenterRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`, `sync`:
		r.list(q, &result)
	case `show`:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all datacenters
func (r *DatacenterRead) list(q *msg.Request, mr *msg.Result) {
	var (
		datacenter string
		rows       *sql.Rows
		err        error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&datacenter); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Datacenter = append(mr.Datacenter, proto.Datacenter{
			Locode: datacenter,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns details about a specific datacenter
func (r *DatacenterRead) show(q *msg.Request, mr *msg.Result) {
	var (
		datacenter string
		err        error
	)

	if err = r.stmtShow.QueryRow(
		q.Datacenter.Locode,
	).Scan(
		&datacenter,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Datacenter = append(mr.Datacenter, proto.Datacenter{
		Locode: datacenter,
	})
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *DatacenterRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
