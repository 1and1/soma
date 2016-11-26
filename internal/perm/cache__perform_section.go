/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/1and1/soma/internal/msg"

// These are the per-Section methods used in Cache.Perform

func (c *Cache) performRepository(q *msg.Request) {
	switch q.Action {
	case `create`:
		c.performRepositoryCreate(q)
	case `destroy`:
		c.performRepositoryDestroy(q)
	}
}

func (c *Cache) performBucket(q *msg.Request) {
	switch q.Action {
	case `create`:
		c.performBucketCreate(q)
	case `destroy`:
		c.performBucketDestroy(q)
	}
}

func (c *Cache) performGroup(q *msg.Request) {
	switch q.Action {
	case `create`:
		c.performGroupCreate(q)
	case `destroy`:
		c.performGroupDestroy(q)
	}
}

func (c *Cache) performCluster(q *msg.Request) {
	switch q.Action {
	case `create`:
		c.performClusterCreate(q)
	case `destroy`:
		c.performClusterDestroy(q)
	}
}

func (c *Cache) performNode(q *msg.Request) {
	switch q.Action {
	case `assign`:
		c.performNodeAssign(q)
	case `unassign`:
		c.performNodeUnassign(q)
	}
}

func (c *Cache) performUser(q *msg.Request) {
	switch q.Action {
	case `add`:
		c.performUserAdd(q)
	case `remove`:
		c.performUserRemove(q)
	}
}

func (c *Cache) performTeam(q *msg.Request) {
	switch q.Action {
	case `add`:
		c.performTeamAdd(q)
	case `remove`:
		c.performTeamRemove(q)
	}
}

func (c *Cache) performRight(q *msg.Request) {
	switch q.Action {
	case `grant`:
		c.performRightGrant(q)
	case `revoke`:
		c.performRightRevoke(q)
	}
}

func (c *Cache) performPermission(q *msg.Request) {
	switch q.Action {
	case `remove`:
		c.performPermissionRemove(q)
	case `map`:
		c.performPermissionMap(q)
	case `unmap`:
		c.performPermissionUnmap(q)
	}
}

func (c *Cache) performCategory(q *msg.Request) {
	switch q.Action {
	case `remove`:
		c.performCategoryRemove(q)
	}
}

func (c *Cache) performSection(q *msg.Request) {
	switch q.Action {
	case `add`:
		c.performSectionAdd(q)
	case `remove`:
		c.performSectionRemove(q)
	}
}

func (c *Cache) performAction(q *msg.Request) {
	switch q.Action {
	case `add`:
		c.performActionAdd(q)
	case `remove`:
		c.performActionRemove(q)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
