package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) permission_category(q *msg.Request) {
	result := msg.FromRequest(q)

	s.reqLog.Printf(LogStrReq, q.Type, fmt.Sprintf("%s/%s", q.Section, q.Action), q.User, q.RemoteAddr)

	if s.readonly && (q.Action == `add` || q.Action == `remove`) {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto abort
	}

	switch q.Action {
	case `list`, `show`:
		s.permission_category_read(q)
	case `add`, `remove`:
		s.permission_category_write(q)
	default:
		result.NotImplemented(fmt.Errorf("Unknown requested action:"+
			" %s/%s/%s", q.Type, q.Section, q.Action))
		goto abort
	}
	return

abort:
	q.Reply <- result
}

func (s *supervisor) permission_category_read(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		rows           *sql.Rows
		err            error
		category, user string
		ts             time.Time
	)

	switch q.Action {
	case `list`:
		result.Category = []proto.Category{}
		if rows, err = s.stmt_ListCategory.Query(); err != nil {
			result.ServerError(err)
			goto dispatch
		}
		defer rows.Close()

		for rows.Next() {
			if err = rows.Scan(
				&category,
			); err != nil {
				result.ServerError(err)
				result.Clear(q.Section)
				goto dispatch
			}
			result.Category = append(result.Category,
				proto.Category{Name: category})
		}
		if err = rows.Err(); err != nil {
			result.ServerError(err)
			result.Clear(q.Section)
		}
		result.OK()
	case `show`:
		if err = s.stmt_ShowCategory.QueryRow(q.Category.Name).Scan(
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
		result.Category = []proto.Category{proto.Category{
			Name: category,
			Details: &proto.CategoryDetails{
				CreatedAt: ts.Format(rfc3339Milli),
				CreatedBy: user,
			},
		}}
		result.OK()
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) permission_category_write(q *msg.Request) {
	result := msg.FromRequest(q)
	userUUID, ok := s.id_user_rev.get(q.User)
	if !ok {
		userUUID = `00000000-0000-0000-0000-000000000000`
	}

	var (
		res    sql.Result
		err    error
		permId string
	)

	switch q.Action {
	case `add`:
		// create requested category
		if res, err = s.stmt_AddCategory.Exec(
			q.Category.Name,
			userUUID,
		); err != nil {
			result.ServerError(err)
			goto dispatch
		}
		// create grant category for requested category
		if res, err = s.stmt_AddCategory.Exec(
			fmt.Sprintf("%s:grant", q.Category.Name),
			userUUID,
		); err != nil {
			result.ServerError(err)
			goto dispatch
		}
		// create system permission for category
		permId = uuid.NewV4().String()
		if res, err = s.stmt_AddPermission.Exec(
			permId,
			q.Category.Name,
			`system`,
			userUUID,
		); err != nil {
			result.ServerError(err)
			goto dispatch
		}
	case `remove`:
		res, err = s.stmt_DelCategory.Exec(
			q.Category.Name,
		)
	}
	if err != nil {
		result.ServerError(err)
		goto dispatch
	}

	if result.RowCnt(res.RowsAffected()) {
		result.Category = []proto.Category{q.Category}
	}

dispatch:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
