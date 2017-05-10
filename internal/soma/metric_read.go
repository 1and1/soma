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

// MetricRead handles read requests for metrics
type MetricRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// newMetricRead return a new MetricRead handler with input buffer of
// length
func newMetricRead(length int) (r *MetricRead) {
	r = &MetricRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (r *MetricRead) register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// run is the event loop for MetricRead
func (r *MetricRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.MetricList: r.stmtList,
		stmt.MetricShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`metric`, err, stmt.Name(statement))
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
func (r *MetricRead) process(q *msg.Request) {
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

// list returns all metrics
func (r *MetricRead) list(q *msg.Request, mr *msg.Result) {
	var (
		metric string
		rows   *sql.Rows
		err    error
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&metric); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		mr.Metric = append(mr.Metric, proto.Metric{
			Path: metric,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show return the details of a specific metrics
func (r *MetricRead) show(q *msg.Request, mr *msg.Result) {
	var (
		metric, unit, description string
		err                       error
	)

	if err = r.stmtShow.QueryRow(
		q.Metric.Path,
	).Scan(
		&metric,
		&unit,
		&description,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Metric = append(mr.Metric, proto.Metric{
		Path:        metric,
		Unit:        unit,
		Description: description,
	})
	mr.OK()
}

// shutdown signals the handler to shut down
func (r *MetricRead) shutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
