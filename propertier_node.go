package somatree

import (
	"log"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (ten *SomaTreeElemNode) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, _ := ten.checkDuplicate(p); dupe && !deleteOK {
		log.Printf("node.SetProperty() detected hard duplicate")
		return // TODO: error out via FaultElement
	} else if dupe && deleteOK {
		// TODO delete inherited value
		// ten.DelProperty(prop)
		log.Printf("node.SetProperty() detected soft duplicate")
		return
	}
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

func (ten *SomaTreeElemNode) inheritProperty(p Property) {
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
	p Property) {
}

func (ten *SomaTreeElemNode) setCustomProperty(
	p Property) {
	ten.PropertyCustom[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setServiceProperty(
	p Property) {
	ten.PropertyService[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setSystemProperty(
	p Property) {
	ten.PropertySystem[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setOncallProperty(
	p Property) {
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

// Checks if this property is already defined on this node, and
// whether it was inherited, ie. can be deleted so it can be
// overwritten
func (ten *SomaTreeElemNode) checkDuplicate(p Property) (
	bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

propswitch:
	switch p.GetType() {
	case "custom":
		for _, pVal := range ten.PropertyCustom {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "service":
		for _, pVal := range ten.PropertyService {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "oncall":
		for _, pVal := range ten.PropertyOncall {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "system":
		for _, pVal := range ten.PropertySystem {
			// tags are only dupes if the value is the same as well
			if p.GetKey() != `tag` {
				dupe, deleteOK, prop = isDupe(pVal, p)
				if dupe {
					break propswitch
				}
			} else if p.GetValue() == pVal.GetValue() {
				// tag and same value, can be a dupe
				dupe, deleteOK, prop = isDupe(pVal, p)
				if dupe {
					break propswitch
				}
			}
			// tag + different value => pass
		}
	default:
		// trigger error path
		dupe = true
		deleteOK = false
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
