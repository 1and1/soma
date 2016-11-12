package main

import "github.com/1and1/soma/internal/msg"

func (s *supervisor) update_map(q *msg.Request) {

	switch q.Super.Object {
	case `team`:
		switch q.Action {
		case `add`:
			s.id_team.insert(q.Super.Team.Id, q.Super.Team.Name)
		case `update`:
			s.id_team.insert(q.Super.Team.Id, q.Super.Team.Name)
		case `delete`:
			s.id_team.remove(q.Super.Team.Id)
		}
	case `user`:
		switch q.Action {
		case `add`:
			s.id_user.insert(q.Super.User.Id, q.Super.User.UserName)
			s.id_user_rev.insert(q.Super.User.UserName, q.Super.User.Id)
			s.id_userteam.insert(q.Super.User.Id, q.Super.User.TeamId)
		case `update`:
			oldname, _ := s.id_user.get(q.Super.User.Id)
			if oldname != q.Super.User.UserName {
				s.id_user_rev.remove(oldname)
			}
			s.id_user.insert(q.Super.User.Id, q.Super.User.UserName)
			s.id_user_rev.insert(q.Super.User.UserName, q.Super.User.Id)
			s.id_userteam.insert(q.Super.User.Id, q.Super.User.TeamId)
		case `delete`:
			if name, ok := s.id_user.get(q.Super.User.Id); ok {
				s.id_user_rev.remove(name)
			}
			s.id_user.remove(q.Super.User.Id)
			s.id_userteam.remove(q.Super.User.Id)
		}
	}

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
