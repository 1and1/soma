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

	if !s.conf.ReadOnly {
		if !s.conf.Observer {
			s.startNodeWrite()
		}
	}
}

// startNodeRead
func (s *Soma) startNodeRead() {
	nodeRead := NodeRead{}
	nodeRead.Input = make(chan msg.Request, 64)
	nodeRead.Shutdown = make(chan bool)
	nodeRead.conn = s.dbConnection
	nodeRead.appLog = s.appLog
	nodeRead.reqLog = s.reqLog
	nodeRead.errLog = s.errLog
	s.handlerMap.Add(`node_r`, &nodeRead)
	go nodeRead.run()
}

// startNodeWrite
func (s *Soma) startNodeWrite() {
	nodeWrite := NodeWrite{}
	nodeWrite.Input = make(chan msg.Request, 64)
	nodeWrite.Shutdown = make(chan bool)
	nodeWrite.conn = s.dbConnection
	nodeWrite.appLog = s.appLog
	nodeWrite.reqLog = s.reqLog
	nodeWrite.errLog = s.errLog
	s.handlerMap.Add(`node_w`, &nodeWrite)
	go nodeWrite.run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
