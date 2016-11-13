package main

import (
	"database/sql"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
)

func (s *supervisor) permission_read(q *msg.Request) {
	result := msg.FromRequest(q)
	var (
		rows                     *sql.Rows
		err                      error
		id, name, category, user string
		ts                       time.Time
	)

	switch q.Action {
	case `list`, `search/name`:
		result.Permission = []proto.Permission{}
		switch q.Action {
		case `list`:
			rows, err = s.stmt_ListPermission.Query()
		case `search/name`:
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
