/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package perm

import "github.com/1and1/soma/lib/proto"

// sectionLookup is the cache data structure for permission sections,
// allowing lookup by ID or name
type sectionLookup struct {
	// see struct permissionMapping for detailed explanaition
	// of this counter
	compactionCounter int64
	byName            map[string]*proto.Section
	byID              map[string]*proto.Section
	byCategory        map[string][]*proto.Section
}

// newSectionLookup returns an initialized sectionLookup
func newSectionLookup() *sectionLookup {
	s := sectionLookup{}
	s.compactionCounter = 0
	s.byName = map[string]*proto.Section{}
	s.byID = map[string]*proto.Section{}
	s.byCategory = map[string][]*proto.Section{}
	return &s
}

// add inserts a section into the cache
func (m *sectionLookup) add(ID, name, category string) {
	s := &proto.Section{
		Id:       ID,
		Name:     name,
		Category: category,
	}
	m.byName[s.Name] = s
	m.byID[s.Id] = s

	if _, ok := m.byCategory[category]; !ok {
		m.byCategory[category] = []*proto.Section{}
	}
	m.byCategory[category] = append(m.byCategory[category], s)
}

// rmByName removes a section from the cache. The section is identified
// by its name.
func (m *sectionLookup) rmByName(name string) {
	if name == `` {
		return
	}

	s, ok := m.byName[name]
	if !ok {
		return
	}

	delete(m.byID, s.Id)
	delete(m.byName, s.Name)
	for i, p := range m.byCategory[s.Category] {
		if p.Id != s.Id {
			continue
		}
		m.byCategory[s.Category] = append(m.byCategory[s.Category][:i],
			m.byCategory[s.Category][i+1:]...)
		m.compactionCounter++
		break
	}
}

// rmByID removes a section from the cache. The section is identified
// by its name.
func (m *sectionLookup) rmByID(id string) {
	if id == `` {
		return
	}

	s, ok := m.byID[id]
	if !ok {
		return
	}

	delete(m.byID, s.Id)
	delete(m.byName, s.Name)
	for i, p := range m.byCategory[s.Category] {
		if p.Id != s.Id {
			continue
		}
		m.byCategory[s.Category] = append(m.byCategory[s.Category][:i],
			m.byCategory[s.Category][i+1:]...)
		m.compactionCounter++
		break
	}
}

// getByID returns a section from the cache. The section is identified
// by its ID. Returns nil if the section was not found.
func (m *sectionLookup) getByID(id string) *proto.Section {
	return m.byID[id]
}

// getByName returns a section from the cache. The section is
// identified by its name. Returns nil if the section was not found.
func (m *sectionLookup) getByName(name string) *proto.Section {
	return m.byName[name]
}

// getCategory returns all sections in that category
func (m *sectionLookup) getCategory(category string) []*proto.Section {
	return m.byCategory[category]
}

// compact copies all slices in the byCategory field into new slices
// to free the underlying array, then resets compactionCounter to
// zero
func (m *sectionLookup) compact() {
	for category, _ := range m.byCategory {
		nsl := make([]*proto.Section, len(m.byCategory[category]))
		copy(nsl, m.byCategory[category])
		m.byCategory[category] = nsl
	}
	m.compactionCounter = 0
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
