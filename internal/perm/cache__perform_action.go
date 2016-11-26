/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/1and1/soma/internal/msg"

// These are the per-Action methods used in Cache.Perform

func (c *Cache) performRepositoryCreate(q *msg.Request) {
}

func (c *Cache) performRepositoryDestroy(q *msg.Request) {
}

func (c *Cache) performBucketCreate(q *msg.Request) {
}

func (c *Cache) performBucketDestroy(q *msg.Request) {
}

func (c *Cache) performGroupCreate(q *msg.Request) {
}

func (c *Cache) performGroupDestroy(q *msg.Request) {
}

func (c *Cache) performClusterCreate(q *msg.Request) {
}

func (c *Cache) performClusterDestroy(q *msg.Request) {
}

func (c *Cache) performNodeAssign(q *msg.Request) {
}

func (c *Cache) performNodeUnassign(q *msg.Request) {
}

func (c *Cache) performUserAdd(q *msg.Request) {
}

func (c *Cache) performUserRemove(q *msg.Request) {
}

func (c *Cache) performTeamAdd(q *msg.Request) {
}

func (c *Cache) performTeamRemove(q *msg.Request) {
}

func (c *Cache) performRightGrant(q *msg.Request) {
}

func (c *Cache) performRightRevoke(q *msg.Request) {
}

func (c *Cache) performPermissionRemove(q *msg.Request) {
}

func (c *Cache) performPermissionMap(q *msg.Request) {
}

func (c *Cache) performPermissionUnmap(q *msg.Request) {
}

func (c *Cache) performCategoryRemove(q *msg.Request) {
}

func (c *Cache) performSectionAdd(q *msg.Request) {
	c.lock.Lock()
	c.section.add(
		q.SectionObj.Id,
		q.SectionObj.Name,
		q.SectionObj.Category,
	)
	c.lock.Unlock()
}

func (c *Cache) performSectionRemove(q *msg.Request) {
}

func (c *Cache) performActionAdd(q *msg.Request) {
	c.lock.Lock()
	c.action.add(
		q.ActionObj.SectionId,
		q.ActionObj.SectionName,
		q.ActionObj.Id,
		q.ActionObj.Name,
		q.ActionObj.Category,
	)
	c.lock.Unlock()
}

func (c *Cache) performActionRemove(q *msg.Request) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
