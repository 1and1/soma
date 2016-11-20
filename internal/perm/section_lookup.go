package perm

import "github.com/1and1/soma/lib/proto"

// sectionLookup is the cache data structure for permission sections,
// allowing lookup by ID or name
type sectionLookup struct {
	byName map[string]*proto.Section
	byID   map[string]*proto.Section
}

// newSectionLookup returns an initialized sectionLookup
func newSectionLookup() *sectionLookup {
	s := sectionLookup{}
	s.byName = map[string]*proto.Section{}
	s.byID = map[string]*proto.Section{}
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
