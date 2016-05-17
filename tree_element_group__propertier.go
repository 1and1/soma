package somatree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

//
// Interface: SomaTreePropertier
func (teg *SomaTreeElemGroup) SetProperty(p SomaTreeProperty) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, _ := teg.checkDuplicate(p); dupe && !deleteOK {
		return // TODO: error out via FaultElement
	} else if dupe && deleteOK {
		// TODO delete inherited value
		//teg.DelProperty(prop)
		// -> passed down deleteProperty(prop) can identify by:
		//    - Inherited true
		//    - same SourceId
		//    - same SourceType
		//    - same InheritedFrom
		// TODO braucht dupeSrcID?
		return
	}
	p.SetId(p.GetInstanceId(teg.Type, teg.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
	log.Printf("SetProperty(Group) created source instance: %s", p.GetID())
	// this property is the source instance
	p.SetInheritedFrom(teg.Id)
	p.SetInherited(false)
	p.SetSourceType(teg.Type)
	if i, e := uuid.FromString(p.GetID()); e == nil {
		p.SetSourceId(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetId(uuid.UUID{})
	teg.inheritPropertyDeep(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	switch p.GetType() {
	case "custom":
		teg.setCustomProperty(p)
	case "service":
		teg.setServiceProperty(p)
	case "system":
		teg.setSystemProperty(p)
	case "oncall":
		teg.setOncallProperty(p)
	}
	teg.actionPropertyNew(p.MakeAction())
}

func (teg *SomaTreeElemGroup) inheritProperty(p SomaTreeProperty) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(teg.Type, teg.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Group) Generated: %s", f.GetID())
	}
	f.clearInstances()

	switch f.GetType() {
	case "custom":
		teg.setCustomProperty(f)
	case "service":
		teg.setServiceProperty(f)
	case "system":
		teg.setSystemProperty(f)
	case "oncall":
		teg.setOncallProperty(f)
	}
	p.SetId(uuid.UUID{})
	teg.inheritPropertyDeep(p)
	teg.actionPropertyNew(f.MakeAction())
}

func (teg *SomaTreeElemGroup) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
	log.Printf("InheritDeep Sending down: %s", p.GetID())
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
		f := teg.PropertyCustom[prop].(*PropertyCustom)
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		teg.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range teg.PropertyOncall {
		if !teg.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := teg.PropertyOncall[prop].(*PropertyOncall).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		teg.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range teg.PropertyService {
		if !teg.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := teg.PropertyService[prop].(*PropertyService).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		teg.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range teg.PropertySystem {
		if !teg.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := teg.PropertySystem[prop].(*PropertySystem).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
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

// Checks if this property is already defined on this node, and
// whether it was inherited, ie. can be deleted so it can be
// overwritten
func (teg *SomaTreeElemGroup) checkDuplicate(p SomaTreeProperty) (
	bool, bool, SomaTreeProperty) {
	var dupe, deleteOK bool
	var prop SomaTreeProperty

propswitch:
	switch p.GetType() {
	case "custom":
		for _, pVal := range teg.PropertyCustom {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "service":
		for _, pVal := range teg.PropertyService {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "oncall":
		for _, pVal := range teg.PropertyOncall {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "system":
		for _, pVal := range teg.PropertySystem {
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
