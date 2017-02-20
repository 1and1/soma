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

const (
	// Rfc3339Milli is a format string for millisecond precision RFC3339
	Rfc3339Milli string = "2006-01-02T15:04:05.000Z07:00"
	// LogStrReq is a format string for logging requests (deprecated)
	LogStrReq = `Subsystem=%s, Request=%s, User=%s, Addr=%s`
	// LogStrSRq is a format string for logging requests
	LogStrSRq = `Section=%s, Action=%s, User=%s, Addr=%s`
	// LogStrArg is a format string for logging scoped requests
	LogStrArg = `Subsystem=%s, Request=%s, User=%s, Addr=%s, Arg=%s`
	// LogStrOK is a format string for logging OK results
	LogStrOK = `Section=%s, Action=%s, InternalCode=%d, ExternalCode=%d`
	// LogStrErr is a format string for logging ERROR results
	LogStrErr = `Section=%s, Action=%s, InternalCode=%d, Error=%s`
)

// Soma application struct
type Soma struct {
	handlerMap   *HandlerMap
	dbConnection *sql.DB
	conf         *Config
}

// New returns a new SOMA application
func New(
	appHandlerMap *HandlerMap,
	dbConnection *sql.DB,
	conf *Config,
) *Soma {
	s := Soma{}
	s.handlerMap = appHandlerMap
	s.dbConnection = dbConnection
	s.conf = conf
	return &s
}

// Start launches all application handlers
func (s *Soma) Start() {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
