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
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemFault.setParent`)
	}
}

func (tef *SomaTreeElemFault) setFaultParent(p SomaTreeFaultReceiver) {
	tef.Parent = p
}

func (tef *SomaTreeElemFault) Destroy() {
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

//
// Interface: SomaTreeReceiver
func (tef *SomaTreeElemFault) Receive(r ReceiveRequest) {
}

//
// Interface: SomaTreeBucketeer
func (tef *SomaTreeElemFault) GetBucket() SomaTreeReceiver {
	return tef
}

//
// Interface: SomaTreeUnlinker
func (tef *SomaTreeElemFault) Unlink(u UnlinkRequest) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
