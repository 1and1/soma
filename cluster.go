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
	"sync"


	"github.com/satori/go.uuid"
)

type Cluster struct {
	Id              uuid.UUID
	Name            string
	State           string
	Team            uuid.UUID
	Type            string
	Parent          ClusterReceiver `json:"-"`
	Fault           *Fault          `json:"-"`
	Action          chan *Action    `json:"-"`
	PropertyOncall  map[string]Property
	PropertyService map[string]Property
	PropertySystem  map[string]Property
	PropertyCustom  map[string]Property
	Checks          map[string]Check
	CheckInstances  map[string][]string
	Instances       map[string]CheckInstance
	Children        map[string]ClusterAttacher          `json:"-"`
	loadedInstances map[string]map[string]CheckInstance `json:"-"`
}

type ClusterSpec struct {
	Id   string
	Name string
	Team string
}

//
// NEW
func NewCluster(spec ClusterSpec) *Cluster {
	if !specClusterCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	tec := new(Cluster)
	tec.Id, _ = uuid.FromString(spec.Id)
	tec.Name = spec.Name
	tec.Team, _ = uuid.FromString(spec.Team)
	tec.Type = "cluster"
	tec.State = "floating"
	tec.Parent = nil
	tec.Children = make(map[string]ClusterAttacher)
	tec.PropertyOncall = make(map[string]Property)
	tec.PropertyService = make(map[string]Property)
	tec.PropertySystem = make(map[string]Property)
	tec.PropertyCustom = make(map[string]Property)
	tec.Checks = make(map[string]Check)
	tec.CheckInstances = make(map[string][]string)
	tec.Instances = make(map[string]CheckInstance)
	tec.loadedInstances = make(map[string]map[string]CheckInstance)

	return tec
}

func (tec Cluster) Clone() *Cluster {
	cl := Cluster{
		Name:  tec.Name,
		State: tec.State,
		Type:  tec.Type,
	}
	cl.Id, _ = uuid.FromString(tec.Id.String())
	cl.Team, _ = uuid.FromString(tec.Team.String())

	f := make(map[string]ClusterAttacher, 0)
	for k, child := range tec.Children {
		f[k] = child.CloneCluster()
	}
	cl.Children = f

	pO := make(map[string]Property)
	for k, prop := range tec.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]Property)
	for k, prop := range tec.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]Property)
	for k, prop := range tec.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]Property)
	for k, prop := range tec.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]Check)
	for k, chk := range tec.Checks {
		cK[k] = chk.clone()
	}
	cl.Checks = cK

	cki := make(map[string]CheckInstance)
	for k, chki := range tec.Instances {
		cki[k] = chki.Clone()
	}
	cl.Instances = cki
	cl.loadedInstances = make(map[string]map[string]CheckInstance)

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

func (tec Cluster) CloneBucket() BucketAttacher {
	return tec.Clone()
}

func (tec Cluster) CloneGroup() GroupAttacher {
	return tec.Clone()
}

//
// Interface: Builder
func (tec *Cluster) GetID() string {
	return tec.Id.String()
}

func (tec *Cluster) GetName() string {
	return tec.Name
}

func (tec *Cluster) GetType() string {
	return tec.Type
}

func (tec *Cluster) setParent(p Receiver) {
	switch p.(type) {
	case *Bucket:
		tec.setClusterParent(p.(ClusterReceiver))
		tec.State = "standalone"
	case *Group:
		tec.setClusterParent(p.(ClusterReceiver))
		tec.State = "grouped"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`Cluster.setParent`)
	}
}

func (tec *Cluster) setAction(c chan *Action) {
	tec.Action = c
}

func (tec *Cluster) setActionDeep(c chan *Action) {
	tec.setAction(c)
	for ch, _ := range tec.Children {
		tec.Children[ch].setActionDeep(c)
	}
}

func (tec *Cluster) updateParentRecursive(p Receiver) {
	tec.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(str Receiver) {
			defer wg.Done()
			tec.Children[c].updateParentRecursive(str)
		}(tec)
	}
	wg.Wait()
}

// ClusterReceiver == can receive Clusters as children
func (tec *Cluster) setClusterParent(p ClusterReceiver) {
	tec.Parent = p
}

func (tec *Cluster) clearParent() {
	tec.Parent = nil
	tec.State = "floating"
}

func (tec *Cluster) setFault(f *Fault) {
	tec.Fault = f
}

