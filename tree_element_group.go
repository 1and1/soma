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
	PropertyOncall  map[string]Property
	PropertyService map[string]Property
	PropertySystem  map[string]Property
	PropertyCustom  map[string]Property
	Checks          map[string]Check
	CheckInstances  map[string][]string
	Instances       map[string]CheckInstance
	Children        map[string]SomaTreeGroupAttacher    `json:"-"`
	loadedInstances map[string]map[string]CheckInstance `json:"-"`
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
	teg.PropertyOncall = make(map[string]Property)
	teg.PropertyService = make(map[string]Property)
	teg.PropertySystem = make(map[string]Property)
	teg.PropertyCustom = make(map[string]Property)
	teg.Checks = make(map[string]Check)
	teg.CheckInstances = make(map[string][]string)
	teg.Instances = make(map[string]CheckInstance)
	teg.loadedInstances = make(map[string]map[string]CheckInstance)

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

	pO := make(map[string]Property)
	for k, prop := range teg.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]Property)
	for k, prop := range teg.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]Property)
	for k, prop := range teg.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]Property)
	for k, prop := range teg.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]Check)
	for k, chk := range teg.Checks {
		cK[k] = chk.clone()
	}
	cl.Checks = cK

	cki := make(map[string]CheckInstance)
	for k, chki := range teg.Instances {
		cki[k] = chki.Clone()
	}
	cl.Instances = cki
	cl.loadedInstances = make(map[string]map[string]CheckInstance)

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
// Interface: Builder
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
// Interface: Bucketeer
func (teg *SomaTreeElemGroup) GetBucket() SomaTreeReceiver {
	if teg.Parent == nil {
		if teg.Fault == nil {
			panic(`SomaTreeElemGroup.GetBucket called without Parent`)
		} else {
			return teg.Fault
		}
	}
	return teg.Parent.(Bucketeer).GetBucket()
}

func (teg *SomaTreeElemGroup) GetRepository() string {
	return teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
}

func (teg *SomaTreeElemGroup) GetRepositoryName() string {
	return teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepositoryName()
}

func (teg *SomaTreeElemGroup) GetEnvironment() string {
	return teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetEnvironment()
}

//
//
func (teg *SomaTreeElemGroup) ComputeCheckInstances() {
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			teg.Children[c].ComputeCheckInstances()
		}()
	}
	wg.Wait()
	teg.updateCheckInstances()
}

//
//
func (teg *SomaTreeElemGroup) ClearLoadInfo() {
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			teg.Children[c].ClearLoadInfo()
		}()
	}
	wg.Wait()
	teg.loadedInstances = map[string]map[string]CheckInstance{}
}

//
//
func (teg *SomaTreeElemGroup) export() proto.Group {
	bucket := teg.Parent.(Bucketeer).GetBucket()
	return proto.Group{
		Id:          teg.Id.String(),
		Name:        teg.Name,
		BucketId:    bucket.(Builder).GetID(),
		ObjectState: teg.State,
		TeamId:      teg.Team.String(),
	}
}

func (teg *SomaTreeElemGroup) actionCreate() {
	teg.Action <- &Action{
		Action: "create",
		Type:   teg.Type,
		Group:  teg.export(),
	}
}

func (teg *SomaTreeElemGroup) actionUpdate() {
	teg.Action <- &Action{
		Action: "update",
		Type:   teg.Type,
		Group:  teg.export(),
	}
}

func (teg *SomaTreeElemGroup) actionDelete() {
	teg.Action <- &Action{
		Action: "delete",
		Type:   teg.Type,
		Group:  teg.export(),
	}
}

func (teg *SomaTreeElemGroup) actionMemberNew(a Action) {
	a.Action = "member_new"
	a.Type = teg.Type
	a.Group = teg.export()

	teg.Action <- &a
}

func (teg *SomaTreeElemGroup) actionMemberRemoved(a Action) {
	a.Action = "member_removed"
	a.Type = teg.Type
	a.Group = teg.export()

	teg.Action <- &a
}

//
func (teg *SomaTreeElemGroup) actionPropertyNew(a Action) {
	a.Action = `property_new`
	teg.actionProperty(a)
}

func (teg *SomaTreeElemGroup) actionPropertyUpdate(a Action) {
	a.Action = `property_update`
	teg.actionProperty(a)
}

func (teg *SomaTreeElemGroup) actionPropertyDelete(a Action) {
	a.Action = `property_delete`
	teg.actionProperty(a)
}

func (teg *SomaTreeElemGroup) actionProperty(a Action) {
	a.Type = teg.Type
	a.Group = teg.export()
	a.Property.RepositoryId = teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Property.BucketId = teg.Parent.(Bucketeer).GetBucket().(Builder).GetID()

	switch a.Property.Type {
	case "custom":
		a.Property.Custom.RepositoryId = a.Property.RepositoryId
	case "service":
		a.Property.Service.TeamId = teg.Team.String()
	}

	teg.Action <- &a
}

//
func (teg *SomaTreeElemGroup) actionCheckNew(a Action) {
	a.Check.RepositoryId = teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = teg.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	teg.actionDispatch("check_new", a)
}

func (teg *SomaTreeElemGroup) actionCheckRemoved(a Action) {
	a.Check.RepositoryId = teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = teg.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	teg.actionDispatch(`check_removed`, a)
}

func (teg *SomaTreeElemGroup) setupCheckAction(c Check) Action {
	return c.MakeAction()
}

func (teg *SomaTreeElemGroup) actionCheckInstanceCreate(a Action) {
	teg.actionDispatch("check_instance_create", a)
}

func (teg *SomaTreeElemGroup) actionCheckInstanceUpdate(a Action) {
	teg.actionDispatch("check_instance_update", a)
}

func (teg *SomaTreeElemGroup) actionCheckInstanceDelete(a Action) {
	teg.actionDispatch("check_instance_delete", a)
}

func (teg *SomaTreeElemGroup) actionDispatch(action string, a Action) {
	a.Action = action
	a.Type = teg.Type
	a.Group = teg.export()

	teg.Action <- &a
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
