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

func (s *supervisor) section(q *msg.Request) {
	result := msg.FromRequest(q)

	s.reqLog.Printf(LogStrReq,
		q.Type,
		fmt.Sprintf("%s/%s", q.Section, q.Action),
		q.User,
		q.RemoteAddr,
	)

	switch q.Action {
	case `list`, `show`, `search`:
		go func() { s.section_read(q) }()
	case `add`, `remove`:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.section_write(q)
	default:
		result.NotImplemented(fmt.Errorf("Unknown requested action:"+
			" %s/%s/%s", q.Type, q.Section, q.Action))
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *supervisor) section_read(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `list`:
		s.section_list(q, &result)
	case `show`:
		s.section_show(q, &result)
	case `search`:
		s.section_search(q, &result)
	}

	q.Reply <- result
}

func (s *supervisor) section_list(q *msg.Request, r *msg.Result) {
	r.SectionObj = []proto.Section{}
	var (
		err                    error
		rows                   *sql.Rows
		sectionId, sectionName string
	)

	if _, err = s.stmt_SectionList.Query(); err != nil {
		r.ServerError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&sectionId,
			&sectionName,
		); err != nil {
			r.ServerError(err)
			r.Clear(q.Section)
			return
		}
		r.SectionObj = append(r.SectionObj, proto.Section{
			Id:   sectionId,
			Name: sectionName,
		})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		r.Clear(q.Section)
	}
}

func (s *supervisor) section_show(q *msg.Request, r *msg.Result) {
	var (
		err                                    error
		sectionId, sectionName, category, user string
		ts                                     time.Time
	)

	if err = s.stmt_SectionShow.QueryRow(q.SectionObj.Id).Scan(
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
	r.SectionObj = []proto.Section{proto.Section{
		Id:       sectionId,
		Name:     sectionName,
		Category: category,
		Details: &proto.DetailsCreation{
			CreatedAt: ts.Format(rfc3339Milli),
			CreatedBy: user,
		},
	}}
}

func (s *supervisor) section_search(q *msg.Request, r *msg.Result) {
	r.SectionObj = []proto.Section{}
	var (
		err                    error
		rows                   *sql.Rows
		sectionId, sectionName string
	)

	if _, err = s.stmt_SectionSearch.Query(
		q.SectionObj.Name); err != nil {
		r.ServerError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&sectionId,
			&sectionName,
		); err != nil {
			r.ServerError(err)
			r.Clear(q.Section)
			return
		}
		r.SectionObj = append(r.SectionObj, proto.Section{
			Id:   sectionId,
			Name: sectionName,
		})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		r.Clear(q.Section)
	}
}

func (s *supervisor) section_write(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `add`:
		s.section_add(q, &result)
	case `remove`:
		s.section_remove(q, &result)
	}

	q.Reply <- result
}

func (s *supervisor) section_add(q *msg.Request, r *msg.Result) {
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

func (s *supervisor) section_remove(q *msg.Request, r *msg.Result) {
	var (
		err      error
		tx       *sql.Tx
		actionId string
		rows     *sql.Rows
		res      sql.Result
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
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.ActionTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err)
			tx.Rollback()
			return
		}
	}

	// remove all actions in this section
	if rows, err = tx.Query(stmt.SectionListActions,
		q.SectionObj.Id); err != nil {
		r.ServerError(err)
		tx.Rollback()
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&actionId,
		); err != nil {
			r.ServerError(err)
			rows.Close()
			tx.Rollback()
			return
		}
		if res, err = s.action_remove_tx(actionId, txMap); err != nil {
			r.ServerError(err)
			rows.Close()
			tx.Rollback()
			return
		}
		if !r.RowCnt(res.RowsAffected()) {
			rows.Close()
			tx.Rollback()
			return
		}
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		tx.Rollback()
		return
	}

	if res, err = s.section_remove_tx(q.SectionObj.Id,
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

func (s *supervisor) section_remove_tx(id string,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err error
		res sql.Result
	)

	// remove section from all permissions
	if res, err = s.tx_exec(id, `section_tx_removeMap`,
		txMap); err != nil {
		return res, err
	}

	// remove section
	return s.tx_exec(id, `section_tx_remove`, txMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix