/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/1and1/soma/lib/proto"

// userLookup is the cache data structure for users, allowing lookup
// by ID or name
type userLookup struct {
	// userName -> proto.User{Id:, TeamId:}
	byName map[string]*proto.User
	// userID -> proto.User{Id:, TeamId:}
	byID map[string]*proto.User
}

// newUserLookup returns and initialized userLookup
func newUserLookup() *userLookup {
	u := userLookup{}
	u.byName = map[string]*proto.User{}
	u.byID = map[string]*proto.User{}
	return &u
}

// add inserts a user into the cache
func (m *userLookup) add(userID, userName, teamID string) {
	u := &proto.User{
		Id:       userID,
		UserName: userName,
		TeamId:   teamID,
	}
	m.byName[userName] = u
	m.byID[userID] = u
}

// getByID returns a user from the cache. The user is identified by
// its ID. Returns nil if the user was not found.
func (m *userLookup) getByID(userID string) *proto.User {
	return m.byID[userID]
}

// getByName returns a user from the cache. The user is identified by
// its username. Returns nil if the user was not found.
func (m *userLookup) getByName(userName string) *proto.User {
	return m.byName[userName]
}

// rmByID removes a user from the cache. The user is identified by
// its ID.
func (m *userLookup) rmByID(userID string) {
	u, ok := m.byID[userID]
	if !ok {
		return
	}

	delete(m.byName, u.UserName)
	delete(m.byID, u.Id)
}

// rmByName removes a user from the cache. The user is identified by
// its username.
func (m *userLookup) rmByName(userName string) {
	u, ok := m.byName[userName]
	if !ok {
		return
	}

	delete(m.byID, u.Id)
	delete(m.byName, u.UserName)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
