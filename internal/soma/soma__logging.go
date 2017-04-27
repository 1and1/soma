/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"github.com/1and1/soma/internal/msg"
	"github.com/Sirupsen/logrus"
)

func msgRequest(l *logrus.Logger, q *msg.Request) {
	l.Printf(LogStrSRq,
		q.Section,
		q.Action,
		q.User,
		q.RemoteAddr,
	)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
