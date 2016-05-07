package somatree

import (
	"fmt"
	"reflect"


	"github.com/satori/go.uuid"
)

type SomaTreeElemNode struct {
	Id              uuid.UUID
	Name            string
	AssetId         uint64
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
	Checks          map[string]Check
	CheckInstances  map[string][]string
	Instances       map[string]CheckInstance
}

type NodeSpec struct {
	Id       string
	AssetId  uint64
	Name     string
	Team     string
	ServerId string
	Online   bool
	Deleted  bool
}

//
// NEW
func NewNode(spec NodeSpec) *SomaTreeElemNode {
	if !specNodeCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	ten := new(SomaTreeElemNode)
	ten.Id, _ = uuid.FromString(spec.Id)
	ten.Name = spec.Name
	ten.AssetId = spec.AssetId
	ten.Team, _ = uuid.FromString(spec.Team)
	ten.ServerId, _ = uuid.FromString(spec.ServerId)
	ten.Online = spec.Online
	ten.Deleted = spec.Deleted
	ten.Type = "node"
	ten.State = "floating"
	ten.Parent = nil
	ten.PropertyOncall = make(map[string]SomaTreeProperty)
	ten.PropertyService = make(map[string]SomaTreeProperty)
	ten.PropertySystem = make(map[string]SomaTreeProperty)
	ten.PropertyCustom = make(map[string]SomaTreeProperty)
	ten.Checks = make(map[string]Check)
	ten.CheckInstances = make(map[string][]string)
	ten.Instances = make(map[string]CheckInstance)

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
	cl.AssetId = ten.AssetId
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

	cK := make(map[string]Check)
	for k, chk := range ten.Checks {
		cK[k] = chk.clone()
	}
	cl.Checks = cK

	cki := make(map[string]CheckInstance)
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

//
//
func (ten *SomaTreeElemNode) ComputeCheckInstances() {
	ten.updateCheckInstances()
}

//
//
func (ten *SomaTreeElemNode) export() proto.Node {
	bucket := ten.Parent.(Bucketeer).GetBucket()
	return proto.Node{
		Id:        ten.Id.String(),
		AssetId:   ten.AssetId,
		Name:      ten.Name,
		TeamId:    ten.Team.String(),
		ServerId:  ten.ServerId.String(),
		State:     ten.State,
		IsOnline:  ten.Online,
		IsDeleted: ten.Deleted,
		Config: &proto.NodeConfig{
			BucketId: bucket.(Builder).GetID(),
		},
	}
}

func (ten *SomaTreeElemNode) actionUpdate() {
	ten.Action <- &Action{
		Action: "update",
		Type:   ten.Type,
		Node:   ten.export(),
	}
}

func (ten *SomaTreeElemNode) actionDelete() {
	ten.Action <- &Action{
		Action: "delete",
		Type:   ten.Type,
		Node:   ten.export(),
	}
}

func (ten *SomaTreeElemNode) actionPropertyNew(a Action) {
	a.Property.RepositoryId = ten.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Property.BucketId = ten.Parent.(Bucketeer).GetBucket().(Builder).GetID()

	switch a.Property.Type {
	case "custom":
		a.Property.Custom.RepositoryId = a.Property.RepositoryId
	case "service":
		a.Property.Service.TeamId = ten.Team.String()
	}

	ten.actionDispatch("property_new", a)
}

//
func (ten *SomaTreeElemNode) setupPropertyAction(p SomaTreeProperty) Action {
	return p.MakeAction()
}

//
func (ten *SomaTreeElemNode) actionCheckNew(a Action) {
	a.Check.RepositoryId = ten.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = ten.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	ten.actionDispatch("check_new", a)
}

func (ten *SomaTreeElemNode) setupCheckAction(c Check) Action {
	return c.MakeAction()
}

func (ten *SomaTreeElemNode) actionCheckInstanceCreate(a Action) {
	ten.actionDispatch("check_instance_create", a)
}

func (ten *SomaTreeElemNode) actionCheckInstanceUpdate(a Action) {
	ten.actionDispatch("check_instance_update", a)
}

func (ten *SomaTreeElemNode) actionCheckInstanceDelete(a Action) {
	ten.actionDispatch("check_instance_delete", a)
}

func (ten *SomaTreeElemNode) actionDispatch(action string, a Action) {
	a.Action = action
	a.Type = ten.Type
	a.Node = ten.export()

	ten.Action <- &a
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
