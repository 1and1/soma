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
	PropertyOncall  map[string]SomaTreeProperty
	PropertyService map[string]SomaTreeProperty
	PropertySystem  map[string]SomaTreeProperty
	PropertyCustom  map[string]SomaTreeProperty
	Checks          map[string]Check
	CheckInstances  map[string][]string
	Instances       map[string]CheckInstance
	Children        map[string]SomaTreeClusterAttacher //`json:"-"`
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
	tec.PropertyOncall = make(map[string]SomaTreeProperty)
	tec.PropertyService = make(map[string]SomaTreeProperty)
	tec.PropertySystem = make(map[string]SomaTreeProperty)
	tec.PropertyCustom = make(map[string]SomaTreeProperty)
	tec.Checks = make(map[string]Check)

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

	pO := make(map[string]SomaTreeProperty)
	for k, prop := range tec.PropertyOncall {
		pO[k] = prop.Clone()
	}
	cl.PropertyOncall = pO

	pSv := make(map[string]SomaTreeProperty)
	for k, prop := range tec.PropertyService {
		pSv[k] = prop.Clone()
	}
	cl.PropertyService = pSv

	pSy := make(map[string]SomaTreeProperty)
	for k, prop := range tec.PropertySystem {
		pSy[k] = prop.Clone()
	}
	cl.PropertySystem = pSy

	pC := make(map[string]SomaTreeProperty)
	for k, prop := range tec.PropertyCustom {
		pC[k] = prop.Clone()
	}
	cl.PropertyCustom = pC

	cK := make(map[string]Check)
	for k, chk := range tec.Checks {
		cK[k] = chk.Clone()
	}
	cl.Checks = cK

	cki := make(map[string]CheckInstance)
	for k, chki := range tec.Instances {
		cki[k] = chki.Clone()
	}
	cl.Instances = cki

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

func (tec *SomaTreeElemCluster) GetEnvironment() string {
	return tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetEnvironment()
}

//
//
func (tec *SomaTreeElemCluster) export() somaproto.ProtoCluster {
	bucket := tec.Parent.(Bucketeer).GetBucket()
	return somaproto.ProtoCluster{
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
	a.Action = "property_new"
	a.Type = tec.Type
	a.Cluster = tec.export()

	tec.Action <- &a
}

//
func (tec *SomaTreeElemCluster) setupPropertyAction(p SomaTreeProperty) Action {
	a := Action{
		Property: somaproto.TreeProperty{
			InstanceId:       p.GetID(),
			SourceInstanceId: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			PropertyType:     p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			RepositoryId:     tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository(),
			BucketId:         tec.Parent.(Bucketeer).GetBucket().(Builder).GetID(),
		},
	}
	switch a.Property.PropertyType {
	case "custom":
		a.Property.Custom = &somaproto.TreePropertyCustom{
			CustomId:     p.(*PropertyCustom).CustomId.String(),
			RepositoryId: a.Property.RepositoryId,
			Name:         p.(*PropertyCustom).Key,
			Value:        p.(*PropertyCustom).Value,
		}
	case "system":
		a.Property.System = &somaproto.TreePropertySystem{
			Name:  p.(*PropertySystem).Key,
			Value: p.(*PropertySystem).Value,
		}
	case "service":
		a.Property.Service = &somaproto.TreePropertyService{
			Name:   p.(*PropertyService).Service,
			TeamId: tec.Team.String(),
		}
		a.Property.Service.Attributes = make([]somaproto.TreeServiceAttribute, 0)
		for _, attr := range p.(*PropertyService).Attributes {
			ta := somaproto.TreeServiceAttribute{
				Attribute: attr.Attribute,
				Value:     attr.Value,
			}
			a.Property.Service.Attributes = append(a.Property.Service.Attributes, ta)
		}
	case "oncall":
		a.Property.Oncall = &somaproto.TreePropertyOncall{
			OncallId: p.(*PropertyOncall).OncallId.String(),
			Name:     p.(*PropertyOncall).Name,
			Number:   p.(*PropertyOncall).Number,
		}
	}
	return a
}

//
func (tec *SomaTreeElemCluster) actionCheckNew(a Action) {
	a.Action = "check_new"
	a.Type = tec.Type
	a.Cluster = tec.export()

	tec.Action <- &a
}

func (tec *SomaTreeElemCluster) setupCheckAction(c Check) Action {
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
	a.Check.RepositoryId = tec.Parent.(Bucketeer).GetBucket().(Bucketeer).GetRepository()
	a.Check.BucketId = tec.Parent.(Bucketeer).GetBucket().(Builder).GetID()
	return a
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
