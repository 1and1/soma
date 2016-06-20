package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/satori/go.uuid"

)

type somaMonitoringRequest struct {
	action     string
	admin      bool
	user       string
	Monitoring proto.Monitoring
	reply      chan somaResult
}

type somaMonitoringResult struct {
	ResultError error
	Monitoring  proto.Monitoring
}

func (a *somaMonitoringResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Systems = append(r.Systems, somaMonitoringResult{ResultError: err})
	}
}

func (a *somaMonitoringResult) SomaAppendResult(r *somaResult) {
	r.Systems = append(r.Systems, *a)
}

/* Read Access
 */
type somaMonitoringReadHandler struct {
	input     chan somaMonitoringRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	scli_stmt *sql.Stmt
}

func (r *somaMonitoringReadHandler) run() {
	var err error

	if r.list_stmt, err = r.conn.Prepare(stmt.ListAllMonitoringSystems); err != nil {
		log.Fatal("monitoring/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_stmt, err = r.conn.Prepare(stmt.ShowMonitoringSystem); err != nil {
		log.Fatal("monitoring/show: ", err)
	}
	defer r.show_stmt.Close()

	if r.scli_stmt, err = r.conn.Prepare(stmt.ListScopedMonitoringSystems); err != nil {
		log.Fatal("monitoring/scoped-list: ", err)
	}
	defer r.scli_stmt.Close()

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

func (r *somaMonitoringReadHandler) process(q *somaMonitoringRequest) {
	var (
		id, name, mode, contact, team string
		rows                          *sql.Rows
		callback                      sql.NullString
		callbackString                string
		err                           error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		if q.admin {
			log.Printf("R: monitorings/list")
			rows, err = r.list_stmt.Query()
		} else {
			log.Printf("R: monitorings/scoped-list for %s", q.user)
			rows, err = r.scli_stmt.Query(q.user)
		}
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
				&id,
				&name,
			)
			result.Append(err, &somaMonitoringResult{
				Monitoring: proto.Monitoring{
					Id:   id,
					Name: name,
				},
			})
		}
	case "show":
		log.Printf("R: monitoring/show for %s", q.Monitoring.Id)
		err = r.show_stmt.QueryRow(q.Monitoring.Id).Scan(
			&id,
			&name,
			&mode,
			&contact,
			&team,
			&callback,
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

		if callback.Valid {
			callbackString = callback.String
		} else {
			callbackString = ""
		}
		result.Append(err, &somaMonitoringResult{
			Monitoring: proto.Monitoring{
				Id:       id,
				Name:     name,
				Mode:     mode,
				Contact:  contact,
				TeamId:   team,
				Callback: callbackString,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaMonitoringWriteHandler struct {
	input    chan somaMonitoringRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaMonitoringWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.monitoring_systems (
	monitoring_id,
	monitoring_name,
	monitoring_system_mode,
	monitoring_contact,
	monitoring_owner_team,
	monitoring_callback_uri)
SELECT $1::uuid, $2::varchar, $3::varchar, $4::uuid, $5::uuid, $6::text
WHERE NOT EXISTS (
	SELECT monitoring_id
	FROM   soma.monitoring_systems
	WHERE  monitoring_id = $1::uuid
    OR     monitoring_name = $2::varchar);`)
	if err != nil {
		log.Fatal("monitoring/add: ", err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.monitoring_systems
WHERE  monitoring_id = $1::uuid;`)
	if err != nil {
		log.Fatal("monitoring/delete: ", err)
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

func (w *somaMonitoringWriteHandler) process(q *somaMonitoringRequest) {
	var (
		callback sql.NullString
		res      sql.Result
		err      error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: monitoring/add for %s", q.Monitoring.Name)
		id := uuid.NewV4()
		if q.Monitoring.Callback == "" {
			callback = sql.NullString{
				String: "",
				Valid:  false,
			}
		} else {
			callback = sql.NullString{
				String: q.Monitoring.Callback,
				Valid:  true,
			}
		}
		res, err = w.add_stmt.Exec(
			id.String(),
			q.Monitoring.Name,
			q.Monitoring.Mode,
			q.Monitoring.Contact,
			q.Monitoring.TeamId,
			callback,
		)
		q.Monitoring.Id = id.String()
	case "delete":
		log.Printf("R: monitoring/delete for %s", q.Monitoring.Id)
		res, err = w.del_stmt.Exec(
			q.Monitoring.Id,
		)
	default:
		log.Printf("R: unimplemented monitorings/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaMonitoringResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaMonitoringResult{})
	default:
		result.Append(nil, &somaMonitoringResult{
			Monitoring: q.Monitoring,
		})
	}
	q.reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