func (tec *Cluster) updateFaultRecursive(f *Fault) {
	tec.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(ptr *Fault) {
			defer wg.Done()
			tec.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

//
// Interface: Bucketeer
func (tec *Cluster) GetBucket() Receiver {
	if tec.Parent == nil {
		if tec.Fault == nil {
			panic(`Cluster.GetBucket called without Parent`)
		} else {
			return tec.Fault
		}
	}
	return tec.Parent.(Bucketeer).GetBucket()
}

func (tec *Cluster) GetRepository() string {
	return tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
}

func (tec *Cluster) GetRepositoryName() string {
	return tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepositoryName()
}

func (tec *Cluster) GetEnvironment() string {
	return tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetEnvironment()
}

//
//
func (tec *Cluster) ComputeCheckInstances() {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			tec.Children[c].ComputeCheckInstances()
		}()
	}
	wg.Wait()
	tec.updateCheckInstances()
}

//
//
func (tec *Cluster) ClearLoadInfo() {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			tec.Children[c].ClearLoadInfo()
		}()
	}
	wg.Wait()
	tec.loadedInstances = map[string]map[string]CheckInstance{}
}

//
//
func (tec *Cluster) export() proto.Cluster {
	bucket := tec.Parent.(Bucketeer).GetBucket()
	return proto.Cluster{
		Id:          tec.Id.String(),
		Name:        tec.Name,
		BucketId:    bucket.(Builder).GetID(),
		ObjectState: tec.State,
		TeamId:      tec.Team.String(),
	}
}

func (tec *Cluster) actionCreate() {
	tec.Action <- &Action{
		Action:  "create",
		Type:    tec.Type,
		Cluster: tec.export(),
	}
}

func (tec *Cluster) actionUpdate() {
	tec.Action <- &Action{
		Action:  "update",
		Type:    tec.Type,
		Cluster: tec.export(),
	}
}

func (tec *Cluster) actionDelete() {
	tec.Action <- &Action{
		Action:  "delete",
		Type:    tec.Type,
		Cluster: tec.export(),
	}
}

func (tec *Cluster) actionMemberNew(a Action) {
	a.Action = "member_new"
	a.Type = tec.Type
	a.Cluster = tec.export()

	tec.Action <- &a
}

func (tec *Cluster) actionMemberRemoved(a Action) {
	a.Action = "member_removed"
	a.Type = tec.Type
	a.Cluster = tec.export()

	tec.Action <- &a
}

//
func (tec *Cluster) actionPropertyNew(a Action) {
	a.Action = `property_new`
	tec.actionProperty(a)
}

func (tec *Cluster) actionPropertyUpdate(a Action) {
	a.Action = `property_update`
	tec.actionProperty(a)
}

func (tec *Cluster) actionPropertyDelete(a Action) {
	a.Action = `property_delete`
	tec.actionProperty(a)
}

func (tec *Cluster) actionProperty(a Action) {
	a.Type = tec.Type
	a.Cluster = tec.export()
	a.Property.RepositoryId = tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Property.BucketId = tec.Parent.(Bucketeer).GetBucket().(Builder).GetID()

	switch a.Property.Type {
	case "custom":
		a.Property.Custom.RepositoryId = a.Property.RepositoryId
	case "service":
		a.Property.Service.TeamId = tec.Team.String()
	}

	tec.Action <- &a
}

//
func (tec *Cluster) actionCheckNew(a Action) {
	a.Check.RepositoryId = tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = tec.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	tec.actionDispatch("check_new", a)
}

func (tec *Cluster) actionCheckRemoved(a Action) {
	a.Check.RepositoryId = tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = tec.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	tec.actionDispatch(`check_removed`, a)
}

func (tec *Cluster) setupCheckAction(c Check) Action {
	return c.MakeAction()
}

func (tec *Cluster) actionCheckInstanceCreate(a Action) {
	tec.actionDispatch("check_instance_create", a)
}

func (tec *Cluster) actionCheckInstanceUpdate(a Action) {
	tec.actionDispatch("check_instance_update", a)
}

func (tec *Cluster) actionCheckInstanceDelete(a Action) {
	tec.actionDispatch("check_instance_delete", a)
}

func (tec *Cluster) actionDispatch(action string, a Action) {
	a.Action = action
	a.Type = tec.Type
	a.Cluster = tec.export()

	tec.Action <- &a
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
