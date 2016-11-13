package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
)


func (s *supervisor) right_globalsystem_read(q *msg.Request) {
	result := msg.FromRequest(q)

	var (
		grantId string
		err     error
	)

	switch q.Action {
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

func (s *supervisor) right_limited_read(q *msg.Request) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
