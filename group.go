package tree

import (
	"fmt"
	"reflect"
	"sync"


	"github.com/satori/go.uuid"
)

type Group struct {
	Id              uuid.UUID
	Name            string
	State           string
	Team            uuid.UUID
	Type            string
	Parent          GroupReceiver `json:"-"`
	Fault           *Fault        `json:"-"`
	Action          chan *Action  `json:"-"`
	PropertyOncall  map[string]Property
	PropertyService map[string]Property
	PropertySystem  map[string]Property
	PropertyCustom  map[string]Property
	Checks          map[string]Check
	CheckInstances  map[string][]string
	Instances       map[string]CheckInstance
	Children        map[string]GroupAttacher            `json:"-"`
	loadedInstances map[string]map[string]CheckInstance `json:"-"`
}

type GroupSpec struct {
	Id   string
	Name string
	Team string
}

//
// NEW
func NewGroup(spec GroupSpec) *Group {
	if !specGroupCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	teg := new(Group)
	teg.Id, _ = uuid.FromString(spec.Id)
	teg.Name = spec.Name
	teg.Team, _ = uuid.FromString(spec.Team)
	teg.Type = "group"
	teg.State = "floating"
	teg.Parent = nil
	teg.Children = make(map[string]GroupAttacher)
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

func (teg Group) Clone() *Group {
	cl := Group{
		Name:  teg.Name,
		State: teg.State,
		Type:  teg.Type,
	}
	cl.Id, _ = uuid.FromString(teg.Id.String())

	f := make(map[string]GroupAttacher, 0)
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

func (teg Group) CloneBucket() BucketAttacher {
	return teg.Clone()
}

func (teg Group) CloneGroup() GroupAttacher {
	return teg.Clone()
}

//
// Interface: Builder
func (teg *Group) GetID() string {
	return teg.Id.String()
}

func (teg *Group) GetName() string {
	return teg.Name
}

func (teg *Group) GetType() string {
	return teg.Type
}

func (teg *Group) setParent(p Receiver) {
	switch p.(type) {
	case *Bucket:
		teg.setGroupParent(p.(GroupReceiver))
		teg.State = "standalone"
	case *Group:
		teg.setGroupParent(p.(GroupReceiver))
		teg.State = "grouped"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`Group.setParent`)
	}
}

func (teg *Group) setAction(c chan *Action) {
	teg.Action = c
}

func (teg *Group) setActionDeep(c chan *Action) {
	teg.setAction(c)
	for ch, _ := range teg.Children {
		teg.Children[ch].setActionDeep(c)
	}
}

// GroupReceiver == can receive Groups as children
func (teg *Group) setGroupParent(p GroupReceiver) {
	teg.Parent = p
}

func (teg *Group) updateParentRecursive(p Receiver) {
	teg.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func(str Receiver) {
			defer wg.Done()
			teg.Children[c].updateParentRecursive(str)
		}(teg)
	}
	wg.Wait()
}

func (teg *Group) clearParent() {
	teg.Parent = nil
	teg.State = "floating"
}

func (teg *Group) setFault(f *Fault) {
	teg.Fault = f
}

func (teg *Group) updateFaultRecursive(f *Fault) {
	teg.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func(ptr *Fault) {
			defer wg.Done()
			teg.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

//
// Interface: Bucketeer
func (teg *Group) GetBucket() Receiver {
	if teg.Parent == nil {
		if teg.Fault == nil {
			panic(`Group.GetBucket called without Parent`)
		} else {
			return teg.Fault
		}
	}
	return teg.Parent.(Bucketeer).GetBucket()
}

func (teg *Group) GetRepository() string {
	return teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
}

func (teg *Group) GetRepositoryName() string {
	return teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepositoryName()
}

func (teg *Group) GetEnvironment() string {
	return teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetEnvironment()
}

//
//
func (teg *Group) ComputeCheckInstances() {
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
func (teg *Group) ClearLoadInfo() {
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
func (teg *Group) export() proto.Group {
	bucket := teg.Parent.(Bucketeer).GetBucket()
	return proto.Group{
		Id:          teg.Id.String(),
		Name:        teg.Name,
		BucketId:    bucket.(Builder).GetID(),
		ObjectState: teg.State,
		TeamId:      teg.Team.String(),
	}
}

func (teg *Group) actionCreate() {
	teg.Action <- &Action{
		Action: "create",
		Type:   teg.Type,
		Group:  teg.export(),
	}
}

func (teg *Group) actionUpdate() {
	teg.Action <- &Action{
		Action: "update",
		Type:   teg.Type,
		Group:  teg.export(),
	}
}

func (teg *Group) actionDelete() {
	teg.Action <- &Action{
		Action: "delete",
		Type:   teg.Type,
		Group:  teg.export(),
	}
}

func (teg *Group) actionMemberNew(a Action) {
	a.Action = "member_new"
	a.Type = teg.Type
	a.Group = teg.export()

	teg.Action <- &a
}

func (teg *Group) actionMemberRemoved(a Action) {
	a.Action = "member_removed"
	a.Type = teg.Type
	a.Group = teg.export()

	teg.Action <- &a
}

//
func (teg *Group) actionPropertyNew(a Action) {
	a.Action = `property_new`
	teg.actionProperty(a)
}

func (teg *Group) actionPropertyUpdate(a Action) {
	a.Action = `property_update`
	teg.actionProperty(a)
}

func (teg *Group) actionPropertyDelete(a Action) {
	a.Action = `property_delete`
	teg.actionProperty(a)
}

func (teg *Group) actionProperty(a Action) {
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
func (teg *Group) actionCheckNew(a Action) {
	a.Check.RepositoryId = teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = teg.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	teg.actionDispatch("check_new", a)
}

func (teg *Group) actionCheckRemoved(a Action) {
	a.Check.RepositoryId = teg.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = teg.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	teg.actionDispatch(`check_removed`, a)
}

func (teg *Group) setupCheckAction(c Check) Action {
	return c.MakeAction()
}

func (teg *Group) actionCheckInstanceCreate(a Action) {
	teg.actionDispatch("check_instance_create", a)
}

func (teg *Group) actionCheckInstanceUpdate(a Action) {
	teg.actionDispatch("check_instance_update", a)
}

func (teg *Group) actionCheckInstanceDelete(a Action) {
	teg.actionDispatch("check_instance_delete", a)
}

func (teg *Group) actionDispatch(action string, a Action) {
	a.Action = action
	a.Type = teg.Type
	a.Group = teg.export()

	teg.Action <- &a
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
