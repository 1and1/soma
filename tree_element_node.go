package somatree

import (
	"fmt"
	"reflect"

	"github.com/satori/go.uuid"
)

type SomaTreeElemNode struct {
	Id              uuid.UUID
	Name            string
	AssetId         uuid.UUID
	Team            uuid.UUID
	ServerId        uuid.UUID
	State           string
	Online          bool
	Deleted         bool
	Type            string
	Parent          SomaTreeNodeReceiver `json:"-"`
	Fault           *SomaTreeElemFault   `json:"-"`
	Action          chan *Action         `json:"-"`
	PropertyOncall  map[string]SomaTreeProperty
	PropertyService map[string]SomaTreeProperty
	PropertySystem  map[string]SomaTreeProperty
	PropertyCustom  map[string]SomaTreeProperty
	Checks          map[string]SomaTreeCheck
	CheckInstances  map[string][]string
	Instances       map[string]SomaTreeCheckInstance
}

type NodeSpec struct {
	Id       uuid.UUID
	Name     string
	AssetId  uuid.UUID
	Team     uuid.UUID
	ServerId uuid.UUID
	State    string
}

//
// NEW
func NewNode(name string, id string) *SomaTreeElemNode {
	ten := new(SomaTreeElemNode)
	if id != "" {
		ten.Id, _ = uuid.FromString(id)
	} else {
		ten.Id = uuid.NewV4()
	}
	ten.Name = name
	ten.Type = "node"
	ten.State = "floating"
	ten.Parent = nil
	ten.PropertyOncall = make(map[string]SomaTreeProperty)
	ten.PropertyService = make(map[string]SomaTreeProperty)
	ten.PropertySystem = make(map[string]SomaTreeProperty)
	ten.PropertyCustom = make(map[string]SomaTreeProperty)
	ten.Checks = make(map[string]SomaTreeCheck)
	ten.CheckInstances = make(map[string][]string)
	ten.Instances = make(map[string]SomaTreeCheckInstance)

	return ten
}

func (ten SomaTreeElemNode) Clone() *SomaTreeElemNode {
	cl := SomaTreeElemNode{
		Name:    ten.Name,
		State:   ten.State,
		Online:  ten.Online,
		Deleted: ten.Deleted,
		Type:    ten.Type,
	}
	cl.Id, _ = uuid.FromString(ten.Id.String())
	cl.AssetId, _ = uuid.FromString(ten.AssetId.String())
	cl.Team, _ = uuid.FromString(ten.Team.String())
	cl.ServerId, _ = uuid.FromString(ten.ServerId.String())

	pO := make(map[string]SomaTreeProperty)
	for k, prop := range ten.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]SomaTreeProperty)
	for k, prop := range ten.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]SomaTreeProperty)
	for k, prop := range ten.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]SomaTreeProperty)
	for k, prop := range ten.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]SomaTreeCheck)
	for k, chk := range ten.Checks {
		cK[k] = chk.Clone()
	}
	cl.Checks = cK

	cki := make(map[string]SomaTreeCheckInstance)
	for k, chki := range ten.Instances {
		cki[k] = chki.Clone()
	}
	cl.Instances = cki

	ci := make(map[string][]string)
	for k, _ := range ten.CheckInstances {
		for _, str := range ten.CheckInstances[k] {
			t := str
			ci[k] = append(ci[k], t)
		}
	}
	cl.CheckInstances = ci

	return &cl
}

func (ten SomaTreeElemNode) CloneBucket() SomaTreeBucketAttacher {
	return ten.Clone()
}

func (ten SomaTreeElemNode) CloneGroup() SomaTreeGroupAttacher {
	return ten.Clone()
}

func (ten SomaTreeElemNode) CloneCluster() SomaTreeClusterAttacher {
	return ten.Clone()
}

//
// Interface:
func (ten *SomaTreeElemNode) GetID() string {
	return ten.Id.String()
}

func (ten *SomaTreeElemNode) GetName() string {
	return ten.Name
}

func (ten *SomaTreeElemNode) GetType() string {
	return ten.Type
}

func (ten *SomaTreeElemNode) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case *SomaTreeElemBucket:
		ten.setNodeParent(p.(SomaTreeNodeReceiver))
		ten.State = "standalone"
	case *SomaTreeElemGroup:
		ten.setNodeParent(p.(SomaTreeNodeReceiver))
		ten.State = "grouped"
	case *SomaTreeElemCluster:
		ten.setNodeParent(p.(SomaTreeNodeReceiver))
		ten.State = "clustered"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemNode.setParent`)
	}
}

func (ten *SomaTreeElemNode) setAction(c chan *Action) {
	ten.Action = c
}

func (ten *SomaTreeElemNode) setActionDeep(c chan *Action) {
	ten.setAction(c)
}

func (ten *SomaTreeElemNode) updateParentRecursive(p SomaTreeReceiver) {
	ten.setParent(p)
}

func (ten *SomaTreeElemNode) setNodeParent(p SomaTreeNodeReceiver) {
	ten.Parent = p
}

func (ten *SomaTreeElemNode) clearParent() {
	ten.Parent = nil
	ten.State = "floating"
}

func (ten *SomaTreeElemNode) setFault(f *SomaTreeElemFault) {
	ten.Fault = f
}

func (ten *SomaTreeElemNode) updateFaultRecursive(f *SomaTreeElemFault) {
	ten.setFault(f)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
