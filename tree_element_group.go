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

type GroupSpec struct {
	Id   string
	Name string
	Team string
}

//
// NEW
func NewGroup(spec GroupSpec) *SomaTreeElemGroup {
	if !specGroupCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	teg := new(SomaTreeElemGroup)
	teg.Id, _ = uuid.FromString(spec.Id)
	teg.Name = spec.Name
	teg.Team, _ = uuid.FromString(spec.Team)
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

func (teg SomaTreeElemGroup) Clone() *SomaTreeElemGroup {
	cl := SomaTreeElemGroup{
		Name:  teg.Name,
		State: teg.State,
		Type:  teg.Type,
	}
	cl.Id, _ = uuid.FromString(teg.Id.String())

	f := make(map[string]SomaTreeGroupAttacher, 0)
	for k, child := range teg.Children {
		f[k] = child.CloneGroup()
	}
	cl.Children = f

	pO := make(map[string]SomaTreeProperty)
	for k, prop := range teg.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]SomaTreeProperty)
	for k, prop := range teg.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]SomaTreeProperty)
	for k, prop := range teg.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]SomaTreeProperty)
	for k, prop := range teg.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]SomaTreeCheck)
	for k, chk := range teg.Checks {
		cK[k] = chk.Clone()
	}
	cl.Checks = cK

	cki := make(map[string]SomaTreeCheckInstance)
	for k, chki := range teg.Instances {
		cki[k] = chki.Clone()
	}
	cl.Instances = cki

	ci := make(map[string][]string)
	for k, _ := range teg.CheckInstances {
		for _, str := range teg.CheckInstances[k] {
			t := str
			ci[k] = append(ci[k], t)
		}
	}
	cl.CheckInstances = ci

	return &cl
}

func (teg SomaTreeElemGroup) CloneBucket() SomaTreeBucketAttacher {
	return teg.Clone()
}

func (teg SomaTreeElemGroup) CloneGroup() SomaTreeGroupAttacher {
	return teg.Clone()
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

func (teg *SomaTreeElemGroup) setActionDeep(c chan *Action) {
	teg.setAction(c)
	for ch, _ := range teg.Children {
		teg.Children[ch].setActionDeep(c)
	}
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
