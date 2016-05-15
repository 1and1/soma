/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

type Result struct {
	Type  string
	Error error
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
