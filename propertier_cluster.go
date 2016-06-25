package somatree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Propertier` interface

//
// Propertier:> Add Property

func (tec *SomaTreeElemCluster) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, _ := tec.checkDuplicate(p); dupe && !deleteOK {
		log.Printf("cluster.SetProperty() detected hard duplicate")
		return // TODO: error out via FaultElement
	} else if dupe && deleteOK {
		// TODO delete inherited value
		// tec.DelProperty(prop)
		log.Printf("cluster.SetProperty() detected soft duplicate")
		return
	}
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
	tec.setPropertyOnChildren(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	tec.addProperty(p)
	tec.actionPropertyNew(p.MakeAction())
}

func (tec *SomaTreeElemCluster) setPropertyInherited(p Property) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(tec.Type, tec.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Cluster) Generated: %s", f.GetID())
	}
	f.clearInstances()

	tec.addProperty(p)
	p.SetId(uuid.UUID{})
	tec.setPropertyOnChildren(p)
	tec.actionPropertyNew(f.MakeAction())
}

func (tec *SomaTreeElemCluster) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	log.Printf("InheritDeep Sending down: %s", p.GetID())
	for child, _ := range tec.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			tec.Children[c].setPropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (tec *SomaTreeElemCluster) addProperty(p Property) {
	switch p.GetType() {
	case `custom`:
		tec.PropertyCustom[p.GetID()] = p
	case `system`:
		tec.PropertySystem[p.GetID()] = p
	case `service`:
		tec.PropertyService[p.GetID()] = p
	case `oncall`:
		tec.PropertyOncall[p.GetID()] = p
	}
}

//
// Propertier:> Update Property

func (tec *SomaTreeElemCluster) UpdateProperty(p Property) {
	if !tec.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		return // XXX faultChannel
	}

	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	tec.switchProperty(f)
	tec.updatePropertyOnChildren(p)
}

func (tec *SomaTreeElemCluster) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	tec.switchProperty(f)
	tec.updatePropertyOnChildren(p)
}

func (tec *SomaTreeElemCluster) updatePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			tec.Children[c].updatePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (tec *SomaTreeElemCluster) switchProperty(p Property) {
	updId, _ := uuid.FromString(tec.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	))
	p.SetId(updId)
	tec.addProperty(p)
	tec.actionPropertyUpdate(p.MakeAction())
}

//
// Propertier:> Delete Property

func (tec *SomaTreeElemCluster) DeleteProperty(p Property) {
	if !tec.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		return // XXX faultChannel
	}

	tec.rmProperty(p)
	tec.deletePropertyOnChildren(p)
}

func (tec *SomaTreeElemCluster) deletePropertyInherited(p Property) {
	tec.rmProperty(p)
	tec.deletePropertyOnChildren(p)
}

func (tec *SomaTreeElemCluster) deletePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			tec.Children[c].deletePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (tec *SomaTreeElemCluster) rmProperty(p Property) {
	delId, _ := uuid.FromString(tec.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	))
	p.SetId(delId)
	tec.actionPropertyDelete(p.MakeAction())

	switch p.GetType() {
	case `custom`:
		delete(tec.PropertyCustom, delId.String())
	case `service`:
		delete(tec.PropertyService, delId.String())
	case `system`:
		delete(tec.PropertySystem, delId.String())
	case `oncall`:
		delete(tec.PropertyOncall, delId.String())
	}
}

//
// Propertier:> Utility

//
func (tec *SomaTreeElemCluster) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := tec.PropertyCustom[id]; !ok {
			return false
		}
		return tec.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := tec.PropertyService[id]; !ok {
			return false
		}
		return tec.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := tec.PropertySystem[id]; !ok {
			return false
		}
		return tec.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := tec.PropertyOncall[id]; !ok {
			return false
		}
		return tec.PropertyOncall[id].GetSourceInstance() == id
	default:
		return false
	}
}

//
func (tec *SomaTreeElemCluster) findIdForSource(source, prop string) string {
	switch prop {
	case `custom`:
		for id, _ := range tec.PropertyCustom {
			if tec.PropertyCustom[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `system`:
		for id, _ := range tec.PropertySystem {
			if tec.PropertySystem[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `service`:
		for id, _ := range tec.PropertyService {
			if tec.PropertyService[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `oncall`:
		for id, _ := range tec.PropertyOncall {
			if tec.PropertyOncall[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	}
	return ``
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (tec *SomaTreeElemCluster) syncProperty(childId string) {
customloop:
	for prop, _ := range tec.PropertyCustom {
		if !tec.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := tec.PropertyCustom[prop].(*PropertyCustom).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		tec.Children[childId].setPropertyInherited(f)
	}
oncallloop:
	for prop, _ := range tec.PropertyOncall {
		if !tec.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := tec.PropertyOncall[prop].(*PropertyOncall).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		tec.Children[childId].setPropertyInherited(f)
	}
serviceloop:
	for prop, _ := range tec.PropertyService {
		if !tec.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := tec.PropertyService[prop].(*PropertyService).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		tec.Children[childId].setPropertyInherited(f)
	}
systemloop:
	for prop, _ := range tec.PropertySystem {
		if !tec.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := tec.PropertySystem[prop].(*PropertySystem).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		tec.Children[childId].setPropertyInherited(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (tec *SomaTreeElemCluster) checkProperty(propType string, propId string) bool {
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

// Checks if this property is already defined on this node, and
// whether it was inherited, ie. can be deleted so it can be
// overwritten
func (tec *SomaTreeElemCluster) checkDuplicate(p Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

propswitch:
	switch p.GetType() {
	case "custom":
		for _, pVal := range tec.PropertyCustom {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "service":
		for _, pVal := range tec.PropertyService {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "oncall":
		for _, pVal := range tec.PropertyOncall {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "system":
		for _, pVal := range tec.PropertySystem {
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
