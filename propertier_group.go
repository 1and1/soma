package tree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Propertier` interface

//
// Propertier:> Add Property

func (teg *Group) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, prop := teg.checkDuplicate(p); dupe && !deleteOK {
		teg.Fault.Error <- &Error{Action: `duplicate_set_property`}
		return
	} else if dupe && deleteOK {
		srcUUID, _ := uuid.FromString(prop.GetSourceInstance())
		switch prop.GetType() {
		case `custom`:
			teg.deletePropertyInherited(&PropertyCustom{
				SourceId: srcUUID,
			})
		case `service`:
			teg.deletePropertyInherited(&PropertyService{
				SourceId: srcUUID,
			})
		case `system`:
			teg.deletePropertyInherited(&PropertySystem{
				SourceId: srcUUID,
			})
		case `oncall`:
			teg.deletePropertyInherited(&PropertyOncall{
				SourceId: srcUUID,
			})
		}
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
	teg.setPropertyOnChildren(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	teg.addProperty(p)
	teg.actionPropertyNew(p.MakeAction())
}

func (teg *Group) setPropertyInherited(p Property) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(teg.Type, teg.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Group) Generated: %s", f.GetID())
	}
	f.clearInstances()

	if !f.GetIsInherited() {
		panic(`not inherited`)
	}
	teg.addProperty(f)
	p.SetId(uuid.UUID{})
	teg.setPropertyOnChildren(p)
	teg.actionPropertyNew(f.MakeAction())
}

func (teg *Group) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	log.Printf("InheritDeep Sending down: %s", p.GetID())
	for child, _ := range teg.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teg.Children[c].setPropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teg *Group) addProperty(p Property) {
	switch p.GetType() {
	case `custom`:
		teg.PropertyCustom[p.GetID()] = p
	case `system`:
		teg.PropertySystem[p.GetID()] = p
	case `service`:
		teg.PropertyService[p.GetID()] = p
	case `oncall`:
		teg.PropertyOncall[p.GetID()] = p
	}
}

//
// Propertier:> Update Property

func (teg *Group) UpdateProperty(p Property) {
	if !teg.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		return // XXX faultChannel
	}

	// keep a copy for ourselves, no shared pointers
	p.SetInheritedFrom(teg.Id)
	p.SetSourceType(teg.Type)
	p.SetInherited(true)
	f := p.Clone()
	f.SetInherited(false)
	teg.switchProperty(f)
	teg.updatePropertyOnChildren(p)
}

func (teg *Group) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	if !f.GetIsInherited() {
		panic(`not inherited`)
	}
	teg.switchProperty(f)
	teg.updatePropertyOnChildren(p)
}

func (teg *Group) updatePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teg.Children[c].updatePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teg *Group) switchProperty(p Property) {
	updId, _ := uuid.FromString(teg.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	))
	p.SetId(updId)
	teg.addProperty(p)
	teg.actionPropertyUpdate(p.MakeAction())
}

//
// Propertier:> Delete Property

func (teg *Group) DeleteProperty(p Property) {
	if !teg.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		return // XXX faultChannel
	}

	teg.rmProperty(p)
	teg.deletePropertyOnChildren(p)
}

func (teg *Group) deletePropertyInherited(p Property) {
	teg.rmProperty(p)
	teg.deletePropertyOnChildren(p)
}

func (teg *Group) deletePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teg.Children[c].deletePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teg *Group) rmProperty(p Property) {
	delId := teg.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)

	switch p.GetType() {
	case `custom`:
		teg.actionPropertyDelete(
			teg.PropertyCustom[delId].MakeAction(),
		)
		delete(teg.PropertyCustom, delId)
	case `service`:
		teg.actionPropertyDelete(
			teg.PropertyService[delId].MakeAction(),
		)
		delete(teg.PropertyService, delId)
	case `system`:
		teg.actionPropertyDelete(
			teg.PropertySystem[delId].MakeAction(),
		)
		delete(teg.PropertySystem, delId)
	case `oncall`:
		teg.actionPropertyDelete(
			teg.PropertyOncall[delId].MakeAction(),
		)
		delete(teg.PropertyOncall, delId)
	}
}

//
// Propertier:> Utility

//
func (teg *Group) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := teg.PropertyCustom[id]; !ok {
			return false
		}
		return teg.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := teg.PropertyService[id]; !ok {
			return false
		}
		return teg.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := teg.PropertySystem[id]; !ok {
			return false
		}
		return teg.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := teg.PropertyOncall[id]; !ok {
			return false
		}
		return teg.PropertyOncall[id].GetSourceInstance() == id
	default:
		return false
	}
}

//
func (teg *Group) findIdForSource(source, prop string) string {
	switch prop {
	case `custom`:
		for id, _ := range teg.PropertyCustom {
			if teg.PropertyCustom[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `system`:
		for id, _ := range teg.PropertySystem {
			if teg.PropertySystem[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `service`:
		for id, _ := range teg.PropertyService {
			if teg.PropertyService[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `oncall`:
		for id, _ := range teg.PropertyOncall {
			if teg.PropertyOncall[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	}
	return ``
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (teg *Group) syncProperty(childId string) {
customloop:
	for prop, _ := range teg.PropertyCustom {
		if !teg.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := teg.PropertyCustom[prop].(*PropertyCustom)
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		teg.Children[childId].setPropertyInherited(f)
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
		teg.Children[childId].setPropertyInherited(f)
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
		teg.Children[childId].setPropertyInherited(f)
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
		teg.Children[childId].setPropertyInherited(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (teg *Group) checkProperty(propType string, propId string) bool {
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
func (teg *Group) checkDuplicate(p Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

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
