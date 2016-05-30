/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

import (
)

type Supervisor struct {
	Verdict      uint16
	VerdictAdmin bool
	RemoteAddr   string
	// Fields for encrypted requests
	KexId string
	Data  []byte
	Kex   auth.Kex
	// Fields for basic authentication requests
	BasicAuthUser  string
	BasicAuthToken string
	Restricted     bool
	// Fields for permission authorization requests
	PermAction     string
	PermRepository string
	PermMonitoring string
	PermNode       string
	// Fields for map update notifications
	Action string
	Object string
	User   proto.User
	Team   proto.Team
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
