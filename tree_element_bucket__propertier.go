package somatree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (teb *SomaTreeElemBucket) SetProperty(p SomaTreeProperty) {
	p.SetId(p.GetInstanceId(teb.Type, teb.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
	// this property is the source instance
	p.SetInheritedFrom(teb.Id)
	p.SetInherited(false)
	p.SetSourceType(teb.Type)
	if i, e := uuid.FromString(p.GetID()); e != nil {
		p.SetSourceId(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetId(uuid.UUID{})
	teb.inheritPropertyDeep(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	switch p.GetType() {
	case "custom":
		teb.setCustomProperty(p)
	case "service":
		teb.setServiceProperty(p)
	case "system":
		teb.setSystemProperty(p)
	case "oncall":
		teb.setOncallProperty(p)
	}
	teb.actionPropertyNew(p.MakeAction())
}

func (teb *SomaTreeElemBucket) inheritProperty(p SomaTreeProperty) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(teb.Type, teb.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Bucket) Generated: %s", f.GetID())
	}
	f.clearInstances()

	switch f.GetType() {
	case "custom":
		teb.setCustomProperty(f)
	case "service":
		teb.setServiceProperty(f)
	case "system":
		teb.setSystemProperty(f)
	case "oncall":
		teb.setOncallProperty(f)
	}
	p.SetId(uuid.UUID{})
	teb.inheritPropertyDeep(p)
	teb.actionPropertyNew(f.MakeAction())
}

func (teb *SomaTreeElemBucket) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
	log.Printf("InheritDeep Sending down: %s", p.GetID())
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func(stp SomaTreeProperty) {
			defer wg.Done()
			teb.Children[c].inheritProperty(stp)
		}(p)
	}
	wg.Wait()
}

func (teb *SomaTreeElemBucket) setCustomProperty(
	p SomaTreeProperty) {
	teb.PropertyCustom[p.GetID()] = p
}

func (teb *SomaTreeElemBucket) setServiceProperty(
	p SomaTreeProperty) {
	teb.PropertyService[p.GetID()] = p
}

func (teb *SomaTreeElemBucket) setSystemProperty(
	p SomaTreeProperty) {
	teb.PropertySystem[p.GetID()] = p
}

func (teb *SomaTreeElemBucket) setOncallProperty(
	p SomaTreeProperty) {
	teb.PropertyOncall[p.GetID()] = p
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (teb *SomaTreeElemBucket) syncProperty(
	childId string) {
customloop:
	for prop, _ := range teb.PropertyCustom {
		if !teb.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := new(PropertyCustom)
		*f = *teb.PropertyCustom[prop].(*PropertyCustom)
		f.Inherited = true
		teb.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range teb.PropertyOncall {
		if !teb.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(PropertyOncall)
		*f = *teb.PropertyOncall[prop].(*PropertyOncall)
		f.Inherited = true
		teb.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range teb.PropertyService {
		if !teb.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(PropertyService)
		*f = *teb.PropertyService[prop].(*PropertyService)
		f.Inherited = true
		teb.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range teb.PropertySystem {
		if !teb.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(PropertySystem)
		*f = *teb.PropertySystem[prop].(*PropertySystem)
		f.Inherited = true
		teb.Children[childId].inheritProperty(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (teb *SomaTreeElemBucket) checkProperty(
	propType string, propId string) bool {
	switch propType {
	case "custom":
		if _, ok := teb.PropertyCustom[propId]; ok {
			return true
		}
	case "service":
		if _, ok := teb.PropertyService[propId]; ok {
			return true
		}
	case "system":
		if _, ok := teb.PropertySystem[propId]; ok {
			return true
		}
	case "oncall":
		if _, ok := teb.PropertyOncall[propId]; ok {
			return true
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
