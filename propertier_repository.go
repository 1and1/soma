package somatree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Propertier` interface

//
// Propertier:> Add Property

func (ter *SomaTreeElemRepository) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, _ := ter.checkDuplicate(p); dupe && !deleteOK {
		log.Printf("repository.SetProperty() detected hard duplicate")
		return // TODO: error out via FaultElement
	} else if dupe && deleteOK {
		// TODO delete inherited value
		// ter.DelProperty(prop)
		log.Printf("repository.SetProperty() detected soft duplicate")
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
	ter.setPropertyOnChildren(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	ter.addProperty(p)
	ter.actionPropertyNew(p.MakeAction())
}

func (ter *SomaTreeElemRepository) setPropertyInherited(p Property) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(ter.Type, ter.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Repository) Generated: %s", f.GetID())
	}
	f.clearInstances()

	ter.addProperty(f)
	p.SetId(uuid.UUID{})
	ter.setPropertyOnChildren(p)
	ter.actionPropertyNew(f.MakeAction())
}

func (ter *SomaTreeElemRepository) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	log.Printf("InheritDeep Sending down: %s", p.GetID())
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func(stp Property) {
			defer wg.Done()
			ter.Children[c].setPropertyInherited(stp)
		}(p)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) addProperty(p Property) {
	switch p.GetType() {
	case `custom`:
		ter.PropertyCustom[p.GetID()] = p
	case `system`:
		ter.PropertySystem[p.GetID()] = p
	case `service`:
		ter.PropertyService[p.GetID()] = p
	case `oncall`:
		ter.PropertyOncall[p.GetID()] = p
	}
}

//
// Propertier:> Update Property

func (ter *SomaTreeElemRepository) UpdateProperty(p Property) {
	if !ter.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		return // XXX faultChannel
	}

	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	ter.switchProperty(f)
	ter.updatePropertyOnChildren(p)
}

func (ter *SomaTreeElemRepository) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	ter.switchProperty(f)
	ter.updatePropertyOnChildren(p)
}

func (ter *SomaTreeElemRepository) updatePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			ter.Children[c].updatePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) switchProperty(p Property) {
	updId, _ := uuid.FromString(ter.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	))
	p.SetId(updId)
	ter.addProperty(p)
	ter.actionPropertyUpdate(p.MakeAction())
}

//
// Propertier:> Delete Property

func (ter *SomaTreeElemRepository) DeleteProperty(p Property) {
	if !ter.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		return // XXX faultChannel
	}

	ter.rmProperty(p)
	ter.deletePropertyOnChildren(p)
}

func (ter *SomaTreeElemRepository) deletePropertyInherited(p Property) {
	ter.rmProperty(p)
	ter.deletePropertyOnChildren(p)
}

func (ter *SomaTreeElemRepository) deletePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			ter.Children[c].deletePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) rmProperty(p Property) {
	delId, _ := uuid.FromString(ter.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	))
	p.SetId(delId)
	ter.actionPropertyDelete(p.MakeAction())

	switch p.GetType() {
	case `custom`:
		delete(ter.PropertyCustom, delId.String())
	case `service`:
		delete(ter.PropertyService, delId.String())
	case `system`:
		delete(ter.PropertySystem, delId.String())
	case `oncall`:
		delete(ter.PropertyOncall, delId.String())
	}
}

//
// Propertier:> Utility

//
func (ter *SomaTreeElemRepository) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := ter.PropertyCustom[id]; !ok {
			return false
		}
		return ter.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := ter.PropertyService[id]; !ok {
			return false
		}
		return ter.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := ter.PropertySystem[id]; !ok {
			return false
		}
		return ter.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := ter.PropertyOncall[id]; !ok {
			return false
		}
		return ter.PropertyOncall[id].GetSourceInstance() == id
	default:
		return false
	}
}

func (ter *SomaTreeElemRepository) findIdForSource(source, prop string) string {
	switch prop {
	case `custom`:
		for id, _ := range ter.PropertyCustom {
			if ter.PropertyCustom[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `system`:
		for id, _ := range ter.PropertySystem {
			if ter.PropertySystem[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `service`:
		for id, _ := range ter.PropertyService {
			if ter.PropertyService[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `oncall`:
		for id, _ := range ter.PropertyOncall {
			if ter.PropertyOncall[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	}
	return ``
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (ter *SomaTreeElemRepository) syncProperty(childId string) {
customloop:
	for prop, _ := range ter.PropertyCustom {
		if !ter.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := ter.PropertyCustom[prop].(*PropertyCustom).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		ter.Children[childId].setPropertyInherited(f)
	}
oncallloop:
	for prop, _ := range ter.PropertyOncall {
		if !ter.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := ter.PropertyOncall[prop].(*PropertyOncall).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		ter.Children[childId].setPropertyInherited(f)
	}
serviceloop:
	for prop, _ := range ter.PropertyService {
		if !ter.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := ter.PropertyService[prop].(*PropertyService).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		ter.Children[childId].setPropertyInherited(f)
	}
systemloop:
	for prop, _ := range ter.PropertySystem {
		if !ter.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := ter.PropertySystem[prop].(*PropertySystem).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		ter.Children[childId].setPropertyInherited(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (ter *SomaTreeElemRepository) checkProperty(propType string, propId string) bool {
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
func (ter *SomaTreeElemRepository) checkDuplicate(p Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

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
