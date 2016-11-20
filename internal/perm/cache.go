/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// Package perm implements the permission cache module for the
// SOMA supervisor. It tracks which actions are mapped to permissions
// and which permissions have been granted.
//
// It can be queried whether a given user is authorized to perform
// an action.
package perm // import "github.com/1and1/soma/internal/perm"

import "sync"

// Cache is a permission cache for the SOMA supervisor
type Cache struct {
	// the entire cache has one global mutex, since many actions
	// requires updates to multiple data structures. Locking each
	// of them individually could therefor lead to deadlocks.
	// A global lock is more robust than a lock order scheme, which
	// could still be adopted later as a performance improvement.
	lock sync.RWMutex

	// general ID<>name lookup maps
	section *sectionLookup
	action  *actionLookup
	user    *userLookup
	team    *teamLookup

	// keeps track which actions are mapped to which permissions
	pmap *permissionMapping

	// keeps track of permission grants
	grantGlobal     *unscopedGrantMap
	grantRepository *scopedGrantMap
	grantTeam       *scopedGrantMap
	grantMonitoring *scopedGrantMap
}

// New returns a new permission cache
func New() *Cache {
	c := Cache{}
	c.lock = sync.RWMutex{}
	c.section = newSectionLookup()
	c.action = newActionLookup()
	c.user = newUserLookup()
	c.team = newTeamLookup()
	c.pmap = newPermissionMapping()
	c.grantGlobal = newUnscopedGrantMap()
	c.grantRepository = newScopedGrantMap()
	c.grantTeam = newScopedGrantMap()
	c.grantMonitoring = newScopedGrantMap()
	return &c
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
