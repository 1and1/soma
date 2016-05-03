package somatree

import (
	"log"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (ten *SomaTreeElemNode) SetProperty(p SomaTreeProperty) {
	p.SetId(p.GetInstanceId(ten.Type, ten.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
	log.Printf("SetProperty(Node) created source instance: %s", p.GetID())
	// this property is the source instance
	p.SetInheritedFrom(ten.Id)
	p.SetInherited(false)
	p.SetSourceType(ten.Type)
	if i, e := uuid.FromString(p.GetID()); e == nil {
		p.SetSourceId(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetId(uuid.UUID{})
	ten.inheritPropertyDeep(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	switch p.GetType() {
	case "custom":
		ten.setCustomProperty(p)
	case "service":
		ten.setServiceProperty(p)
	case "system":
		ten.setSystemProperty(p)
	case "oncall":
		ten.setOncallProperty(p)
	}
	ten.actionPropertyNew(p.MakeAction())
}

func (ten *SomaTreeElemNode) inheritProperty(p SomaTreeProperty) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(ten.Type, ten.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Node) Generated: %s", f.GetID())
	}
	f.clearInstances()

	switch f.GetType() {
	case "custom":
		ten.setCustomProperty(f)
	case "service":
		ten.setServiceProperty(f)
	case "system":
		ten.setSystemProperty(f)
	case "oncall":
		ten.setOncallProperty(f)
	}
	// no inheritPropertyDeep(), nodes have no children
	ten.actionPropertyNew(f.MakeAction())
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
