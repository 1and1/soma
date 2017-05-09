/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

// Start launches all application handlers
func (s *Soma) Start() {
	// grimReaper and supervisor must run first
	// TODO: start grimReaper completely
	// TODO: start supervisor completely

	s.handlerMap.Add(`mode_r`, newModeRead(s.conf.QueueLen))
	s.handlerMap.Add(`node_r`, newNodeRead(s.conf.QueueLen))
	s.handlerMap.Add(`oncall_r`, newOncallRead(s.conf.QueueLen))
	s.handlerMap.Add(`predicate_r`, newPredicateRead(s.conf.QueueLen))
	s.handlerMap.Add(`provider_r`, newProviderRead(s.conf.QueueLen))
	s.handlerMap.Add(`server_r`, newServerRead(s.conf.QueueLen))
	s.handlerMap.Add(`status_r`, newStatusRead(s.conf.QueueLen))
	s.handlerMap.Add(`team_r`, newTeamRead(s.conf.QueueLen))
	s.handlerMap.Add(`unit_r`, newUnitRead(s.conf.QueueLen))
	s.handlerMap.Add(`user_r`, newUserRead(s.conf.QueueLen))
	s.handlerMap.Add(`validity_r`, newValidityRead(s.conf.QueueLen))
	s.handlerMap.Add(`view_r`, newViewRead(s.conf.QueueLen))

	if !s.conf.ReadOnly {
		if !s.conf.Observer {
			s.handlerMap.Add(`mode_w`, newModeWrite(s.conf.QueueLen))
			s.handlerMap.Add(`node_w`, newNodeWrite(s.conf.QueueLen))
			s.handlerMap.Add(`oncall_w`, newOncallWrite(s.conf.QueueLen))
			s.handlerMap.Add(`predicate_w`, newPredicateWrite(s.conf.QueueLen))
			s.handlerMap.Add(`provider_w`, newProviderWrite(s.conf.QueueLen))
			s.handlerMap.Add(`server_w`, newServerWrite(s.conf.QueueLen))
			s.handlerMap.Add(`status_w`, newStatusWrite(s.conf.QueueLen))
			s.handlerMap.Add(`team_w`, newTeamWrite(s.conf.QueueLen, s))
			s.handlerMap.Add(`unit_w`, newUnitWrite(s.conf.QueueLen))
			s.handlerMap.Add(`user_w`, newUserWrite(s.conf.QueueLen, s))
			s.handlerMap.Add(`validity_w`, newValidityWrite(s.conf.QueueLen))
			s.handlerMap.Add(`view_w`, newViewWrite(s.conf.QueueLen))
		}
	}

	// fully initialize the handlers and fire them up
	for handler := range s.handlerMap.Range() {
		switch handler {
		case `supervisor`, `grimReaper`:
			// already running
			continue
		}
		s.handlerMap.Register(
			handler,
			s.dbConnection,
			s.exportLogger(),
		)
		s.handlerMap.Run(handler)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
