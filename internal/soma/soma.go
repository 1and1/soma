/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// Package soma implements the application handlers of the SOMA
// service.
package soma

import "database/sql"

// Soma application struct
type Soma struct {
	handlerMap   *HandlerMap
	dbConnection *sql.DB
}

// New returns a new SOMA application
func New(
	appHandlerMap *HandlerMap,
	dbConnection *sql.DB,
) *Soma {
	s := Soma{}
	s.handlerMap = appHandlerMap
	return &s
}

func (s *Soma) Start() {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
