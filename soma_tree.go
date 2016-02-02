package somatree

import "github.com/satori/go.uuid"

type SomaTree struct {
	Id    uuid.UUID
	Name  string
	Child *SomaTreeElemRepository
}

type SomaTreeAction struct {
	Action   string
	Type     string
	TypeId   string
	Id       string
	SourceId string
}

type SomaTreeCheck struct {
	Id uuid.UUID
}

func New(name string) *SomaTree {
	st := new(SomaTree)
	st.Id = uuid.NewV4()
	st.Name = name
	return st
}

func (st *SomaTree) GetID() string {
	return st.Id.String()
}

// Interface: SomaTreeReceiver
func (st *SomaTree) Receive(r ReceiveRequest) {
	switch {
	case r.ParentType == "root" &&
		r.ParentId == st.Id.String() &&
		r.ChildType == "repository":
		st.ReceiveRepository(r)
	default:
		if st.Child != nil {
			st.Child.Receive(r)
		} else {
			panic("not allowed")
		}
	}
}

// Interface: SomaTreeUnlinker
func (st *SomaTree) Unlink(u UnlinkRequest) {
	switch {
	case u.ParentType == "root" &&
		(u.ParentId == st.Id.String() ||
			u.ParentName == st.Name) &&
		u.ChildType == "repository" &&
		u.ChildName == st.Child.GetName():
		st.UnlinkRepository(u)
	default:
		if st.Child != nil {
			st.Child.Unlink(u)
		} else {
			panic("not allowed")
		}
	}
}

// Interface: SomaTreeRepositoryReceiver
func (st *SomaTree) ReceiveRepository(r ReceiveRequest) {
	switch {
	case r.ParentType == "root" &&
		r.ParentId == st.Id.String() &&
		r.ChildType == "repository":
		st.Child = r.Repository
		r.Repository.SetParent(st)
	default:
		panic("not allowed")
	}
}

// Interface: SomaTreeRepositoryUnlinker
func (st *SomaTree) UnlinkRepository(u UnlinkRequest) {
	switch {
	case u.ParentType == "root" &&
		u.ParentId == st.Id.String() &&
		u.ChildType == "repository" &&
		u.ChildName == st.Child.GetName():
		st.Child = nil
	default:
		panic("not allowed")
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
