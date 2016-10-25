package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
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
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaMonitoringReadHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListAllMonitoringSystems:    r.list_stmt,
		stmt.ShowMonitoringSystem:        r.show_stmt,
		stmt.ListScopedMonitoringSystems: r.scli_stmt,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`monitoring`, err, statement)
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
			r.reqLog.Printf("R: monitorings/list")
			rows, err = r.list_stmt.Query()
		} else {
			r.reqLog.Printf("R: monitorings/scoped-list for %s", q.user)
			rows, err = r.scli_stmt.Query(q.user)
		}
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(
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
		r.reqLog.Printf("R: monitoring/show for %s", q.Monitoring.Id)
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
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaMonitoringWriteHandler) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.MonitoringSystemAdd: w.add_stmt,
		stmt.MonitoringSystemDel: w.del_stmt,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`monitoring`, err, statement)
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

func (w *somaMonitoringWriteHandler) process(q *somaMonitoringRequest) {
	var (
		callback sql.NullString
		res      sql.Result
		err      error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		w.reqLog.Printf("R: monitoring/add for %s", q.Monitoring.Name)
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
		w.reqLog.Printf("R: monitoring/delete for %s", q.Monitoring.Id)
		res, err = w.del_stmt.Exec(
			q.Monitoring.Id,
		)
	default:
		w.reqLog.Printf("R: unimplemented monitorings/%s", q.action)
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

/* Ops Access
 */
func (r *somaMonitoringReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaMonitoringWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
