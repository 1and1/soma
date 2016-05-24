package main


func (s *supervisor) update_map(q *msg.Request) {

	switch q.Super.Object {
	case `team`:
		switch q.Super.Action {
		case `add`:
			s.id_team.insert(q.Super.Team.Id, q.Super.Team.Name)
		case `delete`:
			s.id_team.remove(q.Super.Team.Id)
		}
	case `user`:
		switch q.Super.Action {
		case `add`:
			s.id_user.insert(q.Super.User.Id, q.Super.User.UserName)
			s.id_userteam.insert(q.Super.User.Id, q.Super.User.TeamId)
		case `delete`:
			s.id_user.remove(q.Super.User.Id)
			s.id_userteam.remove(q.Super.User.Id)
		}
	}

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
