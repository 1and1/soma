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
	var subjType, category, actionID, sectionID string
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
	if c.checkOmnipotence(subjType, user.Id) {
		result.Super.Verdict = 200
		result.Super.VerdictAdmin = true
		goto dispatch
	}

	// extract category
	category = c.section.getByName(q.Section).Category

	// lookup sectionID and actionID of the Request, abort for
	// unknown actions
	if action := c.action.getByName(
		q.Section,
		q.Action,
	); action == nil {
		goto dispatch
	} else {
		sectionID = action.SectionId
		actionID = action.Id
	}

	// check if the user has the correct system permission
	if ok, invalid := c.checkSystem(category, subjType,
		user.Id); invalid {
		goto dispatch
	} else if ok {
		result.Super.Verdict = 200
		result.Super.VerdictAdmin = true
		goto dispatch
	}

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
	if c.checkPermission(mergedPermIDs, any, q, subjType, user.Id,
		category) {
		result.Super.Verdict = 200
		result.Super.VerdictAdmin = false
		goto dispatch
	}

	// admin and tool accounts do not inherit team rights,
	// authorization check ends here
	switch subjType {
	case `admin`, `tool`:
		goto dispatch
	}

	// check if the user's team has a specific grant for the action
	if c.checkPermission(mergedPermIDs, any, q, `team`, user.TeamId,
		category) {
		result.Super.Verdict = 200
		result.Super.VerdictAdmin = false
	}

dispatch:
	return result
}

// checkOmnipotence returns true if the subject is omnipotent
func (c *Cache) checkOmnipotence(subjectType, subjectID string) bool {
	return c.grantGlobal.assess(
		subjectType,
		subjectID,
		`omnipotence`,
		`00000000-0000-0000-0000-000000000000`,
	)
}

// checkSystem returns true,false if the subject has the system
// permission for the category. If no system permission exists it
// returns false,true
func (c *Cache) checkSystem(category, subjectType,
	subjectID string) (bool, bool) {
	permID := c.pmap.getIDByName(`system`, category)
	if permID == `` {
		// there must be a system permission for every category,
		// refuse authorization since the permission cache is broken
		return false, true
	}
	return c.grantGlobal.assess(
		subjectType,
		subjectID,
		`system`,
		permID,
	), false
}

// checkPermission returns true if the subject has a grant for the
// requested action
func (c *Cache) checkPermission(permIDs []string, any bool,
	q *msg.Request, subjectType, subjectID, category string) bool {
	var objID string

permloop:
	for _, permID := range permIDs {
		// determine objID
		if any {
			// invalid uuid
			objID = `ffffffff-ffff-3fff-ffff-ffffffffffff`
		} else {
			switch q.Section {
			case `monitoringsystem`, `capability`:
				objID = q.Monitoring.Id
			case `repository`:
				objID = q.Repository.Id
			case `bucket`:
				objID = q.Bucket.Id
			case `node`, `property_service_team`:
				objID = user.TeamId
			}
		}

		// check authorization
		switch q.Section {
		case `monitoringsystem`, `capability`:
			// per-monitoring sections
			if c.grantMonitoring.assess(subjectType, subjectID,
				category, objID, permID, any) {
				return true
			}
		case `repository`, `bucket`:
			// per-repository sections
			if c.grantRepository.assess(subjectType, subjectID,
				category, objID, permID, any) {
				return true
			}
			if q.Section == `bucket` {
				// permission could be on the repository
				objID = c.object.repoForBucket(q.Bucket.Id)
				if objID == `` {
					continue permloop
				}
				if c.grantRepository.assess(subjectType, subjectID,
					category, objID, permID, any) {
					return true
				}
			}
		case `node`, `property_service_team`:
			// per-team sections
			if c.grantTeam.assess(subjectType, subjectID,
				category, objID, permID, any) {
				return true
			}
		default:
			// global sections
			if c.grantGlobal.assess(subjectType, subjectID, category,
				permID) {
				return true
			}
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
