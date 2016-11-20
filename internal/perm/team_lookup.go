/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

// teamLookup is the cache data structure for teams, allowing lookup
// by ID or name
type teamLookup struct {
	// teamName -> teamID
	byName map[string]string
	// teamID -> teamName
	byID map[string]string
	// teamID -> []userID
	members map[string][]string
}

// newTeamLookup returns an initialized teamLookup
func newTeamLookup() *teamLookup {
	t := teamLookup{}
	t.byName = map[string]string{}
	t.byID = map[string]string{}
	t.members = map[string][]string{}
	return &t
}

// add inserts a team into the cache
func (m *teamLookup) add(teamID, teamName string) {
	m.byName[teamName] = teamID
	m.byID[teamID] = teamName

	// do not discard members if a team is added twice
	if _, ok := m.members[teamID]; !ok {
		m.members[teamID] = []string{}
	}
}

// addMember adds a new member to a team
func (m *teamLookup) addMember(teamID, userID string) {
	if _, ok := m.members[teamID]; !ok {
		return
	}

	m.members[teamID] = append(m.members[teamID], userID)
}

// getName returns the teamName for a teamID
func (m *teamLookup) getName(teamID string) string {
	return m.byID[teamID]
}

// getID returns the teamID for a teamName
func (m *teamLookup) getID(teamName string) string {
	return m.byName[teamName]
}

// getMemberByName returns the userIDs of the members of a team.
// If the team is not found, the slice is empty.
func (m *teamLookup) getMemberByName(teamName string) []string {
	return m.getMemberByID(m.byName[teamName])
}

// getMemberByName returns the userIDs of the members of a team.
// If the team is not found, the slice is empty.
func (m *teamLookup) getMemberByID(teamID string) []string {
	if _, ok := m.members[teamID]; !ok {
		return []string{}
	}

	// return a copy
	r := make([]string, len(m.members[teamID]))
	copy(r, m.members[teamID])
	return r
}

// rmByID removes a team from the cache. The team is identified by
// its ID.
func (m *teamLookup) rmByID(teamID string) {
	n, ok := m.byID[teamID]
	if !ok {
		return
	}

	delete(m.byName, n)
	delete(m.members, teamID)
	delete(m.byID, teamID)
}

// rmByName removes a team from the cache. The team is identified by
// its name.
func (m *teamLookup) rmByName(teamName string) {
	id, ok := m.byName[teamName]
	if !ok {
		return
	}

	delete(m.byID, id)
	delete(m.members, id)
	delete(m.byName, teamName)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
