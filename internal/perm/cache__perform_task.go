/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/1and1/soma/internal/msg"

// These are the methods used when an action can have variations

func (c *Cache) performPermissionMapAction(q *msg.Request) {
	for _, a := range *q.Permission.Actions {
		c.pmap.mapAction(
			a.SectionId,
			a.Id,
			q.Permission.Id,
		)
	}
}

func (c *Cache) performPermissionMapSection(q *msg.Request) {
	for _, s := range *q.Permission.Sections {
		c.pmap.mapSection(
			s.Id,
			q.Permission.Id,
		)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
