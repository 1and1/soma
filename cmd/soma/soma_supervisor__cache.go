/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import "github.com/1and1/soma/internal/msg"

// cache handles all requests to update the supervisor
// permission cache
func (s *supervisor) cache(q *msg.Request) {
	s.permCache.Perform(q)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
