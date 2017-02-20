/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"log"

	"github.com/1and1/soma/internal/msg"
)

// NodeWrite handles write requests for nodes
type NodeWrite struct {
	Input    chan msg.Request
	Shutdown chan bool
	conn     *sql.DB
	addStmt  *sql.Stmt
	delStmt  *sql.Stmt
	prgStmt  *sql.Stmt
	updStmt  *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
