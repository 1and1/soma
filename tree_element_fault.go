package somatree

import (
	"fmt"
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
func newFault() *SomaTreeElemFault {
	tef := new(SomaTreeElemFault)
	tef.Id = uuid.NewV4()
	tef.Type = "fault"
	tef.Name = "McFaulty"
	tef.Errors = make([]error, 0)
	tef.State = "floating"
	tef.Parent = nil

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
	if tef.Parent != nil {
		panic(`SomaTreeElemFault.Attach: already attached`)
	}
	switch {
	case a.ParentType == "repository":
		tef.attachToRepository(a)
	}
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

// noop, but satisfy the interface
func (tef *SomaTreeElemFault) setFault(f *SomaTreeElemFault) {
}

// noop, but satisfy the interface
func (tef *SomaTreeElemFault) updateFaultRecursive(f *SomaTreeElemFault) {
}

// noop, but satisfy the interface
func (tef *SomaTreeElemFault) SetProperty(p SomaTreeProperty) {
}
func (tef *SomaTreeElemFault) inheritProperty(p SomaTreeProperty) {
}
func (tef *SomaTreeElemFault) inheritPropertyDeep(p SomaTreeProperty) {
}
func (tef *SomaTreeElemFault) setCustomProperty(p SomaTreeProperty) {
}

func (tef *SomaTreeElemFault) Destroy() {
	if tef.Parent == nil {
		panic(`SomaTreeElemFault.Destroy called without Parent to unlink from`)
	}

	tef.Parent.(SomaTreeAttacher).updateFaultRecursive(nil)

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

/*
 * Fault Handler Special Sauce
 *
 * Elemnts return pointers to the Fault Handler instead of nil pointers
 * when asked for something they do not have.
 *
 * This makes these chains safe:
 *		<foo>.Parent.(SomaTreeReceiver).GetBucket().Unlink()
 *
 * Instead of nil, the parent returns the Fault handler which implements
 * SomaTreeReceiver and SomaTreeUnlinker. Due to the information in the
 * Receive-/UnlinkRequest, it can log what went wrong.
 *
 */

//
// Interface: SomaTreeReceiver
func (tef *SomaTreeElemFault) Receive(r ReceiveRequest) {
	panic(`SomaTreeElemFault.Receive`)
}

//
// Interface: SomaTreeBucketeer
func (tef *SomaTreeElemFault) GetBucket() SomaTreeReceiver {
	panic(`SomaTreeElemFault.GetBucket`)
	return tef
}

//
// Interface: SomaTreeUnlinker
func (tef *SomaTreeElemFault) Unlink(u UnlinkRequest) {
	panic(`SomaTreeElemFault.Unlink`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
