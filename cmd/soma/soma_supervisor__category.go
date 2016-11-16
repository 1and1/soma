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
	"fmt"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) category(q *msg.Request) {
	result := msg.FromRequest(q)

	s.requestLog(q)

	switch q.Action {
	case `list`, `show`:
		go func() { s.category_read(q) }()

	case `add`, `remove`:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.category_write(q)

	default:
		result.UnknownRequest(q)
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *supervisor) category_read(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `list`:
		s.category_list(q, &result)
	case `show`:
		s.category_show(q, &result)
	}
	q.Reply <- result
}

func (s *supervisor) category_write(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `add`:
		s.category_add(q, &result)
	case `remove`:
		//TODO s.category_remove(q, &result)
		s.permission_category_write(q)
		return
	}

	q.Reply <- result
}

func (s *supervisor) category_list(q *msg.Request, r *msg.Result) {
	var (
		err      error
		rows     *sql.Rows
		category string
	)
	if rows, err = s.stmt_ListCategory.Query(); err != nil {
		r.ServerError(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&category,
		); err != nil {
			r.ServerError(err)
			r.Clear(q.Section)
			return
		}
		r.Category = append(r.Category,
			proto.Category{Name: category})
	}
	if err = rows.Err(); err != nil {
		r.ServerError(err)
		r.Clear(q.Section)
	}
	r.OK()
}

func (s *supervisor) category_show(q *msg.Request, r *msg.Result) {
	var (
		err            error
		category, user string
		ts             time.Time
	)
	if err = s.stmt_ShowCategory.QueryRow(q.Category.Name).Scan(
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
	r.Category = []proto.Category{proto.Category{
		Name: category,
		Details: &proto.CategoryDetails{
			CreatedAt: ts.Format(rfc3339Milli),
			CreatedBy: user,
		},
	}}
	r.OK()
}

func (s *supervisor) category_add(q *msg.Request, r *msg.Result) {
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
		`category_add_tx_cat`:  stmt.CategoryAdd,
		`category_add_tx_perm`: stmt.PermissionAdd,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.CategoryTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.category_add_tx(q, txMap); err != nil {
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

	r.Category = []proto.Category{q.Category}
}

func (s *supervisor) category_add_tx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err    error
		res    sql.Result
		permId string
	)

	// create requested category
	if res, err = txMap[`category_add_tx_cat`].Exec(
		q.Category.Name,
		q.User,
	); err != nil {
		return res, err
	}

	// create grant category for requested category
	if res, err = txMap[`category_add_tx_cat`].Exec(
		fmt.Sprintf("%s:grant", q.Category.Name),
		q.User,
	); err != nil {
		return res, err
	}

	// create system permission for category, the category
	// name becomes the permission name in system
	permId = uuid.NewV4().String()
	return txMap[`category_add_tx_perm`].Exec(
		permId,
		q.Category.Name,
		`system`,
		q.User,
	)
}

func (s *supervisor) category_remove(q *msg.Request, r *msg.Result) {
}

func (s *supervisor) category_remove_tx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err error
		res sql.Result
	)
	return res, err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
