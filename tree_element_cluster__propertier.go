package somatree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (tec *SomaTreeElemCluster) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyCustom).Id, uuid.Nil) {
			p.(*PropertyCustom).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyCustom).InheritedFrom = tec.Id
		p.(*PropertyCustom).Inherited = false
		p.(*PropertyCustom).SourceId = p.(*PropertyCustom).Id
		p.(*PropertyCustom).SourceType = tec.Type
		// send a scrubbed copy downward
		f := new(PropertyCustom)
		*f = *p.(*PropertyCustom)
		f.Inherited = true
		f.Id = uuid.Nil
		tec.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyCustom).Instances = nil
		tec.setCustomProperty(p)
	case "service":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyService).Id, uuid.Nil) {
			p.(*PropertyService).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyService).InheritedFrom = tec.Id
		p.(*PropertyService).Inherited = false
		p.(*PropertyService).SourceId = p.(*PropertyService).Id
		p.(*PropertyService).SourceType = tec.Type
		// send a scrubbed copy downward
		f := new(PropertyService)
		*f = *p.(*PropertyService)
		f.Inherited = true
		f.Id = uuid.Nil
		tec.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyService).Instances = nil
		tec.setServiceProperty(p)
	case "system":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertySystem).Id, uuid.Nil) {
			p.(*PropertySystem).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertySystem).InheritedFrom = tec.Id
		p.(*PropertySystem).Inherited = false
		p.(*PropertySystem).SourceId = p.(*PropertySystem).Id
		p.(*PropertySystem).SourceType = tec.Type
		// send a scrubbed copy downward
		f := new(PropertySystem)
		*f = *p.(*PropertySystem)
		f.Inherited = true
		f.Id = uuid.UUID{}
		tec.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertySystem).Instances = nil
		tec.setSystemProperty(p)
	case "oncall":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyOncall).Id, uuid.Nil) {
			p.(*PropertyOncall).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyOncall).InheritedFrom = tec.Id
		p.(*PropertyOncall).Inherited = false
		p.(*PropertyOncall).SourceId = p.(*PropertyOncall).Id
		p.(*PropertyOncall).SourceType = tec.Type
		// send a scrubbed copy downward
		f := new(PropertyOncall)
		*f = *p.(*PropertyOncall)
		f.Inherited = true
		f.Id = uuid.Nil
		tec.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyOncall).Instances = nil
		tec.setOncallProperty(p)
	}
	tec.actionPropertyNew(tec.setupPropertyAction(p))
}

func (tec *SomaTreeElemCluster) inheritProperty(
	p SomaTreeProperty) {

	f := new(PropertySystem)
	switch p.GetType() {
	case "custom":
		f := new(PropertyCustom)
		*f = *p.(*PropertyCustom)
		p.(*PropertyCustom).Id = p.GetInstanceId(tec.Type, tec.Id)
		if uuid.Equal(p.(*PropertyCustom).Id, uuid.Nil) {
			p.(*PropertyCustom).Id = uuid.NewV4()
		}
		p.(*PropertyCustom).Instances = nil
		tec.setCustomProperty(p)
		tec.inheritPropertyDeep(f)
	case "service":
		f := new(PropertyService)
		*f = *p.(*PropertyService)
		p.(*PropertyService).Id = p.GetInstanceId(tec.Type, tec.Id)
		if uuid.Equal(p.(*PropertyService).Id, uuid.Nil) {
			p.(*PropertyService).Id = uuid.NewV4()
		}
		p.(*PropertyService).Instances = nil
		tec.setServiceProperty(p)
		tec.inheritPropertyDeep(f)
	case "system":
		/*
			f := new(PropertySystem)
			*f = *p.(*PropertySystem)
			p.(*PropertySystem).Id = p.GetInstanceId(tec.Type, tec.Id)
			if uuid.Equal(p.(*PropertySystem).Id, uuid.Nil) {
				p.(*PropertySystem).Id = uuid.NewV4()
			}
			p.(*PropertySystem).Instances = nil
			tec.setSystemProperty(p)
			tec.inheritPropertyDeep(f)
		*/
		*f = *p.(*PropertySystem)
		f.Id = f.GetInstanceId(tec.Type, tec.Id)
		if uuid.Equal(f.Id, uuid.Nil) {
			f.Id = uuid.NewV4()
			log.Printf("Inherit (Cluster) Generated: %s", f.Id.String())
		}
		f.Instances = nil
		tec.setSystemProperty(f)
		p.(*PropertySystem).Id = uuid.UUID{}
		log.Printf("Inherit Sending down: %s", p.(*PropertySystem).Id.String())
		tec.inheritPropertyDeep(p)
	case "oncall":
		f := new(PropertyOncall)
		*f = *p.(*PropertyOncall)
		p.(*PropertyOncall).Id = p.GetInstanceId(tec.Type, tec.Id)
		if uuid.Equal(p.(*PropertyOncall).Id, uuid.Nil) {
			p.(*PropertyOncall).Id = uuid.NewV4()
		}
		p.(*PropertyOncall).Instances = nil
		tec.setOncallProperty(p)
		tec.inheritPropertyDeep(f)
	}
	tec.actionPropertyNew(tec.setupPropertyAction(f))
}

func (tec *SomaTreeElemCluster) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
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
