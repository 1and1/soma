package somatree

import (
	"fmt"
	"log"
	"reflect"

	"github.com/satori/go.uuid"
)

type SomaTreeElemFault struct {
	Id     uuid.UUID
	Name   string
	Type   string
	State  string
	Parent SomaTreeFaultReceiver `json:"-"`
	Errors []error
}

//
// NEW
func NewFault() *SomaTreeElemFault {
	tef := new(SomaTreeElemFault)
	tef.Id = uuid.NewV4()
	tef.Type = "fault"
	tef.Name = "McFaulty"
	tef.Errors = make([]error, 0)
	tef.State = "floating"

	return tef
}

//
// Interface: SomaTreeBuilder
func (tef *SomaTreeElemFault) GetID() string {
	return tef.Id.String()
}

func (tef *SomaTreeElemFault) GetName() string {
	return tef.Name
}

func (tef *SomaTreeElemFault) GetType() string {
	return tef.Type
}

func (tef SomaTreeElemFault) CloneRepository() SomaTreeRepositoryAttacher {
	return &tef
}

//
// Interface: SomaTreeAttacher
func (tef *SomaTreeElemFault) Attach(a AttachRequest) {
	switch {
	case a.ParentType == "repository":
		tef.attachToRepository(a)
	}
}

func (tef *SomaTreeElemFault) ReAttach(a AttachRequest) {
	log.Fatal("Not implemented")
}

func (tef *SomaTreeElemFault) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeFaultReceiver:
		tef.setFaultParent(p.(SomaTreeFaultReceiver))
		tef.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemFault.setParent`)
	}
}

func (tef *SomaTreeElemFault) updateParentRecursive(p SomaTreeReceiver) {
	tef.setParent(p)
}

func (tef *SomaTreeElemFault) setFaultParent(p SomaTreeFaultReceiver) {
	tef.Parent = p
}

func (tef *SomaTreeElemFault) clearParent() {
	tef.Parent = nil
	tef.State = "floating"
}

func (tef *SomaTreeElemFault) Destroy() {
	if tef.Parent == nil {
		panic(`SomaTreeElemFault.Destroy called without Parent to unlink from`)
	}

	tef.Parent.Unlink(UnlinkRequest{
		ParentType: tef.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   tef.Parent.(SomaTreeBuilder).GetID(),
		ParentName: tef.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  tef.GetType(),
		ChildName:  tef.GetName(),
		ChildId:    tef.GetID(),
	},
	)
}

func (tef *SomaTreeElemFault) Detach() {
	tef.Destroy()
}

//
// Interface: SomaTreeRepositoryAttacher
func (tef *SomaTreeElemFault) attachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tef.Type,
		Fault:      tef,
	})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
