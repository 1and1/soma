/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
)

// Handler process a specific request type
type Handler interface {
	register(*sql.DB, ...*logrus.Logger)
	run()
	shutdownNow()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
