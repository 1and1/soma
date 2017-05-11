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
	"fmt"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	"github.com/Sirupsen/logrus"
)

// CapabilityRead handles read requests for capabilities
type CapabilityRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// newCapabilityRead return a new CapabilityRead handler with input buffer of length
func newCapabilityRead(length int) (r *CapabilityRead) {
	r = &CapabilityRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *CapabilityRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for CapabilityRead
func (r *CapabilityRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListAllCapabilities: r.stmtList,
		stmt.ShowCapability:      r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`capability`, err, stmt.Name(statement))
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
func (r *CapabilityRead) process(q *msg.Request) {
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

// list returns all capabilities
func (r *CapabilityRead) list(q *msg.Request, mr *msg.Result) {
	var (
		id, monitoring, metric, view, monName string
		rows                                  *sql.Rows
		err                                   error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&id,
			&monitoring,
			&metric,
			&view,
			&monName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Capability = append(mr.Capability, proto.Capability{
			Id:           id,
			MonitoringId: monitoring,
			Metric:       metric,
			View:         view,
			Name: fmt.Sprintf("%s.%s.%s", monName, view,
				metric),
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns the details of a specific capability
func (r *CapabilityRead) show(q *msg.Request, mr *msg.Result) {
	var (
		id, monitoring, metric, view, monName string
		thresholds                            int
		err                                   error
	)

	if err = r.stmtShow.QueryRow(
		q.Capability.Id,
	).Scan(
		&id,
		&monitoring,
		&metric,
		&view,
		&thresholds,
		&monName,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	mr.Capability = append(mr.Capability, proto.Capability{
		Id:           id,
		MonitoringId: monitoring,
		Metric:       metric,
		View:         view,
		Thresholds:   uint64(thresholds),
		Name:         fmt.Sprintf("%s.%s.%s", monName, view, metric),
	})
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *CapabilityRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
