package somatree

import (
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (teg *SomaTreeElemGroup) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyCustom).Id, uuid.Nil) {
			p.(*PropertyCustom).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyCustom).InheritedFrom = teg.Id
		p.(*PropertyCustom).Inherited = false
		p.(*PropertyCustom).SourceId = p.(*PropertyCustom).Id
		// send a scrubbed copy downward
		f := new(PropertyCustom)
		*f = *p.(*PropertyCustom)
		f.Inherited = true
		f.Id = uuid.Nil
		teg.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyCustom).Instances = nil
		teg.setCustomProperty(p)
	case "service":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyService).Id, uuid.Nil) {
			p.(*PropertyService).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyService).InheritedFrom = teg.Id
		p.(*PropertyService).Inherited = false
		p.(*PropertyService).SourceId = p.(*PropertyService).Id
		// send a scrubbed copy downward
		f := new(PropertyService)
		*f = *p.(*PropertyService)
		f.Inherited = true
		f.Id = uuid.Nil
		teg.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyService).Instances = nil
		teg.setServiceProperty(p)
	case "system":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertySystem).Id, uuid.Nil) {
			p.(*PropertySystem).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertySystem).InheritedFrom = teg.Id
		p.(*PropertySystem).Inherited = false
		p.(*PropertySystem).SourceId = p.(*PropertySystem).Id
		// send a scrubbed copy downward
		f := new(PropertySystem)
		*f = *p.(*PropertySystem)
		f.Inherited = true
		f.Id = uuid.Nil
		teg.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertySystem).Instances = nil
		teg.setSystemProperty(p)
	case "oncall":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyOncall).Id, uuid.Nil) {
			p.(*PropertyOncall).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyOncall).InheritedFrom = teg.Id
		p.(*PropertyOncall).Inherited = false
		p.(*PropertyOncall).SourceId = p.(*PropertyOncall).Id
		// send a scrubbed copy downward
		f := new(PropertyOncall)
		*f = *p.(*PropertyOncall)
		f.Inherited = true
		f.Id = uuid.Nil
		teg.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyOncall).Instances = nil
		teg.setOncallProperty(p)
	}
	teg.actionPropertyNew(teg.setupPropertyAction(p))
}

func (teg *SomaTreeElemGroup) inheritProperty(
	p SomaTreeProperty) {

	switch p.GetType() {
	case "custom":
		f := new(PropertyCustom)
		*f = *p.(*PropertyCustom)
		p.(*PropertyCustom).Id = p.GetInstanceId(teg.Type, teg.Id)
		if uuid.Equal(p.(*PropertyCustom).Id, uuid.Nil) {
			p.(*PropertyCustom).Id = uuid.NewV4()
		}
		p.(*PropertyCustom).Instances = nil
		teg.setCustomProperty(p)
		teg.inheritPropertyDeep(f)
	case "service":
		f := new(PropertyService)
		*f = *p.(*PropertyService)
		p.(*PropertyService).Id = p.GetInstanceId(teg.Type, teg.Id)
		if uuid.Equal(p.(*PropertyService).Id, uuid.Nil) {
			p.(*PropertyService).Id = uuid.NewV4()
		}
		p.(*PropertyService).Instances = nil
		teg.setServiceProperty(p)
		teg.inheritPropertyDeep(f)
	case "system":
		f := new(PropertySystem)
		*f = *p.(*PropertySystem)
		p.(*PropertySystem).Id = p.GetInstanceId(teg.Type, teg.Id)
		if uuid.Equal(p.(*PropertySystem).Id, uuid.Nil) {
			p.(*PropertySystem).Id = uuid.NewV4()
		}
		p.(*PropertySystem).Instances = nil
		teg.setSystemProperty(p)
		teg.inheritPropertyDeep(f)
	case "oncall":
		f := new(PropertyOncall)
		*f = *p.(*PropertyOncall)
		p.(*PropertyOncall).Id = p.GetInstanceId(teg.Type, teg.Id)
		if uuid.Equal(p.(*PropertyOncall).Id, uuid.Nil) {
			p.(*PropertyOncall).Id = uuid.NewV4()
		}
		p.(*PropertyOncall).Instances = nil
		teg.setOncallProperty(p)
		teg.inheritPropertyDeep(f)
	}
	teg.actionPropertyNew(teg.setupPropertyAction(p))
}

func (teg *SomaTreeElemGroup) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func(stp SomaTreeProperty) {
			defer wg.Done()
			teg.Children[c].inheritProperty(stp)
		}(p)
	}
	wg.Wait()
}

func (teg *SomaTreeElemGroup) setCustomProperty(
	p SomaTreeProperty) {
	teg.PropertyCustom[p.GetID()] = p
}

func (teg *SomaTreeElemGroup) setServiceProperty(
	p SomaTreeProperty) {
	teg.PropertyService[p.GetID()] = p
}

func (teg *SomaTreeElemGroup) setSystemProperty(
	p SomaTreeProperty) {
	teg.PropertySystem[p.GetID()] = p
}

func (teg *SomaTreeElemGroup) setOncallProperty(
	p SomaTreeProperty) {
	teg.PropertyOncall[p.GetID()] = p
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (teg *SomaTreeElemGroup) syncProperty(
	childId string) {
customloop:
	for prop, _ := range teg.PropertyCustom {
		if !teg.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := new(PropertyCustom)
		*f = *teg.PropertyCustom[prop].(*PropertyCustom)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range teg.PropertyOncall {
		if !teg.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(PropertyOncall)
		*f = *teg.PropertyOncall[prop].(*PropertyOncall)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range teg.PropertyService {
		if !teg.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(PropertyService)
		*f = *teg.PropertyService[prop].(*PropertyService)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range teg.PropertySystem {
		if !teg.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(PropertySystem)
		*f = *teg.PropertySystem[prop].(*PropertySystem)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (teg *SomaTreeElemGroup) checkProperty(
	propType string, propId string) bool {
	switch propType {
	case "custom":
		if _, ok := teg.PropertyCustom[propId]; ok {
			return true
		}
	case "service":
		if _, ok := teg.PropertyService[propId]; ok {
			return true
		}
	case "system":
		if _, ok := teg.PropertySystem[propId]; ok {
			return true
		}
	case "oncall":
		if _, ok := teg.PropertyOncall[propId]; ok {
			return true
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
