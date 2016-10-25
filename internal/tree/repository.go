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

	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

type Repository struct {
	Id              uuid.UUID
	Name            string
	Team            uuid.UUID
	Deleted         bool
	Active          bool
	Type            string
	State           string
	Parent          RepositoryReceiver `json:"-"`
	Fault           *Fault             `json:"-"`
	PropertyOncall  map[string]Property
	PropertyService map[string]Property
	PropertySystem  map[string]Property
	PropertyCustom  map[string]Property
	Checks          map[string]Check
	Children        map[string]RepositoryAttacher `json:"-"`
	Action          chan *Action                  `json:"-"`
	ordNumChildBck  int
	ordChildrenBck  map[int]string
}

type RepositorySpec struct {
	Id      string
	Name    string
	Team    string
	Deleted bool
	Active  bool
}

//
// NEW
func NewRepository(spec RepositorySpec) *Repository {
	if !specRepoCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	ter := new(Repository)
	ter.Id, _ = uuid.FromString(spec.Id)
	ter.Name = spec.Name
	ter.Team, _ = uuid.FromString(spec.Team)
	ter.Type = "repository"
	ter.State = "floating"
	ter.Parent = nil
	ter.Deleted = spec.Deleted
	ter.Active = spec.Active
	ter.Children = make(map[string]RepositoryAttacher)
	ter.PropertyOncall = make(map[string]Property)
	ter.PropertyService = make(map[string]Property)
	ter.PropertySystem = make(map[string]Property)
	ter.PropertyCustom = make(map[string]Property)
	ter.Checks = make(map[string]Check)
	ter.ordNumChildBck = 0
	ter.ordChildrenBck = make(map[int]string)

	// return new repository with attached fault handler
	newFault().Attach(
		AttachRequest{
			Root:       ter,
			ParentType: ter.Type,
			ParentName: ter.Name,
		},
	)
	return ter
}

func (ter Repository) Clone() Repository {
	cl := Repository{
		Name:           ter.Name,
		Deleted:        ter.Deleted,
		Active:         ter.Active,
		Type:           ter.Type,
		State:          ter.State,
		ordNumChildBck: ter.ordNumChildBck,
	}
	cl.Id, _ = uuid.FromString(ter.Id.String())
	cl.Team, _ = uuid.FromString(ter.Id.String())
	f := make(map[string]RepositoryAttacher)
	for k, child := range ter.Children {
		f[k] = child.CloneRepository()
	}
	cl.Children = f

	pO := make(map[string]Property)
	for k, prop := range ter.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]Property)
	for k, prop := range ter.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]Property)
	for k, prop := range ter.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]Property)
	for k, prop := range ter.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]Check)
	for k, chk := range ter.Checks {
		cK[k] = chk.Clone()
	}
	cl.Checks = cK

	chLB := make(map[int]string)
	for i, s := range ter.ordChildrenBck {
		chLB[i] = s
	}
	cl.ordChildrenBck = chLB

	return cl
}

//
// Interface: Builder
func (ter *Repository) GetID() string {
	return ter.Id.String()
}

func (ter *Repository) GetName() string {
	return ter.Name
}

func (ter *Repository) GetType() string {
	return ter.Type
}

func (ter *Repository) setParent(p Receiver) {
	switch p.(type) {
	case RepositoryReceiver:
		ter.setRepositoryParent(p.(RepositoryReceiver))
		ter.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`Bucket.setParent`)
	}
}

func (ter *Repository) setAction(c chan *Action) {
	ter.Action = c
}

func (ter *Repository) setActionDeep(c chan *Action) {
	ter.setAction(c)
	ter.Fault.setActionDeep(c)
	for ch, _ := range ter.Children {
		ter.Children[ch].setActionDeep(c)
	}
}

func (ter *Repository) setError(c chan *Error) {
	if ter.Fault != nil {
		ter.Fault.setError(c)
	}
}

func (ter *Repository) getErrors() []error {
	if ter.Fault != nil {
		return ter.Fault.getErrors()
	}
	return []error{}
}

