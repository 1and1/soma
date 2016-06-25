package somatree

import (
	"fmt"
	"reflect"

	"github.com/satori/go.uuid"
)

type Fault struct {
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
func newFault() *Fault {
	tef := new(Fault)
	tef.Id = uuid.NewV4()
	tef.Type = "fault"
	tef.Name = "McFaulty"
	tef.Errors = make([]error, 0)
	tef.State = "floating"
	tef.Parent = nil

	return tef
}

func (tef *Fault) getErrors() []error {
	err := make([]error, len(tef.Errors))
	copy(err, tef.Errors)
	tef.Errors = make([]error, 0)
	return err
}

//
// Interface: SomaTreeBuilder
func (tef *Fault) GetID() string {
	return tef.Id.String()
}

func (tef *Fault) GetName() string {
	return tef.Name
}

func (tef *Fault) GetType() string {
	return tef.Type
}

func (tef Fault) CloneRepository() SomaTreeRepositoryAttacher {
	return &tef
}

func (tef Fault) CloneBucket() SomaTreeBucketAttacher {
	return &tef
}

func (tef *Fault) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeFaultReceiver:
		tef.setFaultParent(p.(SomaTreeFaultReceiver))
		tef.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`Fault.setParent`)
	}
}

func (tef *Fault) setAction(c chan *Action) {
	tef.Action = c

	tef.Action <- &Action{
		Action: "create",
		Type:   "fault",
		//Id:     tef.Id.String(),
		//Name:   tef.Name,
	}
}

func (tef *Fault) setActionDeep(c chan *Action) {
	tef.Action = c
}

func (tef *Fault) setError(c chan *Error) {
	tef.Error = c

	tef.Action <- &Action{
		Action: "attached",
		Type:   "errorchannel",
	}
}

func (tef *Fault) updateParentRecursive(p SomaTreeReceiver) {
	tef.setParent(p)
}

func (tef *Fault) setFaultParent(p SomaTreeFaultReceiver) {
	tef.Parent = p
}

func (tef *Fault) clearParent() {
	tef.Parent = nil
	tef.State = "floating"
}

// noop, but satisfy the interface
func (tef *Fault) setFault(f *Fault) {
}

// noop, but satisfy the interface
func (tef *Fault) updateFaultRecursive(f *Fault) {
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
func (tef *Fault) GetBucket() SomaTreeReceiver {
	panic(`Fault.GetBucket`)
	return tef
}

func (tef *Fault) GetEnvironment() string {
	panic(`Fault.GetEnvironment`)
	return "none"
}

func (tef *Fault) ComputeCheckInstances() {
}

func (tef *Fault) ClearLoadInfo() {
}

func (tef *Fault) LoadInstance(i CheckInstance) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
