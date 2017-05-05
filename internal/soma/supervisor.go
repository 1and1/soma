/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"

	"github.com/1and1/soma/internal/msg"
	"github.com/Sirupsen/logrus"
)

// XXX this is a stripped down bare-metal definition of the supervisor
// struct to make UserWrite compile until the Supervisor handler is properly
// migrated

// Supervisor is the handler that makes all authentication and authorization
// decisions
type Supervisor struct {
	Input    chan msg.Request
	Update   chan msg.Request
	Shutdown chan bool
	conn     *sql.DB
	appLog   *logrus.Logger
	reqLog   *logrus.Logger
	errLog   *logrus.Logger
}

// register initializes resources provided by the Soma app
func (s *Supervisor) register(c *sql.DB, l ...*logrus.Logger) {
	s.conn = c
	s.appLog = l[0]
	s.reqLog = l[1]
	s.errLog = l[2]
}

// run is the event loop for Supervisor
func (s *Supervisor) run() {
	// TODO required for Handler interface
}

// shutdown signals the handler to shut down
func (s *Supervisor) shutdownNow() {
	close(s.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
