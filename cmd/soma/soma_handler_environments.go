/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß
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
)

type environmentRead struct {
	input    chan msg.Request
	shutdown chan bool
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type environmentWrite struct {
	input      chan msg.Request
	shutdown   chan bool
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
	stmtRename *sql.Stmt
	appLog     *log.Logger
	reqLog     *log.Logger
	errLog     *log.Logger
}

func (r *environmentRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.EnvironmentList: r.stmtList,
		stmt.EnvironmentShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`environment`, err, stmt.Name(statement))
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

func (r *environmentRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case "show":
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

func (r *environmentRead) list(q *msg.Request, mr *msg.Result) {
	var (
		err         error
		rows        *sql.Rows
		environment string
	)
	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&environment); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(`environment`)
			return
		}
		mr.Environment = append(mr.Environment, proto.Environment{
			Name: environment,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

func (r *environmentRead) show(q *msg.Request, mr *msg.Result) {
	var (
		err         error
		environment string
	)
	if err = r.stmtShow.QueryRow(
		q.Environment.Name,
	).Scan(&environment); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	mr.Environment = append(mr.Environment, proto.Environment{
		Name: environment,
	})
	mr.OK()
}

func (r *environmentRead) shutdownNow() {
	r.shutdown <- true
}

func (w *environmentWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.EnvironmentAdd:    w.stmtAdd,
		stmt.EnvironmentRemove: w.stmtRemove,
		stmt.EnvironmentRename: w.stmtRename,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`environment`, err, stmt.Name(statement))
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

func (w *environmentWrite) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(w.reqLog, q)

	switch q.Action {
	case `add`:
		w.add(q, &result)
	case `remove`:
		w.remove(q, &result)
	case `rename`:
		w.rename(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

func (w *environmentWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtAdd.Exec(q.Environment.Name); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Environment = append(mr.Environment, q.Environment)
	}
}

func (w *environmentWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(q.Environment.Name); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Environment = append(mr.Environment, q.Environment)
	}
}

func (w *environmentWrite) rename(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRename.Exec(
		q.Update.Environment.Name,
		q.Environment.Name,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Environment = append(mr.Environment, q.Update.Environment)
	}
}

func (w *environmentWrite) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
