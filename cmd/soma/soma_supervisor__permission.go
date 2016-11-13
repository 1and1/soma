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

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) permission(q *msg.Request) {
	result := msg.FromRequest(q)

	s.reqLog.Printf(LogStrReq,
		q.Type,
		fmt.Sprintf("%s/%s", q.Section,
			q.Action),
		q.User,
		q.RemoteAddr,
	)

	switch q.Action {
	case `list`, `search/name`, `show`:
		go func() { s.permission_read(q) }()

	case `add`, `remove`:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.permission_write(q)

	default:
		result.NotImplemented(fmt.Errorf("Unknown requested action:"+
			" %s/%s/%s", q.Type, q.Section, q.Action))
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *supervisor) permission_write(q *msg.Request) {
	result := msg.FromRequest(q)

	switch q.Action {
	case `add`:
		switch q.Permission.Category {
		case
			`global`, `permission`, `operations`,
			`repository`, `team`, `monitoring`:
			s.permission_add(q, &result)
		}
	case `remove`:
		switch q.Permission.Category {
		case `global`, `permission`, `operations`:
			s.permission_remove_global(q, &result)
		case `repository`:
			s.permission_remove_repository(q, &result)
		case `team`:
			s.permission_remove_team(q, &result)
		case `monitoring`:
			s.permission_remove_monitoring(q, &result)
		}
	}

	q.Reply <- result
}

func (s *supervisor) permission_add(q *msg.Request, r *msg.Result) {
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
		`permission_add_tx_perm`: stmt.PermissionAdd,
		`permission_add_tx_link`: stmt.PermissionLinkGrant,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.PermissionTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.permission_add_tx(q, txMap); err != nil {
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

	r.Permission = []proto.Permission{q.Permission}
}

func (s *supervisor) permission_add_tx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err                        error
		res                        sql.Result
		grantPermId, grantCategory string
	)
	q.Permission.Id = uuid.NewV4().String()
	grantPermId = uuid.NewV4().String()
	switch q.Permission.Category {
	case `global`:
		grantCategory = `global:grant`
	case `permission`:
		grantCategory = `permission:grant`
	case `operations`:
		grantCategory = `operations:grant`
	case `repository`:
		grantCategory = `repository:grant`
	case `team`:
		grantCategory = `team:grant`
	case `monitoring`:
		grantCategory = `monitoring:grant`
	}

	if res, err = txMap[`permission_add_tx_perm`].Exec(
		q.Permission.Id,
		q.Permission.Name,
		q.Permission.Category,
		q.User,
	); err != nil {
		return res, err
	}

	if res, err = txMap[`permission_add_tx_perm`].Exec(
		grantPermId,
		q.Permission.Name,
		grantCategory,
		q.User,
	); err != nil {
		return res, err
	}

	return txMap[`permission_add_tx_link`].Exec(
		grantCategory,
		grantPermId,
		q.Permission.Category,
		q.Permission.Id,
	)
}

func (s *supervisor) permission_remove_global(q *msg.Request,
	r *msg.Result) {
}

func (s *supervisor) permission_remove_repository(q *msg.Request,
	r *msg.Result) {
}

func (s *supervisor) permission_remove_team(q *msg.Request,
	r *msg.Result) {
}

func (s *supervisor) permission_remove_monitoring(q *msg.Request,
	r *msg.Result) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
