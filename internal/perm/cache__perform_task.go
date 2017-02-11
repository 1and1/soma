/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/1and1/soma/internal/msg"

// These are the methods used when an action can have variations.
// Cache locking is performed by the action methods, tasks do not
// lock the cache!

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

func (c *Cache) performRightGrantUnscoped(q *msg.Request) {
	c.grantGlobal.grant(
		q.Grant.RecipientType,
		q.Grant.RecipientId,
		q.Grant.Category,
		q.Grant.PermissionId,
		q.Grant.Id,
	)
}

func (c *Cache) performRightGrantScopeRepository(q *msg.Request) {
	switch q.Grant.ObjectType {
	case `repository`, `bucket`:
		c.grantRepository.grant(
			q.Grant.RecipientType,
			q.Grant.RecipientId,
			q.Grant.Category,
			q.Grant.ObjectId,
			q.Grant.PermissionId,
			q.Grant.Id,
		)
	}
}

func (c *Cache) performRightGrantScopeTeam(q *msg.Request) {
	switch q.Grant.ObjectType {
	case `team`:
		c.grantTeam.grant(
			q.Grant.RecipientType,
			q.Grant.RecipientId,
			q.Grant.Category,
			q.Grant.ObjectId,
			q.Grant.PermissionId,
			q.Grant.Id,
		)
	}
}

func (c *Cache) performRightGrantScopeMonitoring(q *msg.Request) {
	switch q.Grant.ObjectType {
	case `monitoring`:
		c.grantMonitoring.grant(
			q.Grant.RecipientType,
			q.Grant.RecipientId,
			q.Grant.Category,
			q.Grant.ObjectId,
			q.Grant.PermissionId,
			q.Grant.Id,
		)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
