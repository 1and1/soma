package main

import (
	"fmt"
	"log"

)

func (s *supervisor) right(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `right`}

	log.Printf(LogStrReq, q.Type, fmt.Sprintf("%s/%s", q.Action, q.Super.Action), q.User, q.RemoteAddr)

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
			s.right_global_modify(q)
		default:
			s.right_global_read(q)
		}
	case `system`:
		switch q.Super.Action {
		case `grant`:
			fallthrough
		case `revoke`:
			s.right_system_modify(q)
		default:
			s.right_system_read(q)
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

func (s *supervisor) right_global_modify(q *msg.Request) {
}

func (s *supervisor) right_global_read(q *msg.Request) {
}

func (s *supervisor) right_system_modify(q *msg.Request) {
}

func (s *supervisor) right_system_read(q *msg.Request) {
}

func (s *supervisor) right_limited_modify(q *msg.Request) {
}

func (s *supervisor) right_limited_read(q *msg.Request) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
