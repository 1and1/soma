package main

import (
	"database/sql"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
)

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
