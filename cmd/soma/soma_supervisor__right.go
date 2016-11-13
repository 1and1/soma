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

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) right(q *msg.Request) {
	result := msg.FromRequest(q)

	s.reqLog.Printf(LogStrReq,
		q.Type,
		fmt.Sprintf("%s/%s", q.Section, q.Action),
		q.User,
		q.RemoteAddr,
	)

	if q.Grant.RecipientType != `user` {
		result.NotImplemented(fmt.Errorf("Rights for recipient type"+
			" %s are currently not implemented",
			q.Grant.RecipientType))
		goto abort
	}

	switch q.Action {
	case `grant`, `revoke`:
		if s.readonly {
			result.Conflict(fmt.Errorf(`Readonly instance`))
			goto abort
		}
		s.right_write(q)
	default:
		result.NotImplemented(fmt.Errorf("Unknown requested action:"+
			" %s/%s/%s", q.Type, q.Section, q.Action))
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *supervisor) right_write(q *msg.Request) {
	switch q.Action {
	case `grant`:
		switch q.Grant.Category {
		case `system`,
			`global`,
			`global:grant`,
			`permission`,
			`permission:grant`,
			`operations`,
			`operations:grant`:
			s.right_grant_global(q)
		case `repository`,
			`repository:grant`:
			s.right_grant_repository(q)
		case `team`,
			`team:grant`:
			s.right_grant_team(q)
		case `monitoring`,
			`monitoring:grant`:
			s.right_grant_monitoring(q)
		}
	case `revoke`:
		switch q.Grant.Category {
		case `system`,
			`global`,
			`global:grant`,
			`permission`,
			`permission:grant`,
			`operations`,
			`operations:grant`:
			s.right_revoke_global(q)
		case `repository`,
			`repository:grant`:
			s.right_revoke_repository(q)
		case `team`,
			`team:grant`:
			s.right_revoke_team(q)
		case `monitoring`,
			`monitoring:grant`:
			s.right_revoke_monitoring(q)
		}
	}
}

func (s *supervisor) right_grant_global(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                             error
		res                             sql.Result
		adminId, userId, toolId, teamId sql.NullString
	)

	if q.Grant.ObjectType != `` || q.Grant.ObjectId != `` {
		result.BadRequest(fmt.Errorf(
			`Invalid global grant specification`))
		goto dispatch
	}

	switch q.Grant.RecipientType {
	case `admin`:
		adminId.String = q.Grant.RecipientId
		adminId.Valid = true
	case `user`:
		userId.String = q.Grant.RecipientId
		userId.Valid = true
	case `tool`:
		toolId.String = q.Grant.RecipientId
		toolId.Valid = true
	case `team`:
		teamId.String = q.Grant.RecipientId
		teamId.Valid = true
	}

	q.Grant.Id = uuid.NewV4().String()
	if res, err = s.stmt_GrantGlobal.Exec(
		q.Grant.Id,
		adminId,
		userId,
		toolId,
		teamId,
		q.Grant.PermissionId,
		q.Grant.Category,
		q.User,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_grant_repository(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                       error
		res                       sql.Result
		userId, toolId, teamId    sql.NullString
		repoId, bucketId, groupId sql.NullString
		clusterId, nodeId         sql.NullString
		repoName                  string
	)

	switch q.Grant.ObjectType {
	case `repository`:
		repoId.String = q.Grant.ObjectId
		repoId.Valid = true
	case `bucket`:
		if err = s.conn.QueryRow(
			stmt.RepoByBucketId,
			q.Grant.ObjectId,
		).Scan(
			repoId,
			repoName,
		); err != nil {
			result.ServerError(err)
			goto dispatch
		}

		bucketId.String = q.Grant.ObjectId
		bucketId.Valid = true
	case `group`, `cluster`, `node`:
		result.NotImplemented(fmt.Errorf(
			`Deep repository grants currently not implemented.`))
		goto dispatch
	default:
		result.BadRequest(fmt.Errorf(
			`Invalid repository grant specification`))
		goto dispatch
	}

	switch q.Grant.RecipientType {
	case `user`:
		userId.String = q.Grant.RecipientId
		userId.Valid = true
	case `tool`:
		toolId.String = q.Grant.RecipientId
		toolId.Valid = true
	case `team`:
		teamId.String = q.Grant.RecipientId
		teamId.Valid = true
	}

	q.Grant.Id = uuid.NewV4().String()
	if res, err = s.stmt_GrantRepo.Exec(
		q.Grant.Id,
		userId,
		toolId,
		teamId,
		q.Grant.Category,
		q.Grant.PermissionId,
		q.Grant.ObjectType,
		repoId,
		bucketId,
		groupId,
		clusterId,
		nodeId,
		q.User,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_grant_team(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                    error
		res                    sql.Result
		userId, toolId, teamId sql.NullString
	)

	switch q.Grant.RecipientType {
	case `user`:
		userId.String = q.Grant.RecipientId
		userId.Valid = true
	case `tool`:
		toolId.String = q.Grant.RecipientId
		toolId.Valid = true
	case `team`:
		teamId.String = q.Grant.RecipientId
		teamId.Valid = true
	}

	q.Grant.Id = uuid.NewV4().String()
	if res, err = s.stmt_GrantTeam.Exec(
		q.Grant.Id,
		userId,
		toolId,
		teamId,
		q.Grant.Category,
		q.Grant.PermissionId,
		q.Grant.ObjectId,
		q.User,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_grant_monitoring(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		err                    error
		res                    sql.Result
		userId, toolId, teamId sql.NullString
	)

	switch q.Grant.RecipientType {
	case `user`:
		userId.String = q.Grant.RecipientId
		userId.Valid = true
	case `tool`:
		toolId.String = q.Grant.RecipientId
		toolId.Valid = true
	case `team`:
		teamId.String = q.Grant.RecipientId
		teamId.Valid = true
	}

	q.Grant.Id = uuid.NewV4().String()
	if res, err = s.stmt_GrantMonitor.Exec(
		q.Grant.Id,
		userId,
		toolId,
		teamId,
		q.Grant.Category,
		q.Grant.PermissionId,
		q.Grant.ObjectId,
		q.User,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_revoke_global(q *msg.Request) {
	result := msg.FromRequest(q)
	var err error
	var res sql.Result

	if res, err = s.stmt_RevokeGlobal.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_revoke_repository(q *msg.Request) {
	result := msg.FromRequest(q)
	var err error
	var res sql.Result

	if res, err = s.stmt_RevokeRepo.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_revoke_team(q *msg.Request) {
	result := msg.FromRequest(q)
	var err error
	var res sql.Result

	if res, err = s.stmt_RevokeTeam.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_revoke_monitoring(q *msg.Request) {
	result := msg.FromRequest(q)
	var err error
	var res sql.Result

	if res, err = s.stmt_RevokeMonitor.Exec(
		q.Grant.Id,
		q.Grant.PermissionId,
		q.Grant.Category,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
	}

dispatch:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
