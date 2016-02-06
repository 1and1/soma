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
	Id   uuid.UUID
	Name string
	Team uuid.UUID
}

//
// NEW
func NewCluster(name string, id string) *SomaTreeElemCluster {
	tec := new(SomaTreeElemCluster)
	if id != "" {
		tec.Id, _ = uuid.FromString(id)
	} else {
		tec.Id = uuid.NewV4()
	}
	tec.Name = name
	tec.Type = "cluster"
	tec.State = "floating"
	tec.Children = make(map[string]SomaTreeClusterAttacher)
	tec.PropertyOncall = make(map[string]SomaTreeProperty)
	tec.PropertyService = make(map[string]SomaTreeProperty)
	tec.PropertySystem = make(map[string]SomaTreeProperty)
	tec.PropertyCustom = make(map[string]SomaTreeProperty)
	tec.Checks = make(map[string]SomaTreeCheck)

	return tec
}

func (tec SomaTreeElemCluster) CloneBucket() SomaTreeBucketAttacher {
	for k, child := range tec.Children {
		tec.Children[k] = child.CloneCluster()
	}
	return &tec
}

func (tec SomaTreeElemCluster) CloneGroup() SomaTreeGroupAttacher {
	f := make(map[string]SomaTreeClusterAttacher)
	for k, child := range tec.Children {
		f[k] = child.CloneCluster()
	}
	tec.Children = f
	return &tec
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
