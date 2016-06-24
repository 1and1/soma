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
	PropertyOncall  map[string]Property
	PropertyService map[string]Property
	PropertySystem  map[string]Property
	PropertyCustom  map[string]Property
	Checks          map[string]Check
	CheckInstances  map[string][]string
	Instances       map[string]CheckInstance
	Children        map[string]SomaTreeClusterAttacher  `json:"-"`
	loadedInstances map[string]map[string]CheckInstance `json:"-"`
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

func (tec SomaTreeElemCluster) CloneBucket() SomaTreeBucketAttacher {
	return tec.Clone()
}

func (tec SomaTreeElemCluster) CloneGroup() SomaTreeGroupAttacher {
	return tec.Clone()
}

//
// Interface: Builder
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
// Interface: Bucketeer
func (tec *SomaTreeElemCluster) GetBucket() SomaTreeReceiver {
	if tec.Parent == nil {
		if tec.Fault == nil {
			panic(`SomaTreeElemCluster.GetBucket called without Parent`)
		} else {
			return tec.Fault
		}
	}
	return tec.Parent.(Bucketeer).GetBucket()
}

func (tec *SomaTreeElemCluster) GetRepository() string {
	return tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
}

func (tec *SomaTreeElemCluster) GetRepositoryName() string {
	return tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepositoryName()
}

func (tec *SomaTreeElemCluster) GetEnvironment() string {
	return tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetEnvironment()
}

//
//
func (tec *SomaTreeElemCluster) ComputeCheckInstances() {
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
func (tec *SomaTreeElemCluster) ClearLoadInfo() {
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
func (tec *SomaTreeElemCluster) export() proto.Cluster {
	bucket := tec.Parent.(Bucketeer).GetBucket()
	return proto.Cluster{
		Id:          tec.Id.String(),
		Name:        tec.Name,
		BucketId:    bucket.(Builder).GetID(),
		ObjectState: tec.State,
		TeamId:      tec.Team.String(),
	}
}

func (tec *SomaTreeElemCluster) actionCreate() {
	tec.Action <- &Action{
		Action:  "create",
		Type:    tec.Type,
		Cluster: tec.export(),
	}
}

func (tec *SomaTreeElemCluster) actionUpdate() {
	tec.Action <- &Action{
		Action:  "update",
		Type:    tec.Type,
		Cluster: tec.export(),
	}
}

func (tec *SomaTreeElemCluster) actionDelete() {
	tec.Action <- &Action{
		Action:  "delete",
		Type:    tec.Type,
		Cluster: tec.export(),
	}
}

func (tec *SomaTreeElemCluster) actionMemberNew(a Action) {
	a.Action = "member_new"
	a.Type = tec.Type
	a.Cluster = tec.export()

	tec.Action <- &a
}

func (tec *SomaTreeElemCluster) actionMemberRemoved(a Action) {
	a.Action = "member_removed"
	a.Type = tec.Type
	a.Cluster = tec.export()

	tec.Action <- &a
}

func (tec *SomaTreeElemCluster) actionPropertyNew(a Action) {
	a.Property.RepositoryId = tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Property.BucketId = tec.Parent.(Bucketeer).GetBucket().(Builder).GetID()

	switch a.Property.Type {
	case "custom":
		a.Property.Custom.RepositoryId = a.Property.RepositoryId
	case "service":
		a.Property.Service.TeamId = tec.Team.String()
	}

	tec.actionDispatch("property_new", a)
}

//
func (tec *SomaTreeElemCluster) setupPropertyAction(p Property) Action {
	return p.MakeAction()
}

//
func (tec *SomaTreeElemCluster) actionCheckNew(a Action) {
	a.Check.RepositoryId = tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = tec.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	tec.actionDispatch("check_new", a)
}

func (tec *SomaTreeElemCluster) actionCheckRemoved(a Action) {
	a.Check.RepositoryId = tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = tec.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	tec.actionDispatch(`check_removed`, a)
}

func (tec *SomaTreeElemCluster) setupCheckAction(c Check) Action {
	return c.MakeAction()
}

func (tec *SomaTreeElemCluster) actionCheckInstanceCreate(a Action) {
	tec.actionDispatch("check_instance_create", a)
}

func (tec *SomaTreeElemCluster) actionCheckInstanceUpdate(a Action) {
	tec.actionDispatch("check_instance_update", a)
}

func (tec *SomaTreeElemCluster) actionCheckInstanceDelete(a Action) {
	tec.actionDispatch("check_instance_delete", a)
}

func (tec *SomaTreeElemCluster) actionDispatch(action string, a Action) {
	a.Action = action
	a.Type = tec.Type
	a.Cluster = tec.export()

	tec.Action <- &a
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
