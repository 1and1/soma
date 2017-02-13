/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm // import "github.com/1and1/soma/internal/perm"

import (
	"strings"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
)

// isAuthorized implements Cache.IsAuthorized and checks if the
// request is authorized
func (c *Cache) isAuthorized(q *msg.Request) msg.Result {
	result := msg.FromRequest(q)
	// default action is to deny
	result.Super = &msg.Supervisor{
		Verdict:      401,
		VerdictAdmin: false,
	}
	var user *proto.User
	var subjType, category, permID, actionID, sectionID, objID string
	var action *proto.Action
	var sectionPermIDs, actionPermIDs, mergedPermIDs []string
	var any bool

	// determine type of the request subject
	switch {
	case strings.HasPrefix(q.User, `admin_`):
		subjType = `admin`
	case strings.HasPrefix(q.User, `tool_`):
		subjType = `tool`
	default:
		subjType = `user`
	}

	// set readlock on the cache
	c.lock.RLock()
	defer c.lock.RUnlock()

	// look up the user, also handles admin and tool accounts
	if user = c.user.getByName(q.User); user == nil {
		goto dispatch
	}

	// check if the subject has omnipotence
	if c.grantGlobal.assess(
		subjType,
		user.Id,
		`omnipotence`,
		`00000000-0000-0000-0000-000000000000`,
	) {
		result.Super.Verdict = 200
		result.Super.VerdictAdmin = true
		goto dispatch
	}

	// check if the user has the system permission
	category = c.section.getByName(q.Section).Category
	permID = c.pmap.getIDByName(`system`, category)
	if permID == `` {
		goto dispatch
	}
	if c.grantGlobal.assess(
		subjType,
		user.Id,
		`system`,
		permID,
	) {
		result.Super.Verdict = 200
		result.Super.VerdictAdmin = true
		goto dispatch
	}

	// check if the user has a specific grant for the action
	// lookup sectionID and actionID of the Request
	if action = c.action.getByName(
		q.Section,
		q.Action,
	); action == nil {
		goto dispatch
	}
	sectionID = action.SectionId
	actionID = action.Id
	// lookup all permissionIDs that map either section or action
	sectionPermIDs = c.pmap.getSectionPermissionID(sectionID)
	actionPermIDs = c.pmap.getActionPermissionID(sectionID, actionID)
	mergedPermIDs = append(sectionPermIDs, actionPermIDs...)
	// check if we care about the specific object
	switch q.Action {
	case `list`, `search`:
		any = true
	}
	// check if the user has one the permissions that map the
	// requested action
	for _, permID = range mergedPermIDs {
		switch q.Section {
		case `monitoringsystem`:
			switch {
			case any:
				// invalid uuid
				objID = `ffffffff-ffff-3fff-ffff-ffffffffffff`
			default:
				objID = q.Monitoring.Id
			}
			if c.grantMonitoring.assess(subjType, user.Id, category,
				objID, permID, any) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
		case `repository`:
			switch {
			case any:
				// invalid uuid
				objID = `ffffffff-ffff-3fff-ffff-ffffffffffff`
			default:
				objID = q.Repository.Id
			}
			if c.grantRepository.assess(subjType, user.Id, category,
				objID, permID, any) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
		case `bucket`:
			switch {
			case any:
				// invalid uuid
				objID = `ffffffff-ffff-3fff-ffff-ffffffffffff`
			default:
				objID = q.Bucket.Id
			}
			if c.grantRepository.assess(subjType, user.Id, category,
				objID, permID, any) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
			// permission could be on the repository
			objID = c.object.repoForBucket(q.Bucket.Id)
			if objID == `` {
				goto dispatch
			}
			if c.grantRepository.assess(subjType, user.Id, category,
				objID, permID, any) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
		default:
			if c.grantGlobal.assess(subjType, user.Id, category,
				permID) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
		}
	}

	// check if the user's team has a specific grant for the action
	switch subjType {
	case `admin`, `tool`:
		// admin and tool accounts do not inherit team rights
		goto dispatch
	}
	// check if the team has one the permissions that map the
	// requested action
	for _, permID = range mergedPermIDs {
		switch q.Section {
		case `monitoringsystem`:
			switch {
			case any:
				// invalid uuid
				objID = `ffffffff-ffff-3fff-ffff-ffffffffffff`
			default:
				objID = q.Monitoring.Id
			}
			if c.grantMonitoring.assess(`team`, user.TeamId, category,
				objID, permID, any) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
		case `repository`:
			switch {
			case any:
				// invalid uuid
				objID = `ffffffff-ffff-3fff-ffff-ffffffffffff`
			default:
				objID = q.Repository.Id
			}
			if c.grantRepository.assess(subjType, user.Id, category,
				objID, permID, any) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
		case `bucket`:
			switch {
			case any:
				// invalid uuid
				objID = `ffffffff-ffff-3fff-ffff-ffffffffffff`
			default:
				objID = q.Bucket.Id
			}
			if c.grantRepository.assess(subjType, user.Id, category,
				objID, permID, any) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
			// permission could be on the repository
			objID = c.object.repoForBucket(q.Bucket.Id)
			if objID == `` {
				goto dispatch
			}
			if c.grantRepository.assess(subjType, user.Id, category,
				objID, permID, any) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
		default:
			if c.grantGlobal.assess(`team`, user.TeamId, category,
				permID) {
				result.Super.Verdict = 200
				result.Super.VerdictAdmin = false
				goto dispatch
			}
		}
	}

dispatch:
	return result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
