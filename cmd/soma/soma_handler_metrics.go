package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

type somaMetricRequest struct {
	action string
	Metric proto.Metric
	reply  chan somaResult
}

type somaMetricResult struct {
	ResultError error
	Metric      proto.Metric
}

func (a *somaMetricResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Metrics = append(r.Metrics, somaMetricResult{ResultError: err})
	}
}

func (a *somaMetricResult) SomaAppendResult(r *somaResult) {
	r.Metrics = append(r.Metrics, *a)
}

/* Read Access
 */
type somaMetricReadHandler struct {
	input     chan somaMetricRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaMetricReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.MetricList: r.list_stmt,
		stmt.MetricShow: r.show_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`metric`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-r.shutdown:
			break runloop
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaMetricReadHandler) process(q *somaMetricRequest) {
	var (
		metric, unit, description string
		rows                      *sql.Rows
		err                       error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		r.reqLog.Printf("R: metrics/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&metric)
			result.Append(err, &somaMetricResult{
				Metric: proto.Metric{
					Path: metric,
				},
			})
		}
	case "show":
		r.reqLog.Printf("R: metrics/show for %s", q.Metric.Path)
		err = r.show_stmt.QueryRow(q.Metric.Path).Scan(
			&metric,
			&unit,
			&description,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaMetricResult{
			Metric: proto.Metric{
				Path:        metric,
				Unit:        unit,
				Description: description,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaMetricWriteHandler struct {
	input        chan somaMetricRequest
	shutdown     chan bool
	conn         *sql.DB
	add_stmt     *sql.Stmt
	del_stmt     *sql.Stmt
	pkg_add_stmt *sql.Stmt
	pkg_del_stmt *sql.Stmt
	appLog       *log.Logger
	reqLog       *log.Logger
	errLog       *log.Logger
}

func (w *somaMetricWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.MetricAdd:    w.add_stmt,
		stmt.MetricDel:    w.del_stmt,
		stmt.MetricPkgAdd: w.pkg_add_stmt,
		stmt.MetricPkgDel: w.pkg_del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`metric`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-w.shutdown:
			break runloop
		case req := <-w.input:
			w.process(&req)
		}
	}
}

func (w *somaMetricWriteHandler) process(q *somaMetricRequest) {
	var (
		res      sql.Result
		err      error
		tx       *sql.Tx
		pkg      proto.MetricPackage
		rowCnt   int64
		inputVal string
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: metrics/add for %s", q.Metric.Path)
		// test the referenced unit exists, to prettify the error
		w.conn.QueryRow(stmt.UnitVerify, q.Metric.Unit).Scan(&inputVal)
		if err == sql.ErrNoRows {
			err = fmt.Errorf("Unit %s is not registered",
				q.Metric.Unit)
			goto bailout
		} else if err != nil {
			goto bailout
		}

		// test the referenced providers exist
		if q.Metric.Packages != nil && *q.Metric.Packages != nil {
			for _, pkg = range *q.Metric.Packages {
				w.conn.QueryRow(
					stmt.ProviderVerify,
					pkg.Provider).Scan(&inputVal)
				if err == sql.ErrNoRows {
					err = fmt.Errorf(
						"Provider %s is not registered",
						pkg.Provider)
					goto bailout
				} else if err != nil {
					goto bailout
				}
			}
		}

		// start transaction
		tx, err = w.conn.Begin()
		if err != nil {
			goto bailout
		}
		defer tx.Rollback()

		// insert metric
		res, err = tx.Stmt(w.add_stmt).Exec(
			q.Metric.Path,
			q.Metric.Unit,
			q.Metric.Description,
		)
		if err != nil {
			goto bailout
		}

		// get row count while still within the transaction
		rowCnt, _ = res.RowsAffected()
		if rowCnt == 0 {
			goto bailout
		}

		// insert all provider package information
		if q.Metric.Packages != nil && *q.Metric.Packages != nil {
		pkgloop:
			for _, pkg = range *q.Metric.Packages {
				res, err = tx.Stmt(w.pkg_add_stmt).Exec(
					q.Metric.Path,
					pkg.Provider,
					pkg.Name,
				)
				if err != nil {
					break pkgloop
				}
			}
		}
		err = tx.Commit()
	case "delete":
		w.reqLog.Printf("R: metrics/del for %s", q.Metric.Path)

		// start transaction
		tx, err = w.conn.Begin()
		if err != nil {
			goto bailout
		}
		defer tx.Rollback()

		// delete provider package information for this metric
		res, err = tx.Stmt(w.pkg_del_stmt).Exec(
			q.Metric.Path,
		)
		if err != nil {
			goto bailout
		}

		// delete metric that is no longer references
		res, err = tx.Stmt(w.del_stmt).Exec(
			q.Metric.Path,
		)
		if err != nil {
			goto bailout
		}

		// get row count while still within the transaction
		rowCnt, _ = res.RowsAffected()

		err = tx.Commit()
	default:
		w.reqLog.Printf("R: unimplemented metrics/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
bailout:
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaMetricResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaMetricResult{})
	default:
		result.Append(nil, &somaMetricResult{
			Metric: q.Metric,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaMetricReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaMetricWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
