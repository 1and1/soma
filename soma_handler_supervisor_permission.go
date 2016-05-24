package main

import (
	"fmt"
	"log"

)

func (s *supervisor) permission(q *msg.Request) {
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
}

func (s *supervisor) permission_write(q *msg.Request) {
}
