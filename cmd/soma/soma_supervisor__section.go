/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
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

func (s *supervisor) section(q *msg.Request) {
	result := msg.FromRequest(q)

	s.requestLog(q)

	switch q.Action {
	case `list`, `show`, `search`:
		go func() { s.sectionRead(q) }()
	case `add`, `remove`:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.sectionWrite(q)
	default:
		result.UnknownRequest(q)
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *supervisor) sectionRead(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `list`:
		s.sectionList(q, &result)
	case `show`:
		s.sectionShow(q, &result)
	case `search`:
		s.sectionSearch(q, &result)
	}

	q.Reply <- result
}

func (s *supervisor) sectionList(q *msg.Request, r *msg.Result) {
	r.SectionObj = []proto.Section{}
	var (
		err                    error
		rows                   *sql.Rows
		sectionID, sectionName string
	)

	if _, err = s.stmt_SectionList.Query(); err != nil {
		r.ServerError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&sectionID,
			&sectionName,
		); err != nil {
			r.ServerError(err)
			r.Clear(q.Section)
			return
		}
		r.SectionObj = append(r.SectionObj, proto.Section{
			Id:   sectionID,
			Name: sectionName,
		})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		r.Clear(q.Section)
	}
}

func (s *supervisor) sectionShow(q *msg.Request, r *msg.Result) {
	var (
		err                                    error
		sectionID, sectionName, category, user string
		ts                                     time.Time
	)

	if err = s.stmt_SectionShow.QueryRow(q.SectionObj.Id).Scan(
		&sectionID,
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
	r.SectionObj = []proto.Section{proto.Section{
		Id:       sectionID,
		Name:     sectionName,
		Category: category,
		Details: &proto.DetailsCreation{
			CreatedAt: ts.Format(rfc3339Milli),
			CreatedBy: user,
		},
	}}
}

func (s *supervisor) sectionSearch(q *msg.Request, r *msg.Result) {
	r.SectionObj = []proto.Section{}
	var (
		err                    error
		rows                   *sql.Rows
		sectionID, sectionName string
	)

	if _, err = s.stmt_SectionSearch.Query(
		q.SectionObj.Name); err != nil {
		r.ServerError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&sectionID,
			&sectionName,
		); err != nil {
			r.ServerError(err)
			r.Clear(q.Section)
			return
		}
		r.SectionObj = append(r.SectionObj, proto.Section{
			Id:   sectionID,
			Name: sectionName,
		})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		r.Clear(q.Section)
	}
}

func (s *supervisor) sectionWrite(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `add`:
		s.sectionAdd(q, &result)
	case `remove`:
		s.sectionRemove(q, &result)
	}

	if result.IsOK() {
		handlerMap[`supervisor`].(*supervisor).input <- msg.CacheUpdateFromRequest(q)
	}

	q.Reply <- result
}

func (s *supervisor) sectionAdd(q *msg.Request, r *msg.Result) {
	var (
		err error
		res sql.Result
	)
	q.SectionObj.Id = uuid.NewV4().String()
	if res, err = s.stmt_SectionAdd.Exec(
		q.SectionObj.Id,
		q.SectionObj.Name,
		q.SectionObj.Category,
		q.User,
	); err != nil {
		r.ServerError(err)
		return
	}
	if r.RowCnt(res.RowsAffected()) {
		r.SectionObj = []proto.Section{q.SectionObj}
	}
}

func (s *supervisor) sectionRemove(q *msg.Request, r *msg.Result) {
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
		`action_tx_remove`:     stmt.ActionRemove,
		`action_tx_removeMap`:  stmt.ActionRemoveFromMap,
		`section_tx_remove`:    stmt.SectionRemove,
		`section_tx_removeMap`: stmt.SectionRemoveFromMap,
		`section_tx_actlist`:   stmt.SectionListActions,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.SectionTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.sectionRemoveTx(q.SectionObj.Id,
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

func (s *supervisor) sectionRemoveTx(id string,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err      error
		res      sql.Result
		rows     *sql.Rows
		actionID string
		affected int64
	)

	// remove all actions in this section
	if rows, err = txMap[`section_tx_actlist`].Query(
		id); err != nil {
		return res, err
	}

	for rows.Next() {
		if err = rows.Scan(
			&actionID,
		); err != nil {
			rows.Close()
			return res, err
		}
		if res, err = s.actionRemoveTx(actionID, txMap); err != nil {
			rows.Close()
			return res, err
		}
		if affected, err = res.RowsAffected(); err != nil {
			rows.Close()
			return res, err
		} else if affected != 1 {
			rows.Close()
			return res, fmt.Errorf("Delete statement caught %d rows"+
				" of actions instead of 1 (actionID=%s)", affected,
				actionID)
		}
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	// remove section from all permissions
	if res, err = txMap[`section_tx_removeMap`].Exec(id); err != nil {
		return res, err
	}

	// remove section
	return txMap[`section_tx_remove`].Exec(id)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
