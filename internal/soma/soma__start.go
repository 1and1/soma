/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import "github.com/1and1/soma/internal/msg"

// Start launches all application handlers
func (s *Soma) Start() {
	s.startNodeRead()
	s.startValidityRead()
	s.startViewRead()

	if !s.conf.ReadOnly {
		if !s.conf.Observer {
			s.startNodeWrite()
			s.startValidityWrite()
			s.startViewWrite()
		}
	}
}

// startNodeRead
func (s *Soma) startNodeRead() {
	nodeRead := NodeRead{}
	nodeRead.Input = make(chan msg.Request, 64)
	nodeRead.Shutdown = make(chan struct{})
	nodeRead.conn = s.dbConnection
	nodeRead.appLog = s.appLog
	nodeRead.reqLog = s.reqLog
	nodeRead.errLog = s.errLog
	s.handlerMap.Add(`node_r`, &nodeRead)
	go nodeRead.run()
}

// startValidityRead
func (s *Soma) startValidityRead() {
	validityRead := newValidityRead()
	validityRead.register(
		s.dbConnection,
		s.exportLogger()...,
	)
	s.handlerMap.Add(`validity_r`, &validityRead)
	go validityRead.run()
}

// startViewRead
func (s *Soma) startViewRead() {
	viewRead := ViewRead{}
	viewRead.Input = make(chan msg.Request, 64)
	viewRead.Shutdown = make(chan struct{})
	viewRead.conn = s.dbConnection
	viewRead.appLog = s.appLog
	viewRead.reqLog = s.reqLog
	viewRead.errLog = s.errLog
	s.handlerMap.Add(`view_r`, &viewRead)
	go viewRead.run()
}

// startNodeWrite
func (s *Soma) startNodeWrite() {
	nodeWrite := NodeWrite{}
	nodeWrite.Input = make(chan msg.Request, 64)
	nodeWrite.Shutdown = make(chan struct{})
	nodeWrite.conn = s.dbConnection
	nodeWrite.appLog = s.appLog
	nodeWrite.reqLog = s.reqLog
	nodeWrite.errLog = s.errLog
	s.handlerMap.Add(`node_w`, &nodeWrite)
	go nodeWrite.run()
}

// startValidityWrite
func (s *Soma) startValidityWrite() {
	validityWrite := newValidityWrite()
	validityWrite.register(
		s.dbConnection,
		s.exportLogger()...,
	)
	s.handlerMap.Add(`validity_w`, &validityWrite)
	go validityWrite.run()
}

// startViewWrite
func (s *Soma) startViewWrite() {
	viewWrite := ViewWrite{}
	viewWrite.Input = make(chan msg.Request, 64)
	viewWrite.Shutdown = make(chan struct{})
	viewWrite.conn = s.dbConnection
	viewWrite.appLog = s.appLog
	viewWrite.reqLog = s.reqLog
	viewWrite.errLog = s.errLog
	s.handlerMap.Add(`view_w`, &viewWrite)
	go viewWrite.run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
