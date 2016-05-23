package main


func (s *supervisor) update_map(q *msg.Request) {

	switch q.Super.Object {
	case `team`:
		switch q.Super.Action {
		case `add`:
		case `delete`:
		}
	}

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
