/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "fmt"

// unscopedGrantMap is the cache data structure for global permission
// grants. It covers the categories 'omnipotence', 'system', 'global',
// 'permission' and 'operations'.
type unscopedGrantMap struct {
	// subjectId -> category -> permissionID -> grantID
	// The subjectID is built as follows:
	//	user:${userUUID}
	//	admin:${adminUUID}
	//	tool:${toolUUID}
	//	team:${teamUUID}
	grants map[string]map[string]map[string]string
	// grantID -> subject|category|permissionID
	byGrant map[string]map[string]string
}

// newUnscopedGrantMap returns an initialized unscopedGrantMap
func newUnscopedGrantMap() *unscopedGrantMap {
	u := unscopedGrantMap{}
	u.grants = map[string]map[string]map[string]string{}
	u.byGrant = map[string]map[string]string{}
	return &u
}

// grant records a grant of a permission to a subject into the cache
func (m *unscopedGrantMap) grant(subjType, subjID, category,
	permissionID, grantID string) {
	// only accept these four types
	switch subjType {
	case `user`, `admin`, `tool`, `team`:
	default:
		return
	}
	subject := fmt.Sprintf("%s:%s", subjType, subjID)

	// ensure the validity of the deep chain
	if _, ok := m.grants[subject]; !ok {
		m.grants[subject] = map[string]map[string]string{}
	}
	if _, ok := m.grants[subject][category]; !ok {
		m.grants[subject][category] = map[string]string{}
	}
	m.grants[subject][category][permissionID] = grantID
	m.byGrant[grantID] = map[string]string{
		`subjType`:     subjType,
		`subjID`:       subjID,
		`category`:     category,
		`permissionID`: permissionID,
	}
}

// revoke removes a grant of a permission from the cache
func (m *unscopedGrantMap) revoke(grantID string) {
	g, ok := m.byGrant[grantID]
	if !ok {
		return
	}
	subject := fmt.Sprintf("%s:%s", g[`subjType`], g[`subjID`])
	delete(m.grants[subject][g[`category`]], g[`permissionID`])
	delete(m.byGrant, grantID)
}

// getPermissionGrantID returns all grantIDs for a permissionID
func (m *unscopedGrantMap) getPermissionGrantID(
	permissionID string) []string {
	res := []string{}
	for grantID, grant := range m.byGrant {
		if grant[`permissionID`] == permissionID {
			res = append(res, grantID)
		}
	}
	return res
}

// getSubjectGrantID returns all grantIDs for a subjectID
func (m *unscopedGrantMap) getSubjectGrantID(subjType,
	subjID string) []string {
	res := []string{}
	for grantID, grant := range m.byGrant {
		if grant[`subjType`] != subjType {
			continue
		}
		if grant[`subjID`] == subjID {
			res = append(res, grantID)
		}
	}
	return res
}

// assess evaluates whether a subject has been granted a
// specific permission
func (m *unscopedGrantMap) assess(subjType, subjID, category,
	permissionID string) bool {
	// only accept these four types
	switch subjType {
	case `user`, `admin`, `tool`, `team`:
	default:
		return false
	}
	subject := fmt.Sprintf("%s:%s", subjType, subjID)

	if _, ok := m.grants[subject]; !ok {
		// subject has no grants
		return false
	}

	if _, ok := m.grants[subject][category]; !ok {
		// subject has no grants in category
		return false
	}

	if grantID, ok := m.grants[subject][category][permissionID]; ok {
		if grantID != `` {
			// subject has been granted the requested permission
			return true
		}
	}
	return false
}

// scopedGrantMap is the cache data structure for permission grants
// on an object.
type scopedGrantMap struct {
	// subjectID -> category -> permissionID -> objectID -> grantID
	// The subjectID is built as follows:
	//	user:${userUUID}
	//	admin:${adminUUID}
	//	tool:${toolUUID}
	//	team:${teamUUID}
	grants map[string]map[string]map[string]map[string]string
	// grantID -> subject|category|permissionID|objectID
	byGrant map[string]map[string]string
}

// newScopedGrantMap return ans initialized scopedGrantMap
func newScopedGrantMap() *scopedGrantMap {
	s := scopedGrantMap{}
	s.grants = map[string]map[string]map[string]map[string]string{}
	s.byGrant = map[string]map[string]string{}
	return &s
}

// grant records a grant of a permission on an object to a subject
// into the cache
func (m *scopedGrantMap) grant(subjType, subjID, category, objID,
	permissionID, grantID string) {
	// only accept these four types
	switch subjType {
	case `user`, `admin`, `tool`, `team`:
	default:
		return
	}
	subject := fmt.Sprintf("%s:%s", subjType, subjID)

	// ensure the validity of the deep chain
	if _, ok := m.grants[subject]; !ok {
		m.grants[subject] = map[string]map[string]map[string]string{}
	}
	if _, ok := m.grants[subject][category]; !ok {
		m.grants[subject][category] = map[string]map[string]string{}
	}
	if _, ok := m.grants[subject][category][permissionID]; !ok {
		m.grants[subject][category][permissionID] = map[string]string{}
	}
	m.grants[subject][category][permissionID][objID] = grantID
	m.byGrant[grantID] = map[string]string{
		`subjType`:     subjType,
		`subjID`:       subjID,
		`category`:     category,
		`objID`:        objID,
		`permissionID`: permissionID,
	}
}

// revoke removes a grant of a permission from the cache
func (m *scopedGrantMap) revoke(grantID string) {
	g, ok := m.byGrant[grantID]
	if !ok {
		return
	}
	subject := fmt.Sprintf("%s:%s", g[`subjType`], g[`subjID`])
	delete(m.grants[subject][g[`category`]][g[`permissionID`]],
		g[`objID`])
	delete(m.byGrant, grantID)
}

// getPermissionGrantID returns all grantIDs for a permissionID
func (m *scopedGrantMap) getPermissionGrantID(
	permissionID string) []string {
	res := []string{}
	for grantID, grant := range m.byGrant {
		if grant[`permissionID`] == permissionID {
			res = append(res, grantID)
		}
	}
	return res
}

// getObjectGrantID returns all grantIDs for an objectID
func (m *scopedGrantMap) getObjectGrantID(
	objectID string) []string {
	res := []string{}
	for grantID, grant := range m.byGrant {
		if grant[`objID`] == objectID {
			res = append(res, grantID)
		}
	}
	return res
}

// getSubjectGrantID returns all grantIDs for a subjectID
func (m *scopedGrantMap) getSubjectGrantID(subjType,
	subjID string) []string {
	res := []string{}
	for grantID, grant := range m.byGrant {
		if grant[`subjType`] != subjType {
			continue
		}
		if grant[`subjID`] == subjID {
			res = append(res, grantID)
		}
	}
	return res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