func (ter *Repository) setRepositoryParent(p RepositoryReceiver) {
	ter.Parent = p
}

func (ter *Repository) updateParentRecursive(p Receiver) {
	ter.setParent(p.(RepositoryReceiver))
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func(str Receiver) {
			defer wg.Done()
			ter.Children[c].updateParentRecursive(str)
		}(ter)
	}
	wg.Wait()
}

func (ter *Repository) clearParent() {
	ter.Parent = nil
	ter.State = "floating"
}

func (ter *Repository) setFault(f *Fault) {
	ter.Fault = f
}

func (ter *Repository) updateFaultRecursive(f *Fault) {
	ter.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func(ptr *Fault) {
			defer wg.Done()
			ter.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

//
//
func (ter *Repository) ComputeCheckInstances() {
	log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s",
		ter.Name,
		`ComputeCheckInstances`,
		`repository`,
		ter.Id.String(),
	)
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			ter.Children[c].ComputeCheckInstances()
		}()
	}
	wg.Wait()
}

//
//
func (ter *Repository) ClearLoadInfo() {
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			ter.Children[c].ClearLoadInfo()
		}()
	}
	wg.Wait()
}

//
//
func (ter *Repository) export() proto.Repository {
	return proto.Repository{
		Id:        ter.Id.String(),
		Name:      ter.Name,
		TeamId:    ter.Team.String(),
		IsDeleted: ter.Deleted,
		IsActive:  ter.Active,
	}
}

func (ter *Repository) actionCreate() {
	ter.Action <- &Action{
		Action:     "create",
		Type:       ter.Type,
		Repository: ter.export(),
	}
}

func (ter *Repository) actionUpdate() {
	ter.Action <- &Action{
		Action:     "update",
		Type:       ter.Type,
		Repository: ter.export(),
	}
}

func (ter *Repository) actionDelete() {
	ter.Action <- &Action{
		Action:     "delete",
		Type:       ter.Type,
		Repository: ter.export(),
	}
}

//
func (ter *Repository) actionPropertyNew(a Action) {
	a.Action = "property_new"
	ter.actionProperty(a)
}

func (ter *Repository) actionPropertyUpdate(a Action) {
	a.Action = "property_update"
	ter.actionProperty(a)
}

func (ter *Repository) actionPropertyDelete(a Action) {
	a.Action = "property_delete"
	ter.actionProperty(a)
}

func (ter *Repository) actionProperty(a Action) {
	a.Type = ter.Type
	a.Repository = ter.export()
	a.Property.RepositoryId = ter.Id.String()
	a.Property.BucketId = ""

	switch a.Property.Type {
	case "custom":
		a.Property.Custom.RepositoryId = a.Property.RepositoryId
	case "service":
		a.Property.Service.TeamId = ter.Team.String()
	}

	ter.Action <- &a
}

//
func (ter *Repository) actionCheckNew(a Action) {
	a.Action = "check_new"
	a.Type = ter.Type
	a.Repository = ter.export()
	a.Check.RepositoryId = ter.Id.String()
	a.Check.BucketId = ""

	ter.Action <- &a
}

func (ter *Repository) actionCheckRemoved(a Action) {
	a.Action = `check_removed`
	a.Type = ter.Type
	a.Repository = ter.export()
	a.Check.RepositoryId = ter.Id.String()
	a.Check.BucketId = ""

	ter.Action <- &a
}

func (ter *Repository) setupCheckAction(c Check) Action {
	return c.MakeAction()
	/*
			a := Action{
				Check: somaproto.TreeCheck{
					CheckId:       c.GetCheckId(),
					SourceCheckId: c.GetSourceCheckId(),
					CheckConfigId: c.GetCheckConfigId(),
					SourceType:    c.GetSourceType(),
					IsInherited:   c.GetIsInherited(),
					InheritedFrom: c.GetInheritedFrom(),
					Inheritance:   c.GetInheritance(),
					ChildrenOnly:  c.GetChildrenOnly(),
					CapabilityId:  c.GetCapabilityId(),
				},
			}
		return a
	*/
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
