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

// MetricWrite handles write requests for metrics
type MetricWrite struct {
	Input              chan msg.Request
	Shutdown           chan struct{}
	conn               *sql.DB
	stmtAdd            *sql.Stmt
	stmtPkgAdd         *sql.Stmt
	stmtPkgRemove      *sql.Stmt
	stmtRemove         *sql.Stmt
	stmtVerifyProvider *sql.Stmt
	stmtVerifyUnit     *sql.Stmt
	appLog             *logrus.Logger
	reqLog             *logrus.Logger
	errLog             *logrus.Logger
}

// newMetricWrite return a new MetricWrite handler with input buffer of
// length
func newMetricWrite(length int) (w *MetricWrite) {
	w = &MetricWrite{}
	w.Input = make(chan msg.Request, length)
	w.Shutdown = make(chan struct{})
	return
}

// register initializes resources provided by the Soma app
func (w *MetricWrite) register(c *sql.DB, l ...*logrus.Logger) {
	w.conn = c
	w.appLog = l[0]
	w.reqLog = l[1]
	w.errLog = l[2]
}

// run is the event loop for MetricWrite
func (w *MetricWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.MetricAdd:      w.stmtAdd,
		stmt.MetricDel:      w.stmtRemove,
		stmt.MetricPkgAdd:   w.stmtPkgAdd,
		stmt.MetricPkgDel:   w.stmtPkgRemove,
		stmt.ProviderVerify: w.stmtVerifyProvider,
		stmt.UnitVerify:     w.stmtVerifyUnit,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`metric`, err, stmt.Name(statement))
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
func (w *MetricWrite) process(q *msg.Request) {
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

// add inserts a new metric
func (w *MetricWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		res      sql.Result
		err      error
		tx       *sql.Tx
		pkg      proto.MetricPackage
		rowCnt   int64
		inputVal string
	)

	// test the referenced unit exists
	if err = w.stmtVerifyUnit.QueryRow(
		q.Metric.Unit,
	).Scan(
		&inputVal,
	); err == sql.ErrNoRows {
		mr.BadRequest(
			fmt.Errorf("Unit %s is not registered",
				q.Metric.Unit),
			q.Section,
		)
		return
	} else if err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// test the referenced providers exist
	if q.Metric.Packages != nil && *q.Metric.Packages != nil {
		for _, pkg = range *q.Metric.Packages {
			if w.stmtVerifyProvider.QueryRow(
				pkg.Provider,
			).Scan(
				&inputVal,
			); err == sql.ErrNoRows {
				mr.BadRequest(
					fmt.Errorf(
						"Provider %s is not registered",
						pkg.Provider),
					q.Section,
				)
				return
			} else if err != nil {
				mr.ServerError(err, q.Section)
				return
			}
		}
	}

	// start transaction
	if tx, err = w.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	defer tx.Rollback()

	// insert metric
	if res, err = tx.Stmt(w.stmtAdd).Exec(
		q.Metric.Path,
		q.Metric.Unit,
		q.Metric.Description,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}

	// get row count while still within the transaction
	if rowCnt, _ = res.RowsAffected(); rowCnt != 1 {
		tx.Rollback()
		mr.ServerError(
			fmt.Errorf("Metric insertion affected %d"+
				" rows instead of 1", rowCnt),
			q.Section)
		return
	}

	// insert all provider package information
	if q.Metric.Packages != nil && *q.Metric.Packages != nil {
		for _, pkg = range *q.Metric.Packages {
			if res, err = tx.Stmt(w.stmtPkgAdd).Exec(
				q.Metric.Path,
				pkg.Provider,
				pkg.Name,
			); err != nil {
				tx.Rollback()
				mr.ServerError(err, q.Section)
				return
			}
			if rowCnt, _ = res.RowsAffected(); rowCnt != 1 {
				tx.Rollback()
				mr.ServerError(
					fmt.Errorf("Package insertion affected %d"+
						" rows instead of 1", rowCnt),
					q.Section)
				return
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Metric = append(mr.Metric, q.Metric)
	mr.OK()
}

// remove deletes a metric
func (w *MetricWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		res    sql.Result
		err    error
		tx     *sql.Tx
		rowCnt int64
	)

	// start transaction
	if tx, err = w.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	// delete provider package information for this metric
	if res, err = tx.Stmt(w.stmtPkgRemove).Exec(
		q.Metric.Path,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}

	// delete metric
	if res, err = tx.Stmt(w.stmtRemove).Exec(
		q.Metric.Path,
	); err != nil {
		tx.Rollback()
		mr.ServerError(err, q.Section)
		return
	}

	// get row count while still within the transaction
	if rowCnt, _ = res.RowsAffected(); rowCnt != 1 {
		tx.Rollback()
		mr.ServerError(
			fmt.Errorf("Metric deletion affected %d"+
				" rows instead of 1", rowCnt),
			q.Section)
		return
	}

	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.Metric = append(mr.Metric, q.Metric)
	mr.OK()
}

// shutdown signals the handler to shut down
func (w *MetricWrite) shutdownNow() {
	close(w.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
