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

type entityRead struct {
	input    chan msg.Request
	shutdown chan bool
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type entityWrite struct {
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

func (r *entityRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ObjectTypeList: r.stmtList,
		stmt.ObjectTypeShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`entity`, err, stmt.Name(statement))
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

func (r *entityRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case `list`:
		r.list(q, &result)
	case `show`:
		r.show(q, &result)
	default:
	}

	q.Reply <- result
}

func (r *entityRead) list(q *msg.Request, mr *msg.Result) {
	var (
		err    error
		rows   *sql.Rows
		entity string
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(&entity); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.Entity = append(mr.Entity, proto.Entity{
			Name: entity,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

func (r *entityRead) show(q *msg.Request, mr *msg.Result) {
	var entity string
	var err error

	if err = r.stmtShow.QueryRow(
		q.Entity.Name,
	).Scan(&entity); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	mr.Entity = append(mr.Entity, proto.Entity{
		Name: entity,
	})
	mr.OK()
}

func (r *entityRead) shutdownNow() {
	r.shutdown <- true
}

func (w *entityWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.ObjectTypeAdd:    w.stmtAdd,
		stmt.ObjectTypeDel:    w.stmtRemove,
		stmt.ObjectTypeRename: w.stmtRename,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`entity`, err, stmt.Name(statement))
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

func (w *entityWrite) process(q *msg.Request) {
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

func (w *entityWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtAdd.Exec(q.Entity.Name); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Entity = append(mr.Entity, q.Entity)
	}
}

func (w *entityWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(q.Entity.Name); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Entity = append(mr.Entity, q.Entity)
	}
}

func (w *entityWrite) rename(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Update.Entity.Name,
		q.Entity.Name,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Entity = append(mr.Entity, q.Update.Entity)
	}
}

func (w *entityWrite) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
