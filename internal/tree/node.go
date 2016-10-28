/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"fmt"
	"reflect"

	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type Node struct {
	Id              uuid.UUID
	Name            string
	AssetId         uint64
	Team            uuid.UUID
	ServerId        uuid.UUID
	State           string
	Online          bool
	Deleted         bool
	Type            string
	Parent          NodeReceiver `json:"-"`
	Fault           *Fault       `json:"-"`
	Action          chan *Action `json:"-"`
	PropertyOncall  map[string]Property
	PropertyService map[string]Property
	PropertySystem  map[string]Property
	PropertyCustom  map[string]Property
	Checks          map[string]Check
	CheckInstances  map[string][]string
	Instances       map[string]CheckInstance
	loadedInstances map[string]map[string]CheckInstance
	hasUpdate       bool
	log             *log.Logger
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
func NewNode(spec NodeSpec) *Node {
	if !specNodeCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	ten := new(Node)
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
	ten.PropertyOncall = make(map[string]Property)
	ten.PropertyService = make(map[string]Property)
	ten.PropertySystem = make(map[string]Property)
	ten.PropertyCustom = make(map[string]Property)
	ten.Checks = make(map[string]Check)
	ten.CheckInstances = make(map[string][]string)
	ten.Instances = make(map[string]CheckInstance)
	ten.loadedInstances = make(map[string]map[string]CheckInstance)

	return ten
}

func (ten Node) Clone() *Node {
	cl := Node{
		Name:    ten.Name,
		State:   ten.State,
		Online:  ten.Online,
		Deleted: ten.Deleted,
		Type:    ten.Type,
		log:     ten.log,
	}
	cl.Id, _ = uuid.FromString(ten.Id.String())
	cl.AssetId = ten.AssetId
	cl.Team, _ = uuid.FromString(ten.Team.String())
	cl.ServerId, _ = uuid.FromString(ten.ServerId.String())

	pO := make(map[string]Property)
	for k, prop := range ten.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]Property)
	for k, prop := range ten.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]Property)
	for k, prop := range ten.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]Property)
	for k, prop := range ten.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]Check)
	for k, chk := range ten.Checks {
		cK[k] = chk.Clone()
	}
	cl.Checks = cK

	cki := make(map[string]CheckInstance)
	for k, chki := range ten.Instances {
		cki[k] = chki.Clone()
	}
	cl.Instances = cki
	cl.loadedInstances = make(map[string]map[string]CheckInstance)

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

func (ten Node) CloneBucket() BucketAttacher {
	return ten.Clone()
}

func (ten Node) CloneGroup() GroupAttacher {
	return ten.Clone()
}

func (ten Node) CloneCluster() ClusterAttacher {
	return ten.Clone()
}

//
// Interface:
func (ten *Node) GetID() string {
	return ten.Id.String()
}

func (ten *Node) GetName() string {
	return ten.Name
}

func (ten *Node) GetType() string {
	return ten.Type
}

func (ten *Node) setParent(p Receiver) {
	switch p.(type) {
	case *Bucket:
		ten.setNodeParent(p.(NodeReceiver))
		ten.State = "standalone"
	case *Group:
		ten.setNodeParent(p.(NodeReceiver))
		ten.State = "grouped"
	case *Cluster:
		ten.setNodeParent(p.(NodeReceiver))
		ten.State = "clustered"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`Node.setParent`)
	}
}

func (ten *Node) setAction(c chan *Action) {
	ten.Action = c
}

func (ten *Node) setActionDeep(c chan *Action) {
	ten.setAction(c)
}

func (n *Node) setLog(newlog *log.Logger) {
	n.log = newlog
}

func (n *Node) setLoggerDeep(newlog *log.Logger) {
	n.setLog(newlog)
}

func (ten *Node) updateParentRecursive(p Receiver) {
	ten.setParent(p)
}

func (ten *Node) setNodeParent(p NodeReceiver) {
	ten.Parent = p
}

func (ten *Node) clearParent() {
	ten.Parent = nil
	ten.State = "floating"
}

func (ten *Node) setFault(f *Fault) {
	ten.Fault = f
}

func (ten *Node) updateFaultRecursive(f *Fault) {
	ten.setFault(f)
}

//
//
func (ten *Node) ComputeCheckInstances() {
	ten.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s",
		ten.Parent.(Bucketeer).GetRepositoryName(),
		`ComputeCheckInstances`,
		`node`,
		ten.Id.String(),
	)
	ten.updateCheckInstances()
}

//
//
func (ten *Node) ClearLoadInfo() {
	ten.loadedInstances = map[string]map[string]CheckInstance{}
}

//
//
func (ten *Node) export() proto.Node {
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

func (ten *Node) actionUpdate() {
	ten.Action <- &Action{
		Action: "update",
		Type:   ten.Type,
		Node:   ten.export(),
	}
}

func (ten *Node) actionDelete() {
	ten.Action <- &Action{
		Action: "delete",
		Type:   ten.Type,
		Node:   ten.export(),
	}
}

//
func (ten *Node) actionPropertyNew(a Action) {
	a.Action = `property_new`
	ten.actionProperty(a)
}

func (ten *Node) actionPropertyUpdate(a Action) {
	a.Action = `property_update`
	ten.actionProperty(a)
}

func (ten *Node) actionPropertyDelete(a Action) {
	a.Action = `property_delete`
	ten.actionProperty(a)
}

func (ten *Node) actionProperty(a Action) {
	a.Type = ten.Type
	a.Node = ten.export()
	a.Property.RepositoryId = ten.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Property.BucketId = ten.Parent.(Bucketeer).GetBucket().(Builder).GetID()

	switch a.Property.Type {
	case "custom":
		a.Property.Custom.RepositoryId = a.Property.RepositoryId
	case "service":
		a.Property.Service.TeamId = ten.Team.String()
	}

	ten.Action <- &a
}

//
func (ten *Node) actionCheckNew(a Action) {
	a.Check.RepositoryId = ten.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = ten.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	ten.actionDispatch("check_new", a)
}

func (ten *Node) actionCheckRemoved(a Action) {
	a.Check.RepositoryId = ten.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = ten.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	ten.actionDispatch(`check_removed`, a)
}

func (ten *Node) setupCheckAction(c Check) Action {
	return c.MakeAction()
}

func (ten *Node) actionCheckInstanceCreate(a Action) {
	ten.actionDispatch("check_instance_create", a)
}

func (ten *Node) actionCheckInstanceUpdate(a Action) {
	ten.actionDispatch("check_instance_update", a)
}

func (ten *Node) actionCheckInstanceDelete(a Action) {
	ten.actionDispatch("check_instance_delete", a)
}

func (ten *Node) actionDispatch(action string, a Action) {
	a.Action = action
	a.Type = ten.Type
	a.Node = ten.export()

	ten.Action <- &a
}

func (ten *Node) repositoryName() string {
	return ten.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepositoryName()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
