package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

)

type somaMetricRequest struct {
	action string
	Metric somaproto.ProtoMetric
	reply  chan somaResult
}

type somaMetricResult struct {
	ResultError error
	Metric      somaproto.ProtoMetric
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
}

func (r *somaMetricReadHandler) run() {
	var err error

	log.Println("Prepare: metric/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT metric
FROM   soma.metrics;`)
	if err != nil {
		log.Fatal("metric/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: metric/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT metric,
       metric_unit,
       description
FROM   soma.metrics
WHERE  metric = $1::varchar;`)
	if err != nil {
		log.Fatal("metric/show: ", err)
	}
	defer r.show_stmt.Close()

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
		log.Printf("R: metrics/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&metric)
			result.Append(err, &somaMetricResult{
				Metric: somaproto.ProtoMetric{
					Metric: metric,
				},
			})
		}
	case "show":
		log.Printf("R: metrics/show for %s", q.Metric.Metric)
		err = r.show_stmt.QueryRow(q.Metric.Metric).Scan(
			&metric,
			&unit,
			&description,
		)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaMetricResult{
			Metric: somaproto.ProtoMetric{
				Metric:      metric,
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
}

func (w *somaMetricWriteHandler) run() {
	var err error

	log.Println("Prepare: metric/add")
	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.metrics (
	metric,
	metric_unit,
	description)
SELECT $1::varchar, $2::varchar, $3::text WHERE NOT EXISTS (
	SELECT metric
	FROM   soma.metrics
	WHERE  metric = $1::varchar);`)
	if err != nil {
		log.Fatal("metric/add: ", err)
	}
	defer w.add_stmt.Close()

	log.Println("Prepare: metric/delete")
	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.metrics
WHERE  metric = $1::varchar;`)
	if err != nil {
		log.Fatal("metric/delete: ", err)
	}
	defer w.del_stmt.Close()

	log.Println("Prepare: metric/package-add")
	w.pkg_add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.metric_packages (
	metric,
	metric_provider,
	package)
SELECT $1::varchar, $2::varchar, $3::varchar WHERE NOT EXISTS (
	SELECT metric
	FROM   soma.metric_packages
	WHERE  metric = $1::varchar
	AND    metric_provider = $2::varchar);`)
	if err != nil {
		log.Fatal("metric/package-add")
	}
	defer w.pkg_add_stmt.Close()

	log.Println("Prepare: metric/package-del")
	w.pkg_del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.metric_packages
WHERE  metric = $1::varchar;`)
	if err != nil {
		log.Fatal("metric/package-del")
	}
	defer w.pkg_del_stmt.Close()

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
		res    sql.Result
		err    error
		tx     *sql.Tx
		pkg    somaproto.ProtoMetricProviderPackage
		rowCnt int64
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: metrics/add for %s", q.Metric.Metric)

		// start transaction
		tx, err = w.conn.Begin()
		if err != nil {
			goto bailout
		}
		defer tx.Rollback()

		// insert metric
		res, err = tx.Stmt(w.add_stmt).Exec(
			q.Metric.Metric,
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
	pkgloop:
		for _, pkg = range *q.Metric.Packages {
			res, err = tx.Stmt(w.pkg_add_stmt).Exec(
				q.Metric.Metric,
				pkg.Provider,
				pkg.Package,
			)
			if err != nil {
				break pkgloop
			}
		}
		err = tx.Commit()
	case "delete":
		log.Printf("R: metrics/del for %s", q.Metric.Metric)

		// start transaction
		tx, err = w.conn.Begin()
		if err != nil {
			goto bailout
		}
		defer tx.Rollback()

		// delete provider package information for this metric
		res, err = tx.Stmt(w.pkg_del_stmt).Exec(
			q.Metric.Metric,
		)
		if err != nil {
			goto bailout
		}

		// delete metric that is no longer references
		res, err = tx.Stmt(w.del_stmt).Exec(
			q.Metric.Metric,
		)
		if err != nil {
			goto bailout
		}

		// get row count while still within the transaction
		rowCnt, _ = res.RowsAffected()

		err = tx.Commit()
	default:
		log.Printf("R: unimplemented metrics/%s", q.action)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
