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
		default:
			result.ServerError(fmt.Errorf(`Illegal category`))
		}
	case `remove`:
		switch q.Permission.Category {
		case
			`global`, `permission`, `operations`,
			`repository`, `team`, `monitoring`:
			s.permission_remove(q, &result)
		default:
			result.ServerError(fmt.Errorf(`Illegal category`))
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

func (s *supervisor) permission_remove(q *msg.Request, r *msg.Result) {
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
		`permission_rm_tx_rev_gl`: stmt.PermissionRevokeGlobal,
		`permission_rm_tx_rev_rp`: stmt.PermissionRevokeRepository,
		`permission_rm_tx_rev_tm`: stmt.PermissionRevokeTeam,
		`permission_rm_tx_rev_mn`: stmt.PermissionRevokeMonitoring,
		`permission_rm_tx_lookup`: stmt.PermissionLookupGrantId,
		`permission_rm_tx_unlink`: stmt.PermissionRemoveLink,
		`permission_rm_tx_remove`: stmt.PermissionRemove,
		`permission_rm_tx_unmapa`: stmt.PermissionUnmapAll,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.PermissionTx.Prepare(%s) error: %s",
				name, err.Error())
			r.ServerError(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.permission_remove_tx(q, txMap); err != nil {
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

func (s *supervisor) permission_remove_tx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err                  error
		res                  sql.Result
		grantingPermissionId string
		revocation           string
	)

	// select correct revocation statement scope
	switch q.Permission.Category {
	case `global`, `permission`, `operations`:
		revocation = `permission_rm_tx_rev_gl`
	case `repository`:
		revocation = `permission_rm_tx_rev_rp`
	case `team`:
		revocation = `permission_rm_tx_rev_tm`
	case `monitoring`:
		revocation = `permission_rm_tx_rev_mn`
	}

	// lookup which permission grants this permission
	if err = txMap[`permission_rm_tx_lookup`].QueryRow(
		q.Permission.Id,
	).Scan(
		&grantingPermissionId,
	); err != nil {
		return res, err
	}

	// revoke all grants of the granting permission
	if res, err = txMap[revocation].Exec(
		grantingPermissionId,
	); err != nil {
		return res, err
	}

	// sever the link between permission and granting permission
	if res, err = txMap[`permission_rm_tx_unlink`].Exec(
		q.Permission.Id,
	); err != nil {
		return res, err
	}

	// remove granting permission
	if res, err = txMap[`permission_rm_tx_remove`].Exec(
		grantingPermissionId,
	); err != nil {
		return res, err
	}

	// revoke all grants of the permission
	if res, err = txMap[revocation].Exec(
		q.Permission.Id,
	); err != nil {
		return res, err
	}

	// unmap all actions from the permission
	if res, err = txMap[`permission_rm_tx_unmapa`].Exec(
		q.Permission.Id,
	); err != nil {
		return res, err
	}

	// remove permission
	return txMap[`permission_rm_tx_remove`].Exec(q.Permission.Id)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
