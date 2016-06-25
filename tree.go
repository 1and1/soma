package somatree

import (
	"log"

	"github.com/satori/go.uuid"
)

type Tree struct {
	Id     uuid.UUID
	Name   string
	Type   string
	Child  *Repository
	Snap   *Repository
	Action chan *Action `json:"-"`
}

type TreeSpec struct {
	Id     string
	Name   string
	Action chan *Action
}

func New(spec TreeSpec) *Tree {
	st := new(Tree)
	st.Id, _ = uuid.FromString(spec.Id)
	st.Name = spec.Name
	st.Action = spec.Action
	st.Type = "root"
	return st
}

func (st *Tree) Begin() {
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

func (st *Tree) Rollback() {
	err := st.Child.Fault.Error
	ac := st.Child.Action

	st.Child = st.Snap
	st.Snap = nil
	st.Child.setActionDeep(ac)
	st.Child.setError(err)
}

func (st *Tree) Commit() {
	st.Snap = nil
}

func (st *Tree) AttachError(err Error) {
	if st.Child != nil {
		st.Child.Fault.Error <- &err
	}
}

//
// Interface: Builder
func (st *Tree) GetID() string {
	return st.Id.String()
}

func (st *Tree) GetName() string {
	return st.Name
}

func (st *Tree) GetType() string {
	return st.Type
}

//
func (st *Tree) SetError(c chan *Error) {
	if st.Child != nil {
		st.Child.setError(c)
	}
}

func (st *Tree) GetErrors() []error {
	if st.Child != nil {
		return st.Child.getErrors()
	}
	return []error{}
}

// Interface: Receiver
func (st *Tree) Receive(r ReceiveRequest) {
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

// Interface: Unlinker
func (st *Tree) Unlink(u UnlinkRequest) {
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

// Interface: RepositoryReceiver
func (st *Tree) receiveRepository(r ReceiveRequest) {
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

// Interface: RepositoryUnlinker
func (st *Tree) unlinkRepository(u UnlinkRequest) {
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

// Interface: Finder
func (st *Tree) Find(f FindRequest, b bool) Attacher {
	if !b {
		panic(`Tree.Find: root element can never inherit a Find request`)
	}

	res := st.Child.Find(f, false)
	if res != nil {
		return res
	}
	return st.Child.Fault
}

//
func (st *Tree) ComputeCheckInstances() {
	if st.Child == nil {
		panic(`Tree.ComputeCheckInstances: no repository registered`)
	}

	log.Printf("Tree[%s]: Action=%s, ObjectType=%s, ObjectId=%s",
		st.Name,
		`ComputeCheckInstances`,
		`tree`,
		st.Id.String(),
	)
	st.Child.ComputeCheckInstances()
	return
}

//
func (st *Tree) ClearLoadInfo() {
	if st.Child == nil {
		panic(`Tree.ClearLoadInfo: no repository registered`)
	}

	st.Child.ClearLoadInfo()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
