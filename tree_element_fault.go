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
	Action chan *Action `json:"-"`
	Error  chan *Error  `json:"-"`
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

func (tef *SomaTreeElemFault) getErrors() []error {
	err := make([]error, len(tef.Errors))
	copy(err, tef.Errors)
	tef.Errors = make([]error, 0)
	return err
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

func (tef *SomaTreeElemFault) setAction(c chan *Action) {
	tef.Action = c

	tef.Action <- &Action{
		Action: "create",
		Type:   "fault",
		Id:     tef.Id.String(),
		Name:   tef.Name,
	}
}

func (tef *SomaTreeElemFault) setActionDeep(c chan *Action) {
	tef.Action = c
}

func (tef *SomaTreeElemFault) setError(c chan *Error) {
	tef.Error = c

	tef.Action <- &Action{
		Action: "attached",
		Type:   "errorchannel",
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
// Interface: SomaTreeBucketeer
func (tef *SomaTreeElemFault) GetBucket() SomaTreeReceiver {
	panic(`SomaTreeElemFault.GetBucket`)
	return tef
}

func (tef *SomaTreeElemFault) GetEnvironment() string {
	panic(`SomaTreeElemFault.GetEnvironment`)
	return "none"
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
