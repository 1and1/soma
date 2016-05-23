package main


func (s *supervisor) authorize(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `verdict`}

	switch q.Super.PermAction {
	default:
		goto unauthorized
	}

unauthorized:
	result.Super = &msg.Supervisor{Verdict: 401}
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
