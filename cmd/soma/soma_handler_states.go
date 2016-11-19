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

type stateRead struct {
	input    chan msg.Request
	shutdown chan bool
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type stateWrite struct {
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

func (r *stateRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ObjectStateList: r.stmtList,
		stmt.ObjectStateShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`state`, err, stmt.Name(statement))
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

func (r *stateRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case `show`:
		r.show(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

func (r *stateRead) list(q *msg.Request, mr *msg.Result) {
	var (
		err   error
		rows  *sql.Rows
		state string
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&state); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.State = append(mr.State, proto.State{
			Name: state,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

func (r *stateRead) show(q *msg.Request, mr *msg.Result) {
	var state string
	var err error

	if err = r.stmtShow.QueryRow(
		q.State.Name,
	).Scan(&state); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	mr.State = append(mr.State, proto.State{
		Name: state,
	})
	mr.OK()
}

func (r *stateRead) shutdownNow() {
	r.shutdown <- true
}

func (w *stateWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ObjectStateAdd:    w.stmtAdd,
		stmt.ObjectStateRemove: w.stmtRemove,
		stmt.ObjectStateRename: w.stmtRename,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`state`, err, stmt.Name(statement))
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

func (w *stateWrite) process(q *msg.Request) {
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

func (w *stateWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtAdd.Exec(q.State.Name); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.State = append(mr.State, q.State)
	}
}

func (w *stateWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(q.State.Name); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.State = append(mr.State, q.State)
	}
}

func (w *stateWrite) rename(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Update.State.Name,
		q.State.Name,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.State = append(mr.State, q.Update.State)
	}
}

func (w *stateWrite) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
