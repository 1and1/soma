package somatree

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/satori/go.uuid"
)

type SomaTreeElemCluster struct {
	Id              uuid.UUID
	Name            string
	State           string
	Team            uuid.UUID
	Type            string
	Parent          SomaTreeClusterReceiver `json:"-"`
	Fault           *SomaTreeElemFault      `json:"-"`
	Action          chan *Action            `json:"-"`
	PropertyOncall  map[string]SomaTreeProperty
	PropertyService map[string]SomaTreeProperty
	PropertySystem  map[string]SomaTreeProperty
	PropertyCustom  map[string]SomaTreeProperty
	Checks          map[string]SomaTreeCheck
	CheckInstances  map[string][]string
	Instances       map[string]SomaTreeCheckInstance
	Children        map[string]SomaTreeClusterAttacher //`json:"-"`
}

type ClusterSpec struct {
	Id   string
	Name string
	Team string
}

//
// NEW
func NewCluster(spec ClusterSpec) *SomaTreeElemCluster {
	if !specClusterCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	tec := new(SomaTreeElemCluster)
	tec.Id, _ = uuid.FromString(spec.Id)
	tec.Name = spec.Name
	tec.Team, _ = uuid.FromString(spec.Team)
	tec.Type = "cluster"
	tec.State = "floating"
	tec.Parent = nil
	tec.Children = make(map[string]SomaTreeClusterAttacher)
	tec.PropertyOncall = make(map[string]SomaTreeProperty)
	tec.PropertyService = make(map[string]SomaTreeProperty)
	tec.PropertySystem = make(map[string]SomaTreeProperty)
	tec.PropertyCustom = make(map[string]SomaTreeProperty)
	tec.Checks = make(map[string]SomaTreeCheck)

	return tec
}

func (tec SomaTreeElemCluster) Clone() *SomaTreeElemCluster {
	cl := SomaTreeElemCluster{
		Name:  tec.Name,
		State: tec.State,
		Type:  tec.Type,
	}
	cl.Id, _ = uuid.FromString(tec.Id.String())
	cl.Team, _ = uuid.FromString(tec.Team.String())

	f := make(map[string]SomaTreeClusterAttacher, 0)
	for k, child := range tec.Children {
		f[k] = child.CloneCluster()
	}
	cl.Children = f

	pO := make(map[string]SomaTreeProperty)
	for k, prop := range tec.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]SomaTreeProperty)
	for k, prop := range tec.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]SomaTreeProperty)
	for k, prop := range tec.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]SomaTreeProperty)
	for k, prop := range tec.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]SomaTreeCheck)
	for k, chk := range tec.Checks {
		cK[k] = chk.Clone()
	}
	cl.Checks = cK

	cki := make(map[string]SomaTreeCheckInstance)
	for k, chki := range tec.Instances {
		cki[k] = chki.Clone()
	}
	cl.Instances = cki

	ci := make(map[string][]string)
	for k, _ := range tec.CheckInstances {
		for _, str := range tec.CheckInstances[k] {
			t := str
			ci[k] = append(ci[k], t)
		}
	}
	cl.CheckInstances = ci

	return &cl
}

func (tec SomaTreeElemCluster) CloneBucket() SomaTreeBucketAttacher {
	return tec.Clone()
}

func (tec SomaTreeElemCluster) CloneGroup() SomaTreeGroupAttacher {
	return tec.Clone()
}

//
// Interface: SomaTreeBuilder
func (tec *SomaTreeElemCluster) GetID() string {
	return tec.Id.String()
}

func (tec *SomaTreeElemCluster) GetName() string {
	return tec.Name
}

func (tec *SomaTreeElemCluster) GetType() string {
	return tec.Type
}

func (tec *SomaTreeElemCluster) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case *SomaTreeElemBucket:
		tec.setClusterParent(p.(SomaTreeClusterReceiver))
		tec.State = "standalone"
	case *SomaTreeElemGroup:
		tec.setClusterParent(p.(SomaTreeClusterReceiver))
		tec.State = "grouped"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemCluster.setParent`)
	}
}

func (tec *SomaTreeElemCluster) setAction(c chan *Action) {
	tec.Action = c
}

func (tec *SomaTreeElemCluster) setActionDeep(c chan *Action) {
	tec.setAction(c)
	for ch, _ := range tec.Children {
		tec.Children[ch].setActionDeep(c)
	}
}

func (tec *SomaTreeElemCluster) updateParentRecursive(p SomaTreeReceiver) {
	tec.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(str SomaTreeReceiver) {
			defer wg.Done()
			tec.Children[c].updateParentRecursive(str)
		}(tec)
	}
	wg.Wait()
}

// SomaTreeClusterReceiver == can receive Clusters as children
func (tec *SomaTreeElemCluster) setClusterParent(p SomaTreeClusterReceiver) {
	tec.Parent = p
}

func (tec *SomaTreeElemCluster) clearParent() {
	tec.Parent = nil
	tec.State = "floating"
}

func (tec *SomaTreeElemCluster) setFault(f *SomaTreeElemFault) {
	tec.Fault = f
}

func (tec *SomaTreeElemCluster) updateFaultRecursive(f *SomaTreeElemFault) {
	tec.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(ptr *SomaTreeElemFault) {
			defer wg.Done()
			tec.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

//
// Interface: SomaTreeBucketeer
func (tec *SomaTreeElemCluster) GetBucket() SomaTreeReceiver {
	if tec.Parent == nil {
		if tec.Fault == nil {
			panic(`SomaTreeElemCluster.GetBucket called without Parent`)
		} else {
			return tec.Fault
		}
	}
	return tec.Parent.(SomaTreeBucketeer).GetBucket()
}

func (tec *SomaTreeElemCluster) GetEnvironment() string {
	return tec.Parent.(SomaTreeBucketeer).GetBucket().(SomaTreeBucketeer).GetEnvironment()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
