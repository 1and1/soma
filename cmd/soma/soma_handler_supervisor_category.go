package main

import (
	"database/sql"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
)

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

	var (
		res sql.Result
		err error
	)

	switch q.Action {
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
