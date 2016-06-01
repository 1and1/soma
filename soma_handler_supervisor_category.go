package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

)

func (s *supervisor) permission_category(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `category`}

	log.Printf(LogStrReq, q.Type, fmt.Sprintf("%s/%s", q.Action, q.Super.Action), q.User, q.RemoteAddr)

	if s.readonly && (q.Super.Action == `add` || q.Super.Action == `delete`) {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto dispatch
	}

	switch q.Super.Action {
	case `list`:
		fallthrough
	case `show`:
		s.permission_category_read(q)
		return
	case `add`:
		fallthrough
	case `delete`:
		s.permission_category_write(q)
		return
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) permission_category_read(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `category`}
	var (
		rows           *sql.Rows
		err            error
		category, user string
		ts             time.Time
	)

	switch q.Super.Action {
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
				result.Clear(q.Action)
				goto dispatch
			}
			result.Category = append(result.Category,
				proto.Category{Name: category})
		}
		if err = rows.Err(); err != nil {
			result.ServerError(err)
			result.Clear(q.Action)
		}
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
	}

dispatch:
	q.Reply <- result
}

func (s *supervisor) permission_category_write(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `category`}
	userUUID, ok := s.id_user_rev.get(q.User)
	if !ok {
		userUUID = `00000000-0000-0000-0000-000000000000`
	}

	var (
		res sql.Result
		err error
	)

	switch q.Super.Action {
	case `add`:
		res, err = s.stmt_AddCategory.Exec(
			q.Category.Name,
			userUUID,
		)
	case `delete`:
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
