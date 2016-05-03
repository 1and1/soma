package somatree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (tec *SomaTreeElemCluster) SetProperty(p SomaTreeProperty) {
	p.SetId(p.GetInstanceId(tec.Type, tec.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
	log.Printf("SetProperty(Cluster) created source instance: %s", p.GetID())
	// this property is the source instance
	p.SetInheritedFrom(tec.Id)
	p.SetInherited(false)
	p.SetSourceType(tec.Type)
	if i, e := uuid.FromString(p.GetID()); e == nil {
		p.SetSourceId(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetId(uuid.UUID{})
	tec.inheritPropertyDeep(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	switch p.GetType() {
	case "custom":
		tec.setCustomProperty(p)
	case "service":
		tec.setServiceProperty(p)
	case "system":
		tec.setSystemProperty(p)
	case "oncall":
		tec.setOncallProperty(p)
	}
	tec.actionPropertyNew(p.MakeAction())
}

func (tec *SomaTreeElemCluster) inheritProperty(p SomaTreeProperty) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(tec.Type, tec.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Cluster) Generated: %s", f.GetID())
	}
	f.clearInstances()

	switch f.GetType() {
	case "custom":
		tec.setCustomProperty(f)
	case "service":
		tec.setServiceProperty(f)
	case "system":
		tec.setSystemProperty(f)
	case "oncall":
		tec.setOncallProperty(f)
	}
	p.SetId(uuid.UUID{})
	tec.inheritPropertyDeep(p)
	tec.actionPropertyNew(f.MakeAction())
}

func (tec *SomaTreeElemCluster) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
	log.Printf("InheritDeep Sending down: %s", p.GetID())
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(stp SomaTreeProperty) {
			defer wg.Done()
			tec.Children[c].inheritProperty(stp)
		}(p)
	}
	wg.Wait()
}

func (tec *SomaTreeElemCluster) setCustomProperty(
	p SomaTreeProperty) {
	tec.PropertyCustom[p.GetID()] = p
}

func (tec *SomaTreeElemCluster) setServiceProperty(
	p SomaTreeProperty) {
	tec.PropertyService[p.GetID()] = p
}

func (tec *SomaTreeElemCluster) setSystemProperty(
	p SomaTreeProperty) {
	tec.PropertySystem[p.GetID()] = p
}

func (tec *SomaTreeElemCluster) setOncallProperty(
	p SomaTreeProperty) {
	tec.PropertyOncall[p.GetID()] = p
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (tec *SomaTreeElemCluster) syncProperty(
	childId string) {
customloop:
	for prop, _ := range tec.PropertyCustom {
		if !tec.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := new(PropertyCustom)
		*f = *tec.PropertyCustom[prop].(*PropertyCustom)
		f.Inherited = true
		f.Id = uuid.Nil
		tec.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range tec.PropertyOncall {
		if !tec.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(PropertyOncall)
		*f = *tec.PropertyOncall[prop].(*PropertyOncall)
		f.Inherited = true
		f.Id = uuid.Nil
		tec.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range tec.PropertyService {
		if !tec.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(PropertyService)
		*f = *tec.PropertyService[prop].(*PropertyService)
		f.Inherited = true
		f.Id = uuid.Nil
		tec.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range tec.PropertySystem {
		if !tec.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(PropertySystem)
		*f = *tec.PropertySystem[prop].(*PropertySystem)
		f.Inherited = true
		f.Id = uuid.Nil
		tec.Children[childId].inheritProperty(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (tec *SomaTreeElemCluster) checkProperty(
	propType string, propId string) bool {
	switch propType {
	case "custom":
		if _, ok := tec.PropertyCustom[propId]; ok {
			return true
		}
	case "service":
		if _, ok := tec.PropertyService[propId]; ok {
			return true
		}
	case "system":
		if _, ok := tec.PropertySystem[propId]; ok {
			return true
		}
	case "oncall":
		if _, ok := tec.PropertyOncall[propId]; ok {
			return true
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
