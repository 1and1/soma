package somatree

import (
	"log"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (ten *SomaTreeElemNode) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyCustom).Id, uuid.Nil) {
			p.(*PropertyCustom).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyCustom).InheritedFrom = ten.Id
		p.(*PropertyCustom).Inherited = false
		p.(*PropertyCustom).SourceId = p.(*PropertyCustom).Id
		p.(*PropertyCustom).SourceType = ten.Type
		// send a scrubbed copy downward
		f := new(PropertyCustom)
		*f = *p.(*PropertyCustom)
		f.Inherited = true
		f.Id = uuid.Nil
		ten.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyCustom).Instances = nil
		ten.setCustomProperty(p)
	case "service":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyService).Id, uuid.Nil) {
			p.(*PropertyService).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyService).InheritedFrom = ten.Id
		p.(*PropertyService).Inherited = false
		p.(*PropertyService).SourceId = p.(*PropertyService).Id
		p.(*PropertyService).SourceType = ten.Type
		// send a scrubbed copy downward
		f := new(PropertyService)
		*f = *p.(*PropertyService)
		f.Inherited = true
		f.Id = uuid.UUID{}
		ten.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyService).Instances = nil
		ten.setServiceProperty(p)
	case "system":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertySystem).Id, uuid.Nil) {
			p.(*PropertySystem).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertySystem).InheritedFrom = ten.Id
		p.(*PropertySystem).Inherited = false
		p.(*PropertySystem).SourceId = p.(*PropertySystem).Id
		p.(*PropertySystem).SourceType = ten.Type
		// send a scrubbed copy downward
		f := new(PropertySystem)
		*f = *p.(*PropertySystem)
		f.Inherited = true
		f.Id = uuid.Nil
		ten.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertySystem).Instances = nil
		ten.setSystemProperty(p)
	case "oncall":
		// generate uuid if none is set
		if uuid.Equal(p.(*PropertyOncall).Id, uuid.Nil) {
			p.(*PropertyOncall).Id = uuid.NewV4()
		}
		// this property is the source instance
		p.(*PropertyOncall).InheritedFrom = ten.Id
		p.(*PropertyOncall).Inherited = false
		p.(*PropertyOncall).SourceId = p.(*PropertyOncall).Id
		p.(*PropertyOncall).SourceType = ten.Type
		// send a scrubbed copy downward
		f := new(PropertyOncall)
		*f = *p.(*PropertyOncall)
		f.Inherited = true
		f.Id = uuid.Nil
		ten.inheritPropertyDeep(f)
		// scrub instance startup information prior to setting
		p.(*PropertyOncall).Instances = nil
		ten.setOncallProperty(p)
	}
	ten.actionPropertyNew(ten.setupPropertyAction(p))
}

func (ten *SomaTreeElemNode) inheritProperty(
	p SomaTreeProperty) {
	f := new(PropertySystem)
	switch p.GetType() {
	case "custom":
		p.(*PropertyCustom).Id = p.GetInstanceId(ten.Type, ten.Id)
		if uuid.Equal(p.(*PropertyCustom).Id, uuid.Nil) {
			p.(*PropertyCustom).Id = uuid.NewV4()
		}
		ten.setCustomProperty(p)
	case "service":
		p.(*PropertyService).Id = p.GetInstanceId(ten.Type, ten.Id)
		if uuid.Equal(p.(*PropertyService).Id, uuid.Nil) {
			p.(*PropertyService).Id = uuid.NewV4()
		}
		ten.setServiceProperty(p)
	case "system":
		/*
			p.(*PropertySystem).Id = p.GetInstanceId(ten.Type, ten.Id)
			if uuid.Equal(p.(*PropertySystem).Id, uuid.Nil) {
				p.(*PropertySystem).Id = uuid.NewV4()
			}
			ten.setSystemProperty(p)
		*/
		*f = *p.(*PropertySystem)
		f.Id = f.GetInstanceId(ten.Type, ten.Id)
		if uuid.Equal(f.Id, uuid.Nil) {
			f.Id = uuid.NewV4()
			log.Printf("Inherit (Node) Generated: %s", f.Id.String())
		}
		f.Instances = nil
		ten.setSystemProperty(f)
	case "oncall":
		p.(*PropertyOncall).Id = p.GetInstanceId(ten.Type, ten.Id)
		if uuid.Equal(p.(*PropertyOncall).Id, uuid.Nil) {
			p.(*PropertyOncall).Id = uuid.NewV4()
		}
		ten.setOncallProperty(p)
	}
	ten.actionPropertyNew(ten.setupPropertyAction(f))
	// no inheritPropertyDeep(), nodes have no children
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) inheritPropertyDeep(
	p SomaTreeProperty) {
}

func (ten *SomaTreeElemNode) setCustomProperty(
	p SomaTreeProperty) {
	ten.PropertyCustom[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setServiceProperty(
	p SomaTreeProperty) {
	ten.PropertyService[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setSystemProperty(
	p SomaTreeProperty) {
	ten.PropertySystem[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setOncallProperty(
	p SomaTreeProperty) {
	ten.PropertyOncall[p.GetID()] = p
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) syncProperty(
	childId string) {
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) checkProperty(
	propType string, propId string) bool {
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
