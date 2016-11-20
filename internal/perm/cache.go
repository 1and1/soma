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
	lock    sync.RWMutex
	section *sectionLookup
	action  *actionLookup
}

// New returns a new permission cache
func New() *Cache {
	c := Cache{}
	c.lock = sync.RWMutex{}
	c.section = newSectionLookup()
	c.action = newActionLookup()
	return &c
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
