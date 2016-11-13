/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) action(q *msg.Request) {
	result := msg.FromRequest(q)

	s.reqLog.Printf(LogStrReq,
		q.Type,
		fmt.Sprintf("%s/%s", q.Section, q.Action),
		q.User,
		q.RemoteAddr,
	)

	switch q.Action {
	case `list`, `show`, `search`:
		go func() { s.action_read(q) }()
	case `add`, `remove`:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.action_write(q)
	default:
		result.NotImplemented(fmt.Errorf("Unknown requested action:"+
			" %s/%s/%s", q.Type, q.Section, q.Action))
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *supervisor) action_read(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `list`:
		s.action_list(q, &result)
	case `show`:
		s.action_show(q, &result)
	case `search`:
		s.action_search(q, &result)
	}

	q.Reply <- result
}

func (s *supervisor) action_list(q *msg.Request, r *msg.Result) {
	r.ActionObj = []proto.Action{}
	var (
		err                             error
		rows                            *sql.Rows
		actionId, actionName, sectionId string
	)
	if rows, err = s.stmt_ActionList.Query(); err != nil {
		r.ServerError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&actionId,
			&actionName,
			&sectionId,
		); err != nil {
			r.ServerError(err)
			r.Clear(q.Section)
			return
		}
		r.ActionObj = append(r.ActionObj,
			proto.Action{
				Id:        actionId,
				Name:      actionName,
				SectionId: sectionId,
			})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		r.Clear(q.Section)
		return
	}
	r.OK()
}

func (s *supervisor) action_show(q *msg.Request, r *msg.Result) {
	var (
		err                             error
		ts                              time.Time
		actionId, actionName, sectionId string
		category, user, sectionName     string
	)
	if err = s.stmt_ActionShow.QueryRow(q.ActionObj.Id).Scan(
		&actionId,
		&actionName,
		&sectionId,
		&sectionName,
		&category,
		&user,
		&ts,
	); err == sql.ErrNoRows {
		r.NotFound(err)
		return
	} else if err != nil {
		r.ServerError(err)
		return
	}
	r.ActionObj = []proto.Action{proto.Action{
		Id:          actionId,
		Name:        actionName,
		SectionId:   sectionId,
		SectionName: sectionName,
		Category:    category,
		Details: &proto.DetailsCreation{
			CreatedBy: user,
			CreatedAt: ts.Format(rfc3339Milli),
		},
	}}
	r.OK()
}

func (s *supervisor) action_search(q *msg.Request, r *msg.Result) {
	r.ActionObj = []proto.Action{}
	var (
		err                             error
		rows                            *sql.Rows
		actionId, actionName, sectionId string
	)
	if rows, err = s.stmt_ActionList.Query(
		q.ActionObj.Name,
		q.ActionObj.SectionId,
	); err != nil {
		r.ServerError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&actionId,
			&actionName,
			&sectionId,
		); err != nil {
			r.ServerError(err)
			r.Clear(q.Section)
			return
		}
		r.ActionObj = append(r.ActionObj,
			proto.Action{
				Id:        actionId,
				Name:      actionName,
				SectionId: sectionId,
			})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		r.Clear(q.Section)
		return
	}
	r.OK()
}

func (s *supervisor) action_write(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `add`:
		s.action_add(q, &result)
	case `remove`:
		s.action_remove(q, &result)
	}

	q.Reply <- result
}

func (s *supervisor) action_add(q *msg.Request, r *msg.Result) {
	var (
		err error
		res sql.Result
	)
	q.ActionObj.Id = uuid.NewV4().String()
	if res, err = s.stmt_ActionAdd.Exec(
		q.ActionObj.Id,
		q.ActionObj.Name,
		q.ActionObj.SectionId,
		q.User,
	); err != nil {
		r.ServerError(err)
		return
	}
	if r.RowCnt(res.RowsAffected()) {
		r.ActionObj = []proto.Action{q.ActionObj}
	}
}

func (s *supervisor) action_remove(q *msg.Request, r *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		r.ServerError(err)
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`action_tx_remove`:    stmt.ActionRemove,
		`action_tx_removeMap`: stmt.ActionRemoveFromMap,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.ActionTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.action_remove_tx(q.ActionObj.Id,
		txMap); err != nil {
		r.ServerError(err)
		tx.Rollback()
		return
	}
	// sets r.OK()
	if !r.RowCnt(res.RowsAffected()) {
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		r.ServerError(err)
		return
	}

	r.ActionObj = []proto.Action{q.ActionObj}
}

func (s *supervisor) action_remove_tx(id string,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err error
		res sql.Result
	)

	// remove action from all permissions
	if res, err = s.tx_exec(id, `action_tx_removeMap`,
		txMap); err != nil {
		return res, err
	}

	// remove action
	return s.tx_exec(id, `action_tx_remove`, txMap)
}

func (s *supervisor) tx_exec(id, name string,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	return txMap[name].Exec(id)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
