package somatree

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/satori/go.uuid"
)

type SomaTreeElemRepository struct {
	Id              uuid.UUID
	Name            string
	Team            uuid.UUID
	Deleted         bool
	Active          bool
	Type            string
	State           string
	Parent          SomaTreeRepositoryReceiver `json:"-"`
	Fault           *SomaTreeElemFault         `json:"-"`
	PropertyOncall  map[string]SomaTreeProperty
	PropertyService map[string]SomaTreeProperty
	PropertySystem  map[string]SomaTreeProperty
	PropertyCustom  map[string]SomaTreeProperty
	Checks          map[string]Check
	Children        map[string]SomaTreeRepositoryAttacher // `json:"-"`
	Action          chan *Action                          `json:"-"`
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
func NewRepository(spec RepositorySpec) *SomaTreeElemRepository {
	if !specRepoCheck(spec) {
		fmt.Printf("%#v\n", spec) // XXX DEBUG
		panic(`No.`)
	}

	ter := new(SomaTreeElemRepository)
	ter.Id, _ = uuid.FromString(spec.Id)
	ter.Name = spec.Name
	ter.Team, _ = uuid.FromString(spec.Team)
	ter.Type = "repository"
	ter.State = "floating"
	ter.Parent = nil
	ter.Deleted = spec.Deleted
	ter.Active = spec.Active
	ter.Children = make(map[string]SomaTreeRepositoryAttacher)
	ter.PropertyOncall = make(map[string]SomaTreeProperty)
	ter.PropertyService = make(map[string]SomaTreeProperty)
	ter.PropertySystem = make(map[string]SomaTreeProperty)
	ter.PropertyCustom = make(map[string]SomaTreeProperty)
	ter.Checks = make(map[string]Check)

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

func (ter SomaTreeElemRepository) Clone() SomaTreeElemRepository {
	cl := SomaTreeElemRepository{
		Name:    ter.Name,
		Deleted: ter.Deleted,
		Active:  ter.Active,
		Type:    ter.Type,
		State:   ter.State,
	}
	cl.Id, _ = uuid.FromString(ter.Id.String())
	cl.Team, _ = uuid.FromString(ter.Id.String())
	f := make(map[string]SomaTreeRepositoryAttacher)
	for k, child := range ter.Children {
		f[k] = child.CloneRepository()
	}
	cl.Children = f

	pO := make(map[string]SomaTreeProperty)
	for k, prop := range ter.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]SomaTreeProperty)
	for k, prop := range ter.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]SomaTreeProperty)
	for k, prop := range ter.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]SomaTreeProperty)
	for k, prop := range ter.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]Check)
	for k, chk := range ter.Checks {
		cK[k] = chk.Clone()
	}
	cl.Checks = cK

	return cl
}

//
// Interface: SomaTreeBuilder
func (ter *SomaTreeElemRepository) GetID() string {
	return ter.Id.String()
}

func (ter *SomaTreeElemRepository) GetName() string {
	return ter.Name
}

func (ter *SomaTreeElemRepository) GetType() string {
	return ter.Type
}

func (ter *SomaTreeElemRepository) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeRepositoryReceiver:
		ter.setRepositoryParent(p.(SomaTreeRepositoryReceiver))
		ter.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemBucket.setParent`)
	}
}

func (ter *SomaTreeElemRepository) setAction(c chan *Action) {
	ter.Action = c
}

func (ter *SomaTreeElemRepository) setActionDeep(c chan *Action) {
	ter.setAction(c)
	ter.Fault.setActionDeep(c)
	for ch, _ := range ter.Children {
		ter.Children[ch].setActionDeep(c)
	}
}

func (ter *SomaTreeElemRepository) setError(c chan *Error) {
	if ter.Fault != nil {
		ter.Fault.setError(c)
	}
}

func (ter *SomaTreeElemRepository) getErrors() []error {
	if ter.Fault != nil {
		return ter.Fault.getErrors()
	}
	return []error{}
}

func (ter *SomaTreeElemRepository) setRepositoryParent(p SomaTreeRepositoryReceiver) {
	ter.Parent = p
}

func (ter *SomaTreeElemRepository) updateParentRecursive(p SomaTreeReceiver) {
	ter.setParent(p.(SomaTreeRepositoryReceiver))
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func(str SomaTreeReceiver) {
			defer wg.Done()
			ter.Children[c].updateParentRecursive(str)
		}(ter)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) clearParent() {
	ter.Parent = nil
	ter.State = "floating"
}

func (ter *SomaTreeElemRepository) setFault(f *SomaTreeElemFault) {
	ter.Fault = f
}

func (ter *SomaTreeElemRepository) updateFaultRecursive(f *SomaTreeElemFault) {
	ter.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func(ptr *SomaTreeElemFault) {
			defer wg.Done()
			ter.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

//
func (ter *SomaTreeElemRepository) actionCheckNew(a Action) {
	a.Action = "check_new"
	a.Type = ter.Type
	a.Group = ter.export()

	ter.Action <- &a
}

func (ter *SomaTreeElemRepository) setupCheckAction(c Check) Action {
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
	a.Check.RepositoryId = ter.Id.String()
	a.Check.BucketId = ""
	return a
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
