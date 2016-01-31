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
	Capability somaproto.ProtoCapability
	reply      chan somaResult
}

type somaCapabilityResult struct {
	ResultError error
	Capability  somaproto.ProtoCapability
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

	log.Println("Prepare: capability/list")
	r.list_stmt, err = r.conn.Prepare(`
SELECT    mc.capability_id,
          mc.capability_monitoring,
	      mc.capability_metric,
	      mc.capability_view,
	      ms.monitoring_name
FROM      soma.monitoring_capabilities mc
LEFT JOIN soma.monitoring_systems ms
ON        mc.capability_monitoring = ms.monitoring_id;`)
	if err != nil {
		log.Fatal("capability/list: ", err)
	}
	defer r.list_stmt.Close()

	log.Println("Prepare: capability/show")
	r.show_stmt, err = r.conn.Prepare(`
SELECT    mc.capability_id,
          mc.capability_monitoring,
          mc.capability_metric,
	      mc.capability_view,
	      mc.threshold_amount,
	      ms.monitoring_name
FROM      soma.monitoring_capabilities mc
LEFT JOIN soma.monitoring_systems ms
ON        mc.capability_monitoring = ms.monitoring_id
WHERE     mc.capability_id = $1::uuid;`)
	if err != nil {
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
				Capability: somaproto.ProtoCapability{
					Id:         id,
					Monitoring: monitoring,
					Metric:     metric,
					View:       view,
					Name:       fmt.Sprintf("%s.%s.%s", monName, view, metric),
				},
			})
		}
	case "show":
		log.Printf("R: capability/show for %s", q.Capability.Id)
		err = r.show_stmt.QueryRow(q.Capability.Id).Scan(
			&id,
			&monitoring,
			&metric,
			&view,
			&thresholds,
			&monName,
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

		result.Append(err, &somaCapabilityResult{
			Capability: somaproto.ProtoCapability{
				Id:         id,
				Monitoring: monitoring,
				Metric:     metric,
				View:       view,
				Thresholds: uint64(thresholds),
				Name:       fmt.Sprintf("%s.%s.%s", monName, view, metric),
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

	log.Println("Prepare: capability/add")
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

	log.Println("Prepare: capability/delete")
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
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: capability/add for %s.%s.%s",
			q.Capability.Monitoring,
			q.Capability.View,
			q.Capability.Metric,
		)
		id := uuid.NewV4()
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Capability.Monitoring,
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
