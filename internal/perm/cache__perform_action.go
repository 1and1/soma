/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/1and1/soma/internal/msg"

// These are the per-Action methods used in Cache.Perform

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

// performActionRemove removes an action from the cache after
// removing it from all permission maps
func (c *Cache) performActionRemove(q *msg.Request) {
	c.lock.Lock()
	c.performActionRemoveTask(
		q.ActionObj.SectionId,
		q.ActionObj.Id,
	)
	c.lock.Unlock()
}

func (c *Cache) performBucketCreate(q *msg.Request) {
}

func (c *Cache) performBucketDestroy(q *msg.Request) {
}

func (c *Cache) performCategoryRemove(q *msg.Request) {
}

func (c *Cache) performClusterCreate(q *msg.Request) {
}

func (c *Cache) performClusterDestroy(q *msg.Request) {
}

func (c *Cache) performGroupCreate(q *msg.Request) {
}

func (c *Cache) performGroupDestroy(q *msg.Request) {
}

func (c *Cache) performNodeAssign(q *msg.Request) {
}

func (c *Cache) performNodeUnassign(q *msg.Request) {
}

func (c *Cache) performPermissionMap(q *msg.Request) {
	c.lock.Lock()
	// map request can contain either actions or sections, not a mix
	if q.Permission.Actions != nil {
		c.performPermissionMapAction(q)
	}
	if q.Permission.Sections != nil {
		c.performPermissionMapSection(q)
	}
	c.lock.Unlock()
}

func (c *Cache) performPermissionRemove(q *msg.Request) {
}

func (c *Cache) performPermissionUnmap(q *msg.Request) {
}

func (c *Cache) performRepositoryCreate(q *msg.Request) {
}

func (c *Cache) performRepositoryDestroy(q *msg.Request) {
}

func (c *Cache) performRightGrant(q *msg.Request) {
	c.lock.Lock()
	switch q.Grant.Category {
	case `omnipotence`, `system`, `global`, `permission`, `operations`:
		c.performRightGrantUnscoped(q)
	case `repository`:
		c.performRightGrantScopeRepository(q)
	case `team`:
		c.performRightGrantScopeTeam(q)
	case `monitoring`:
		c.performRightGrantScopeMonitoring(q)
	}
	c.lock.Unlock()
}

func (c *Cache) performRightRevoke(q *msg.Request) {
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

// performSectionRemove removes a section after removing all its
// actions and permission mappings
func (c *Cache) performSectionRemove(q *msg.Request) {
	c.lock.Lock()
	c.performSectionRemoveTask(q.SectionObj.Id)
	c.lock.Unlock()
}

func (c *Cache) performTeamAdd(q *msg.Request) {
	c.lock.Lock()
	c.team.add(
		q.Team.Id,
		q.Team.Name,
	)
	c.lock.Unlock()
}

func (c *Cache) performTeamRemove(q *msg.Request) {
}

func (c *Cache) performUserAdd(q *msg.Request) {
	c.lock.Lock()
	c.user.add(
		q.UserObj.Id,
		q.UserObj.UserName,
		q.UserObj.TeamId,
	)
	c.team.addMember(
		q.UserObj.TeamId,
		q.UserObj.Id,
	)
	c.lock.Unlock()
}

func (c *Cache) performUserRemove(q *msg.Request) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
