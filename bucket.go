package tree

import (
	"fmt"
	"log"
	"sync"


	"github.com/satori/go.uuid"
)

type Bucket struct {
	Id              uuid.UUID
	Name            string
	Environment     string
	Type            string
	State           string
	Frozen          bool
	Deleted         bool
	Repository      uuid.UUID
	Team            uuid.UUID
	Parent          BucketReceiver `json:"-"`
	Fault           *Fault         `json:"-"`
	PropertyOncall  map[string]Property
	PropertyService map[string]Property
	PropertySystem  map[string]Property
	PropertyCustom  map[string]Property
	Checks          map[string]Check
	Children        map[string]BucketAttacher //`json:"-"`
	Action          chan *Action              `json:"-"`
}

type BucketSpec struct {
	Id          string
	Name        string
	Environment string
	Team        string
	Repository  string
	Deleted     bool
	Frozen      bool
}

//
// NEW
func NewBucket(spec BucketSpec) *Bucket {
	if !specBucketCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	teb := new(Bucket)
	teb.Id, _ = uuid.FromString(spec.Id)
	teb.Name = spec.Name
	teb.Team, _ = uuid.FromString(spec.Team)
	teb.Environment = spec.Environment
	teb.Frozen = spec.Frozen
	teb.Deleted = spec.Deleted
	teb.Repository, _ = uuid.FromString(spec.Repository)
	teb.Type = "bucket"
	teb.State = "floating"
	teb.Parent = nil
	teb.Children = make(map[string]BucketAttacher)
	teb.PropertyOncall = make(map[string]Property)
	teb.PropertyService = make(map[string]Property)
	teb.PropertySystem = make(map[string]Property)
	teb.PropertyCustom = make(map[string]Property)
	teb.Checks = make(map[string]Check)

	return teb
}

func (teb Bucket) CloneRepository() RepositoryAttacher {
	cl := Bucket{
		Name:        teb.Name,
		Environment: teb.Environment,
		Type:        teb.Type,
		State:       teb.State,
		Frozen:      teb.Frozen,
		Deleted:     teb.Deleted,
	}
	cl.Id, _ = uuid.FromString(teb.Id.String())
	cl.Team, _ = uuid.FromString(teb.Team.String())
	cl.Repository, _ = uuid.FromString(teb.Repository.String())

	f := make(map[string]BucketAttacher)
	for k, child := range teb.Children {
		f[k] = child.CloneBucket()
	}
	cl.Children = f

	pO := make(map[string]Property)
	for k, prop := range teb.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]Property)
	for k, prop := range teb.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]Property)
	for k, prop := range teb.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]Property)
	for k, prop := range teb.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]Check)
	for k, chk := range teb.Checks {
		cK[k] = chk.clone()
	}
	cl.Checks = cK

	return &cl
}

//
// Interface: Builder
func (teb *Bucket) GetID() string {
	return teb.Id.String()
}

func (teb *Bucket) GetName() string {
	return teb.Name
}

func (teb *Bucket) GetType() string {
	return teb.Type
}

func (teb *Bucket) setAction(c chan *Action) {
	teb.Action = c
}

func (teb *Bucket) setActionDeep(c chan *Action) {
	teb.setAction(c)
	for ch, _ := range teb.Children {
		teb.Children[ch].setActionDeep(c)
	}
}

//
// Interface: Bucketeer
func (teb *Bucket) GetBucket() Receiver {
	return teb
}

func (teb *Bucket) GetEnvironment() string {
	return teb.Environment
}

func (teb *Bucket) GetRepository() string {
	return teb.Repository.String()
}

func (teb *Bucket) GetRepositoryName() string {
	return teb.Parent.(*Repository).GetName()
}

//
//
func (teb *Bucket) ComputeCheckInstances() {
	log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s",
		teb.GetRepositoryName(),
		`ComputeCheckInstances`,
		`bucket`,
		teb.Id.String(),
	)
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			teb.Children[c].ComputeCheckInstances()
		}()
	}
	wg.Wait()
}

//
//
func (teb *Bucket) ClearLoadInfo() {
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			teb.Children[c].ClearLoadInfo()
		}()
	}
	wg.Wait()
}

//
//
func (teb *Bucket) export() proto.Bucket {
	return proto.Bucket{
		Id:           teb.Id.String(),
		Name:         teb.Name,
		RepositoryId: teb.Repository.String(),
		TeamId:       teb.Team.String(),
		Environment:  teb.Environment,
		IsDeleted:    teb.Deleted,
		IsFrozen:     teb.Frozen,
	}
}

func (teb *Bucket) actionCreate() {
	teb.Action <- &Action{
		Action: "create",
		Type:   teb.Type,
		Bucket: teb.export(),
	}
}

func (teb *Bucket) actionUpdate() {
	teb.Action <- &Action{
		Action: "update",
		Type:   teb.Type,
		Bucket: teb.export(),
	}
}

func (teb *Bucket) actionDelete() {
	teb.Action <- &Action{
		Action: "delete",
		Type:   teb.Type,
		Bucket: teb.export(),
	}
}

func (teb *Bucket) actionAssignNode(a Action) {
	a.Action = "node_assignment"
	a.Type = teb.Type
	a.Bucket = teb.export()

	teb.Action <- &a
}

//
func (teb *Bucket) actionPropertyNew(a Action) {
	a.Action = "property_new"
	teb.actionProperty(a)
}

func (teb *Bucket) actionPropertyUpdate(a Action) {
	a.Action = `property_update`
	teb.actionProperty(a)
}

func (teb *Bucket) actionPropertyDelete(a Action) {
	a.Action = `property_delete`
	teb.actionProperty(a)
}

func (teb *Bucket) actionProperty(a Action) {
	a.Type = teb.Type
	a.Bucket = teb.export()

	a.Property.RepositoryId = teb.Repository.String()
	a.Property.BucketId = teb.Id.String()
	switch a.Property.Type {
	case `custom`:
		a.Property.Custom.RepositoryId = a.Property.RepositoryId
	case `service`:
		a.Property.Service.TeamId = teb.Team.String()
	}

	teb.Action <- &a
}

//
func (teb *Bucket) actionCheckNew(a Action) {
	a.Action = "check_new"
	a.Type = teb.Type
	a.Bucket = teb.export()
	a.Check.RepositoryId = teb.Repository.String()
	a.Check.BucketId = teb.Id.String()

	teb.Action <- &a
}

func (teb *Bucket) actionCheckRemoved(a Action) {
	a.Action = `check_removed`
	a.Type = teb.Type
	a.Bucket = teb.export()
	a.Check.RepositoryId = teb.Repository.String()
	a.Check.BucketId = teb.Id.String()

	teb.Action <- &a
}

func (teb *Bucket) setupCheckAction(c Check) Action {
	return c.MakeAction()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
