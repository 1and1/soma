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
	"github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

// CapabilityWrite handles write requests for capabilities
type CapabilityWrite struct {
	Input                chan msg.Request
	Shutdown             chan struct{}
	conn                 *sql.DB
	stmtAdd              *sql.Stmt
	stmtRemove           *sql.Stmt
	stmtVerifyMetric     *sql.Stmt
	stmtVerifyMonitoring *sql.Stmt
	stmtVerifyView       *sql.Stmt
	appLog               *logrus.Logger
	reqLog               *logrus.Logger
	errLog               *logrus.Logger
}

// newCapabilityWrite return a new CapabilityWrite handler with
// input buffer of length
func newCapabilityWrite(length int) (w *CapabilityWrite) {
	w = &CapabilityWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *CapabilityWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for CapabilityWrite
func (w *CapabilityWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AddCapability:          w.stmtAdd,
		stmt.DelCapability:          w.stmtRemove,
		stmt.MetricVerify:           w.stmtVerifyMetric,
		stmt.VerifyMonitoringSystem: w.stmtVerifyMonitoring,
		stmt.ViewVerify:             w.stmtVerifyView,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`capability`, err, stmt.Name(statement))
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
func (w *CapabilityWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case `add`:
		w.add(q, &result)
	case `remove`:
		w.remove(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// add inserts a new capability
func (w *CapabilityWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		inputVal string
		res      sql.Result
		err      error
	)

	// input validation: MonitoringID
	if w.stmtVerifyMonitoring.QueryRow(
		q.Capability.MonitoringId,
	).Scan(
		&inputVal,
	); err == sql.ErrNoRows {
		mr.NotFound(fmt.Errorf(
			"Monitoring system with ID %s is not registered",
			q.Capability.MonitoringId),
			q.Section,
		)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// input validation: metric
	if w.stmtVerifyMetric.QueryRow(
		q.Capability.Metric,
	).Scan(
		&inputVal,
	); err == sql.ErrNoRows {
		mr.NotFound(fmt.Errorf("Metric %s is not registered",
			q.Capability.Metric),
			q.Section,
		)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// input validation: view
	if w.stmtVerifyView.QueryRow(
		q.Capability.View,
	).Scan(
		&inputVal,
	); err == sql.ErrNoRows {
		mr.NotFound(fmt.Errorf("View %s is not registered",
			q.Capability.View),
			q.Section,
		)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	q.Capability.Id = uuid.NewV4().String()
	if res, err = w.stmtAdd.Exec(
		q.Capability.Id,
		q.Capability.MonitoringId,
		q.Capability.Metric,
		q.Capability.View,
		q.Capability.Thresholds,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Capability = append(mr.Capability, q.Capability)
	}
}

// remove deletes a capability
func (w *CapabilityWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		res sql.Result
		err error
	)

	if res, err = w.stmtRemove.Exec(
		q.Capability.Id,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Capability = append(mr.Capability, q.Capability)
	}
}

// shutdown signals the handler to shut down
func (w *CapabilityWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
