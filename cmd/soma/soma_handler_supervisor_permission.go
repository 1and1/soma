package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) permission(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `permission`}

	s.reqLog.Printf(LogStrReq, q.Type, fmt.Sprintf("%s/%s", q.Action, q.Super.Action), q.User, q.RemoteAddr)

	if s.readonly && (q.Super.Action == `add` || q.Super.Action == `delete`) {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto dispatch
	}

	switch q.Super.Action {
	case `list`:
		fallthrough
	case `search/name`:
		fallthrough
	case `show`:
		s.permission_read(q)
		return
	case `add`:
		fallthrough
	case `delete`:
		s.permission_write(q)
		return
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) permission_read(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `permission`, Super: &msg.Supervisor{}}
	var (
		rows                     *sql.Rows
		err                      error
		id, name, category, user string
		ts                       time.Time
	)

	switch q.Super.Action {
	case `list`:
		fallthrough
	case `search/name`:
		result.Permission = []proto.Permission{}
		switch q.Super.Action {
		case `list`:
			result.Super.Action = `list`
			rows, err = s.stmt_ListPermission.Query()
		case `search/name`:
			result.Super.Action = `search/name`
			rows, err = s.stmt_SearchPerm.Query(q.Permission.Name)
		}
		if err != nil {
			result.ServerError(err)
			goto dispatch
		}
		defer rows.Close()

		for rows.Next() {
			if err = rows.Scan(
				&id,
				&name,
			); err != nil {
				result.ServerError(err)
				result.Clear(q.Action)
				goto dispatch
			}
			result.Permission = append(result.Permission,
				proto.Permission{Id: id, Name: name})
		}
		if err = rows.Err(); err != nil {
			result.ServerError(err)
			result.Clear(q.Action)
		}
		result.OK()
	case `show`:
		result.Super.Action = `show`
		if err = s.stmt_ShowPermission.QueryRow(q.Permission.Name).Scan(
			&id,
			&name,
			&category,
			&user,
			&ts,
		); err == sql.ErrNoRows {
			result.NotFound(err)
			goto dispatch
		} else if err != nil {
			result.ServerError(err)
			goto dispatch
		}
		result.Permission = []proto.Permission{proto.Permission{
			Id:       id,
			Name:     name,
			Category: category,
			Details: &proto.DetailsCreation{
				CreatedAt: ts.Format(rfc3339Milli),
				CreatedBy: user,
			},
		}}
		result.OK()
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) permission_write(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `permission`, Super: &msg.Supervisor{}}
	userUUID, ok := s.id_user_rev.get(q.User)
	if !ok {
		userUUID = `00000000-0000-0000-0000-000000000000`
	}

	var (
		res sql.Result
		err error
		id  string
	)

	switch q.Super.Action {
	case `add`:
		result.Super.Action = `add`
		q.Permission.Id = uuid.NewV4().String()
		res, err = s.stmt_AddPermission.Exec(
			q.Permission.Id,
			q.Permission.Name,
			q.Permission.Category,
			userUUID,
		)
	case `delete`:
		result.Super.Action = `delete`
		if id, ok = s.id_permission.get(q.Permission.Name); !ok {
			result.NotFound(fmt.Errorf(`Supervisor: unknown`))
			goto dispatch
		}
		res, err = s.stmt_DelPermission.Exec(
			id,
		)
	}
	if err != nil {
		result.ServerError(err)
		goto dispatch
	}

	if result.RowCnt(res.RowsAffected()) {
		result.Permission = []proto.Permission{q.Permission}
		// keep lookup maps in sync
		switch q.Super.Action {
		case `add`:
			s.id_permission.insert(q.Permission.Name, q.Permission.Id)
		case `delete`:
			s.id_permission.remove(q.Permission.Name)
		}
	}

dispatch:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
