package main

import (
	"database/sql"
	"fmt"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) right(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `right`, Super: &msg.Supervisor{Action: q.Super.Action}}

	s.reqLog.Printf(LogStrReq, q.Type, fmt.Sprintf("%s/%s", q.Action, q.Super.Action), q.User, q.RemoteAddr)

	if s.readonly && (q.Super.Action == `grant` || q.Super.Action == `revoke`) {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto dispatch
	}

	switch q.Grant.Category {
	case `global`:
		switch q.Super.Action {
		case `grant`:
			fallthrough
		case `revoke`:
			s.right_globalsystem_modify(q)
		default:
			s.right_globalsystem_read(q)
		}
	case `system`:
		switch q.Super.Action {
		case `grant`:
			fallthrough
		case `revoke`:
			s.right_globalsystem_modify(q)
		default:
			s.right_globalsystem_read(q)
		}
	case `limited`:
		switch q.Super.Action {
		case `grant`:
			fallthrough
		case `revoke`:
			s.right_limited_modify(q)
		default:
			s.right_limited_read(q)
		}
	}
	return

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_globalsystem_modify(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `right`, Super: &msg.Supervisor{Action: q.Super.Action}}
	userUUID, ok := s.id_user_rev.get(q.User)
	if !ok {
		userUUID = `00000000-0000-0000-0000-000000000000`
	}

	//XXX msg.Grant.RecipientType == user

	var (
		res  sql.Result
		err  error
		data []string
	)

	switch q.Super.Action {
	case `grant`:
		q.Grant.Id = uuid.NewV4().String()
		res, err = s.stmt_GrantSysGlUser.Exec(
			q.Grant.Id,
			q.Grant.RecipientId,
			q.Grant.PermissionId,
			q.Grant.Category,
			userUUID,
		)
	case `revoke`:
		// data = []string{userID, permissionID}
		if data = s.global_grants.get(q.Grant.Id); data == nil {
			result.NotFound(fmt.Errorf(`Supervisor: unknown`))
			goto dispatch
		}
		q.Grant.RecipientId = data[0]
		q.Grant.PermissionId = data[1]

		res, err = s.stmt_RevkSysGlUser.Exec(
			q.Grant.Id,
		)
	}
	if err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if result.RowCnt(res.RowsAffected()) {
		result.Grant = []proto.Grant{q.Grant}
		// keep lookup maps in sync
		switch q.Super.Action {
		case `grant`:
			s.global_permissions.grant(q.Grant.RecipientId,
				q.Grant.PermissionId, q.Grant.Id)
			s.global_grants.record(q.Grant.RecipientId,
				q.Grant.PermissionId, q.Grant.Id)
		case `revoke`:
			s.global_grants.discard(q.Grant.Id)
			s.global_permissions.revoke(q.Grant.RecipientId,
				q.Grant.PermissionId)
		}
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_globalsystem_read(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `right`, Super: &msg.Supervisor{Action: q.Super.Action}}

	var (
		grantId string
		err     error
	)

	switch q.Super.Action {
	case `search`:
		if err = s.stmt_SrchGlSysGrant.QueryRow(
			q.Grant.PermissionId,
			q.Grant.Category,
			q.Grant.RecipientId,
		).Scan(grantId); err == sql.ErrNoRows {
			result.NotFound(err)
			goto dispatch
		} else if err != nil {
			result.ServerError(err)
			goto dispatch
		}
		result.Grant = []proto.Grant{proto.Grant{
			Id:            grantId,
			PermissionId:  q.Grant.PermissionId,
			Category:      q.Grant.Category,
			RecipientId:   q.Grant.RecipientId,
			RecipientType: q.Grant.RecipientType,
		}}
	default:
		result.ServerError(nil)
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) right_limited_modify(q *msg.Request) {
}

func (s *supervisor) right_limited_read(q *msg.Request) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
