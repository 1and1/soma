package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/satori/go.uuid"

)

type somaCapabilityRequest struct {
	action     string
	Capability proto.Capability
	reply      chan somaResult
}

type somaCapabilityResult struct {
	ResultError error
	Capability  proto.Capability
}

func (a *somaCapabilityResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Capabilities = append(r.Capabilities, somaCapabilityResult{ResultError: err})
	}
}

func (a *somaCapabilityResult) SomaAppendResult(r *somaResult) {
	r.Capabilities = append(r.Capabilities, *a)
}

/* Read Access
 */
type somaCapabilityReadHandler struct {
	input     chan somaCapabilityRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaCapabilityReadHandler) run() {
	var err error

	if r.list_stmt, err = r.conn.Prepare(stmt.ListAllCapabilities); err != nil {
		log.Fatal("capability/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_stmt, err = r.conn.Prepare(stmt.ShowCapability); err != nil {
		log.Fatal("capability/show: ", err)
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

func (r *somaCapabilityReadHandler) process(q *somaCapabilityRequest) {
	var (
		id, monitoring, metric, view, monName string
		rows                                  *sql.Rows
		thresholds                            int
		err                                   error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: capability/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(
				&id,
				&monitoring,
				&metric,
				&view,
				&monName,
			)
			result.Append(err, &somaCapabilityResult{
				Capability: proto.Capability{
					Id:           id,
					MonitoringId: monitoring,
					Metric:       metric,
					View:         view,
					Name:         fmt.Sprintf("%s.%s.%s", monName, view, metric),
				},
			})
		}
	case "show":
		log.Printf("R: capability/show for %s", q.Capability.Id)
		if err = r.show_stmt.QueryRow(q.Capability.Id).Scan(
			&id,
			&monitoring,
			&metric,
			&view,
			&thresholds,
			&monName,
		); err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaCapabilityResult{
			Capability: proto.Capability{
				Id:           id,
				MonitoringId: monitoring,
				Metric:       metric,
				View:         view,
				Thresholds:   uint64(thresholds),
				Name:         fmt.Sprintf("%s.%s.%s", monName, view, metric),
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaCapabilityWriteHandler struct {
	input    chan somaCapabilityRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaCapabilityWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.monitoring_capabilities (
	capability_id,
	capability_monitoring,
	capability_metric,
	capability_view,
	threshold_amount)
SELECT $1::uuid, $2::uuid, $3::varchar, $4::varchar, $5::integer
WHERE NOT EXISTS (
	SELECT capability_id
	FROM   soma.monitoring_capabilities
	WHERE  capability_id = $1::uuid
	OR     (    capability_monitoring = $2::uuid
	        AND capability_metric     = $3::varchar
			AND capability_view       = $4::varchar));`)
	if err != nil {
		log.Fatal("capability/add: ", err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.monitoring_capabilities
WHERE  capability_id = $1::uuid;`)
	if err != nil {
		log.Fatal("capability/delete: ", err)
	}
	defer w.del_stmt.Close()

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

func (w *somaCapabilityWriteHandler) process(q *somaCapabilityRequest) {
	var (
		inputVal string
		res      sql.Result
		err      error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: capability/add for %s.%s.%s",
			q.Capability.MonitoringId,
			q.Capability.View,
			q.Capability.Metric,
		)
		// input validation: MonitoringId
		if w.conn.QueryRow("SELECT monitoring_id FROM soma.monitoring_systems WHERE monitoring_id = $1::uuid;",
			q.Capability.MonitoringId).Scan(&inputVal); err == sql.ErrNoRows {
			err = fmt.Errorf("Monitoring system with ID %s is not registered", q.Capability.MonitoringId)
			goto bailout
		} else if err != nil {
			goto bailout
		}

		// input validation: metric
		if w.conn.QueryRow("SELECT metric FROM soma.metrics WHERE metric = $1::varchar;",
			q.Capability.Metric).Scan(&inputVal); err == sql.ErrNoRows {
			err = fmt.Errorf("Metric %s is not registered", q.Capability.Metric)
			goto bailout
		} else if err != nil {
			goto bailout
		}

		// input validation: view
		if w.conn.QueryRow("SELECT view FROM soma.views WHERE view = $1::varchar;",
			q.Capability.View).Scan(&inputVal); err == sql.ErrNoRows {
			err = fmt.Errorf("View %s is not registered", q.Capability.View)
			goto bailout
		} else if err != nil {
			goto bailout
		}

		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Capability.MonitoringId,
			q.Capability.Metric,
			q.Capability.View,
			q.Capability.Thresholds,
		)
		q.Capability.Id = id.String()
	case "delete":
		log.Printf("R: capability/delete for %s", q.Capability.Id)
		res, err = w.del_stmt.Exec(
			q.Capability.Id,
		)
	default:
		log.Printf("R: unimplemented capability/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
bailout:
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaCapabilityResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaCapabilityResult{})
	default:
		result.Append(nil, &somaCapabilityResult{
			Capability: q.Capability,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaCapabilityReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaCapabilityWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
