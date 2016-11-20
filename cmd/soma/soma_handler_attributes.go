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
)

type attributeRead struct {
	input    chan msg.Request
	shutdown chan bool
	conn     *sql.DB
	stmtList *sql.Stmt
	stmtShow *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

type attributeWrite struct {
	input      chan msg.Request
	shutdown   chan bool
	conn       *sql.DB
	stmtAdd    *sql.Stmt
	stmtRemove *sql.Stmt
	appLog     *log.Logger
	reqLog     *log.Logger
	errLog     *log.Logger
}

func (r *attributeRead) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AttributeList: r.stmtList,
		stmt.AttributeShow: r.stmtShow,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`attribute`, err, stmt.Name(statement))
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

func (r *attributeRead) process(q *msg.Request) {
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

func (r *attributeRead) list(q *msg.Request, mr *msg.Result) {
	var (
		err                    error
		attribute, cardinality string
		rows                   *sql.Rows
	)

	if rows, err = r.stmtList.Query(); err != nil {
		mr.ServerError(err)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&attribute,
			&cardinality,
		); err != nil {
			rows.Close()
			mr.ServerError(err)
			mr.Clear(q.Section)
			return
		}
		mr.Attribute = append(mr.Attribute, proto.Attribute{
			Name:        attribute,
			Cardinality: cardinality,
		})
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err)
		return
	}
	mr.OK()
}

func (r *attributeRead) show(q *msg.Request, mr *msg.Result) {
	var (
		err                    error
		attribute, cardinality string
	)

	if err = r.stmtShow.QueryRow(q.Attribute.Name).Scan(
		&attribute,
		&cardinality,
	); err == sql.ErrNoRows {
		mr.NotFound(err)
		return
	} else if err != nil {
		mr.ServerError(err)
		return
	}
	mr.Attribute = append(mr.Attribute, proto.Attribute{
		Name:        attribute,
		Cardinality: cardinality,
	})
	mr.OK()
}

func (r *attributeRead) shutdownNow() {
	r.shutdown <- true
}

func (w *attributeWrite) run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.AttributeAdd:    w.stmtAdd,
		stmt.AttributeRemove: w.stmtRemove,
	} {
		if prepStmt, err = w.conn.Prepare(statement); err != nil {
			w.errLog.Fatal(`attribute`, err, stmt.Name(statement))
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

func (w *attributeWrite) process(q *msg.Request) {
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

func (w *attributeWrite) add(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtAdd.Exec(
		q.Attribute.Name,
		q.Attribute.Cardinality,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Attribute = append(mr.Attribute, q.Attribute)
	}
}

func (w *attributeWrite) remove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		res sql.Result
	)

	if res, err = w.stmtRemove.Exec(
		q.Attribute.Name,
	); err != nil {
		mr.ServerError(err)
		return
	}
	if mr.RowCnt(res.RowsAffected()) {
		mr.Attribute = append(mr.Attribute, q.Attribute)
	}
}

func (w *attributeWrite) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
