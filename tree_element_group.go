package somatree

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/satori/go.uuid"
)

type SomaTreeElemGroup struct {
	Id              uuid.UUID
	Name            string
	State           string
	Team            uuid.UUID
	Type            string
	Parent          SomaTreeGroupReceiver `json:"-"`
	Fault           *SomaTreeElemFault    `json:"-"`
	Action          chan *Action          `json:"-"`
	PropertyOncall  map[string]SomaTreeProperty
	PropertyService map[string]SomaTreeProperty
	PropertySystem  map[string]SomaTreeProperty
	PropertyCustom  map[string]SomaTreeProperty
	Checks          map[string]SomaTreeCheck
	CheckInstances  map[string][]string
	Instances       map[string]SomaTreeCheckInstance
	Children        map[string]SomaTreeGroupAttacher //`json:"-"`
}

type GroupSepc struct {
	Id   uuid.UUID
	Name string
	Team uuid.UUID
}

//
// NEW
func NewGroup(name string, id string) *SomaTreeElemGroup {
	teg := new(SomaTreeElemGroup)
	if id != "" {
		teg.Id, _ = uuid.FromString(id)
	} else {
		teg.Id = uuid.NewV4()
	}
	teg.Name = name
	teg.Type = "group"
	teg.State = "floating"
	teg.Parent = nil
	teg.Children = make(map[string]SomaTreeGroupAttacher)
	teg.PropertyOncall = make(map[string]SomaTreeProperty)
	teg.PropertyService = make(map[string]SomaTreeProperty)
	teg.PropertySystem = make(map[string]SomaTreeProperty)
	teg.PropertyCustom = make(map[string]SomaTreeProperty)
	teg.Checks = make(map[string]SomaTreeCheck)

	return teg
}

func (teg SomaTreeElemGroup) CloneBucket() SomaTreeBucketAttacher {
	for k, child := range teg.Children {
		teg.Children[k] = child.CloneGroup()
	}
	return &teg
}

func (teg SomaTreeElemGroup) CloneGroup() SomaTreeGroupAttacher {
	f := make(map[string]SomaTreeGroupAttacher)
	for k, child := range teg.Children {
		f[k] = child.CloneGroup()
	}
	teg.Children = f
	return &teg
}

//
// Interface: SomaTreeBuilder
func (teg *SomaTreeElemGroup) GetID() string {
	return teg.Id.String()
}

func (teg *SomaTreeElemGroup) GetName() string {
	return teg.Name
}

func (teg *SomaTreeElemGroup) GetType() string {
	return teg.Type
}

func (teg *SomaTreeElemGroup) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case *SomaTreeElemBucket:
		teg.setGroupParent(p.(SomaTreeGroupReceiver))
		teg.State = "standalone"
	case *SomaTreeElemGroup:
		teg.setGroupParent(p.(SomaTreeGroupReceiver))
		teg.State = "grouped"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemGroup.setParent`)
	}
}

func (teg *SomaTreeElemGroup) setAction(c chan *Action) {
	teg.Action = c
}

// SomaTreeGroupReceiver == can receive Groups as children
func (teg *SomaTreeElemGroup) setGroupParent(p SomaTreeGroupReceiver) {
	teg.Parent = p
}

func (teg *SomaTreeElemGroup) updateParentRecursive(p SomaTreeReceiver) {
	teg.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func(str SomaTreeReceiver) {
			defer wg.Done()
			teg.Children[c].updateParentRecursive(str)
		}(teg)
	}
	wg.Wait()
}

func (teg *SomaTreeElemGroup) clearParent() {
	teg.Parent = nil
	teg.State = "floating"
}

func (teg *SomaTreeElemGroup) setFault(f *SomaTreeElemFault) {
	teg.Fault = f
}

func (teg *SomaTreeElemGroup) updateFaultRecursive(f *SomaTreeElemFault) {
	teg.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func(ptr *SomaTreeElemFault) {
			defer wg.Done()
			teg.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

//
// Interface: SomaTreeBucketeer
func (teg *SomaTreeElemGroup) GetBucket() SomaTreeReceiver {
	if teg.Parent == nil {
		if teg.Fault == nil {
			panic(`SomaTreeElemGroup.GetBucket called without Parent`)
		} else {
			return teg.Fault
		}
	}
	return teg.Parent.(SomaTreeBucketeer).GetBucket()
}

func (teg *SomaTreeElemGroup) GetEnvironment() string {
	return teg.Parent.(SomaTreeBucketeer).GetBucket().(SomaTreeBucketeer).GetEnvironment()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
