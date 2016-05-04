package somatree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (ter *SomaTreeElemRepository) SetProperty(p SomaTreeProperty) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, _ := ter.checkDuplicate(p); dupe && !deleteOK {
		return // TODO: error out via FaultElement
	} else if dupe && deleteOK {
		// TODO delete inherited value
		// ter.DelProperty(prop)
		return
	}
	p.SetId(p.GetInstanceId(ter.Type, ter.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
	log.Printf("SetProperty(Repository) created source instance: %s", p.GetID())
	// this property is the source instance
	p.SetInheritedFrom(ter.Id)
	p.SetInherited(false)
	p.SetSourceType(ter.Type)
	if i, e := uuid.FromString(p.GetID()); e == nil {
		p.SetSourceId(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetId(uuid.UUID{})
	ter.inheritPropertyDeep(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	switch p.GetType() {
	case "custom":
		ter.setCustomProperty(p)
	case "service":
		ter.setServiceProperty(p)
	case "system":
		ter.setSystemProperty(p)
	case "oncall":
		ter.setOncallProperty(p)
	}
	ter.actionPropertyNew(p.MakeAction())
}

func (ter *SomaTreeElemRepository) inheritProperty(p SomaTreeProperty) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(ter.Type, ter.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Repository) Generated: %s", f.GetID())
	}
	f.clearInstances()

	switch f.GetType() {
	case "custom":
		ter.setCustomProperty(f)
	case "service":
		ter.setServiceProperty(f)
	case "system":
		ter.setSystemProperty(f)
	case "oncall":
		ter.setOncallProperty(f)
	}
	p.SetId(uuid.UUID{})
	ter.inheritPropertyDeep(p)
	ter.actionPropertyNew(f.MakeAction())
}

func (ter *SomaTreeElemRepository) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
	log.Printf("InheritDeep Sending down: %s", p.GetID())
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func(stp SomaTreeProperty) {
			defer wg.Done()
			ter.Children[c].inheritProperty(stp)
		}(p)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) setCustomProperty(
	p SomaTreeProperty) {
	ter.PropertyCustom[p.GetID()] = p
}

func (ter *SomaTreeElemRepository) setServiceProperty(
	p SomaTreeProperty) {
	ter.PropertyService[p.GetID()] = p
}

func (ter *SomaTreeElemRepository) setSystemProperty(
	p SomaTreeProperty) {
	ter.PropertySystem[p.GetID()] = p
}

func (ter *SomaTreeElemRepository) setOncallProperty(
	p SomaTreeProperty) {
	ter.PropertyOncall[p.GetID()] = p
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (ter *SomaTreeElemRepository) syncProperty(
	childId string) {
customloop:
	for prop, _ := range ter.PropertyCustom {
		if !ter.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := new(PropertyCustom)
		*f = *ter.PropertyCustom[prop].(*PropertyCustom)
		f.Inherited = true
		ter.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range ter.PropertyOncall {
		if !ter.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(PropertyOncall)
		*f = *ter.PropertyOncall[prop].(*PropertyOncall)
		f.Inherited = true
		ter.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range ter.PropertyService {
		if !ter.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(PropertyService)
		*f = *ter.PropertyService[prop].(*PropertyService)
		f.Inherited = true
		ter.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range ter.PropertySystem {
		if !ter.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(PropertySystem)
		*f = *ter.PropertySystem[prop].(*PropertySystem)
		f.Inherited = true
		ter.Children[childId].inheritProperty(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (ter *SomaTreeElemRepository) checkProperty(
	propType string, propId string) bool {
	switch propType {
	case "custom":
		if _, ok := ter.PropertyCustom[propId]; ok {
			return true
		}
	case "service":
		if _, ok := ter.PropertyService[propId]; ok {
			return true
		}
	case "system":
		if _, ok := ter.PropertySystem[propId]; ok {
			return true
		}
	case "oncall":
		if _, ok := ter.PropertyOncall[propId]; ok {
			return true
		}
	}
	return false
}

// Checks if this property is already defined on this node, and
// whether it was inherited, ie. can be deleted so it can be
// overwritten
func (ter *SomaTreeElemRepository) checkDuplicate(p SomaTreeProperty) (
	bool, bool, SomaTreeProperty) {
	var dupe, deleteOK bool
	var prop SomaTreeProperty

propswitch:
	switch p.GetType() {
	case "custom":
		for _, pVal := range ter.PropertyCustom {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "service":
		for _, pVal := range ter.PropertyService {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "oncall":
		for _, pVal := range ter.PropertyOncall {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "system":
		for _, pVal := range ter.PropertySystem {
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
