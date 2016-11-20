/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type monitoringRead struct {
	input            chan msg.Request
	shutdown         chan bool
	conn             *sql.DB
	stmtListAll      *sql.Stmt
	stmtListScoped   *sql.Stmt
	stmtShow         *sql.Stmt
	stmtSearchAll    *sql.Stmt
	stmtSearchScoped *sql.Stmt
	appLog           *log.Logger
	reqLog           *log.Logger
	errLog           *log.Logger
}

type monitoringWrite struct {
	input      chan msg.Request
	shutdown   chan bool
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
	appLog     *log.Logger
	reqLog     *log.Logger
	errLog     *log.Logger
}

func (r *monitoringRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ListAllMonitoringSystems:      r.stmtListAll,
		stmt.ListScopedMonitoringSystems:   r.stmtListScoped,
		stmt.ShowMonitoringSystem:          r.stmtShow,
		stmt.SearchAllMonitoringSystems:    r.stmtSearchAll,
		stmt.SearchScopedMonitoringSystems: r.stmtSearchScoped,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`monitoring`, err, stmt.Name(statement))
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

func (r *monitoringRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		switch {
		case q.Flag.Unscoped:
			r.listAll(q, &result)
		default:
			r.listScoped(q, &result)
		}
	case `search`:
		switch {
		case q.Flag.Unscoped:
			r.searchAll(q, &result)
		default:
			r.searchScoped(q, &result)
		}
	case `show`:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

func (r *monitoringRead) listAll(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		monitoringID   string
		monitoringName string
		rows           *sql.Rows
	)
	if rows, err = r.stmtListAll.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&monitoringID,
			&monitoringName,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
			Id:   monitoringID,
			Name: monitoringName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

func (r *monitoringRead) listScoped(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		monitoringID   string
		monitoringName string
		rows           *sql.Rows
	)
	if rows, err = r.stmtListScoped.Query(q.User); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&monitoringID,
			&monitoringName,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
			Id:   monitoringID,
			Name: monitoringName,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

func (r *monitoringRead) show(q *msg.Request, mr *msg.Result) {
	var (
		err                      error
		monitoringID, name, mode string
		contact, teamID          string
		callbackNull             sql.NullString
		callback                 string
	)
	if err = r.stmtShow.QueryRow(q.Monitoring.Id).Scan(
		&monitoringID,
		&name,
		&mode,
		&contact,
		&teamID,
		&callbackNull,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}

	if callbackNull.Valid {
		callback = callbackNull.String
	}
	mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
		Id:       monitoringID,
		Name:     name,
		Mode:     mode,
		Contact:  contact,
		TeamId:   teamID,
		Callback: callback,
	})
	mr.OK()
}

func (r *monitoringRead) searchAll(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		monitoringID   string
		monitoringName string
	)
	// search condition has unique constraint
	if err = r.stmtSearchAll.QueryRow(
		q.Monitoring.Name,
	).Scan(
		&monitoringID,
		&monitoringName,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
		Id:   monitoringID,
		Name: monitoringName,
	})
	mr.OK()
}

func (r *monitoringRead) searchScoped(q *msg.Request, mr *msg.Result) {
	var (
		err            error
		monitoringID   string
		monitoringName string
	)
	// search condition has unique constraint
	if err = r.stmtSearchScoped.QueryRow(
		q.User,
		q.Monitoring.Name,
	).Scan(
		&monitoringID,
		&monitoringName,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	mr.Monitoring = append(mr.Monitoring, proto.Monitoring{
		Id:   monitoringID,
		Name: monitoringName,
	})
	mr.OK()
}

func (r *monitoringRead) shutdownNow() {
	r.shutdown <- true
}

func (w *monitoringWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.MonitoringSystemAdd:    w.stmtAdd,
		stmt.MonitoringSystemRemove: w.stmtRemove,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`monitoring`, err, stmt.Name(statement))
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

func (w *monitoringWrite) process(q *msg.Request) {
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

func (w *monitoringWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		res      sql.Result
		callback sql.NullString
	)

	q.Monitoring.Id = uuid.NewV4().String()
	if q.Monitoring.Callback != `` {
		callback = sql.NullString{
			String: q.Monitoring.Callback,
			Valid:  true,
		}
	}
	if res, err = w.stmtAdd.Exec(
		q.Monitoring.Id,
		q.Monitoring.Name,
		q.Monitoring.Mode,
		q.Monitoring.Contact,
		q.Monitoring.TeamId,
		callback,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Monitoring = append(mr.Monitoring, q.Monitoring)
	}
}

func (w *monitoringWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Monitoring.Id,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Monitoring = append(mr.Monitoring, q.Monitoring)
	}
}

func (w *monitoringWrite) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
