package somatree

import (
	"log"

	"github.com/satori/go.uuid"
)

type SomaTree struct {
	Id     uuid.UUID
	Name   string
	Type   string
	Child  *SomaTreeElemRepository
	Snap   *SomaTreeElemRepository
	Action chan *Action `json:"-"`
}

type TreeSpec struct {
	Id     string
	Name   string
	Action chan *Action
}

func New(spec TreeSpec) *SomaTree {
	st := new(SomaTree)
	st.Id, _ = uuid.FromString(spec.Id)
	st.Name = spec.Name
	st.Action = spec.Action
	st.Type = "root"
	return st
}

func (st *SomaTree) Begin() {
	t := st.Child.Clone()
	st.Snap = &t
	st.Snap.updateParentRecursive(st)

	newFault().Attach(
		AttachRequest{
			Root:       st.Snap,
			ParentType: st.Snap.Type,
			ParentName: st.Snap.Name,
		},
	)
}

func (st *SomaTree) Rollback() {
	err := st.Child.Fault.Error
	ac := st.Child.Action

	st.Child = st.Snap
	st.Snap = nil
	st.Child.setActionDeep(ac)
	st.Child.setError(err)
}

func (st *SomaTree) Commit() {
	st.Snap = nil
}

//
// Interface: SomaTreeBuilder
func (st *SomaTree) GetID() string {
	return st.Id.String()
}

func (st *SomaTree) GetName() string {
	return st.Name
}

func (st *SomaTree) GetType() string {
	return st.Type
}

//
func (st *SomaTree) SetError(c chan *Error) {
	if st.Child != nil {
		st.Child.setError(c)
	}
}

func (st *SomaTree) GetErrors() []error {
	if st.Child != nil {
		return st.Child.getErrors()
	}
	return []error{}
}

// Interface: SomaTreeReceiver
func (st *SomaTree) Receive(r ReceiveRequest) {
	switch {
	case r.ParentType == "root" &&
		r.ParentId == st.Id.String() &&
		r.ChildType == "repository":
		st.receiveRepository(r)
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
		st.unlinkRepository(u)
	default:
		if st.Child != nil {
			st.Child.Unlink(u)
		} else {
			panic("not allowed")
		}
	}
}

// Interface: SomaTreeRepositoryReceiver
func (st *SomaTree) receiveRepository(r ReceiveRequest) {
	switch {
	case r.ParentType == "root" &&
		r.ParentId == st.Id.String() &&
		r.ChildType == "repository":
		st.Child = r.Repository
		r.Repository.setParent(st)
		r.Repository.setAction(st.Action)
	default:
		panic("not allowed")
	}
}

// Interface: SomaTreeRepositoryUnlinker
func (st *SomaTree) unlinkRepository(u UnlinkRequest) {
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

// Interface: SomaTreeFinder
func (st *SomaTree) Find(f FindRequest, b bool) SomaTreeAttacher {
	if !b {
		panic(`SomaTree.Find: root element can never inherit a Find request`)
	}

	res := st.Child.Find(f, false)
	if res != nil {
		return res
	}
	return st.Child.Fault
}

//
func (st *SomaTree) ComputeCheckInstances() {
	if st.Child == nil {
		panic(`SomaTree.ComputeCheckInstances: no repository registered`)
	}

	log.Printf("SomaTree[%s]: Action=%s, ObjectType=%s, ObjectId=%s",
		st.Name,
		`ComputeCheckInstances`,
		`tree`,
		st.Id.String(),
	)
	st.Child.ComputeCheckInstances()
	return
}

//
func (st *SomaTree) ClearLoadInfo() {
	if st.Child == nil {
		panic(`SomaTree.ClearLoadInfo: no repository registered`)
	}

	st.Child.ClearLoadInfo()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
