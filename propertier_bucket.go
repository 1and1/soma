package tree

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Propertier` interface

//
// Propertier:> Add Property

func (teb *Bucket) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, prop := teb.checkDuplicate(p); dupe && !deleteOK {
		teb.Fault.Error <- &Error{Action: `duplicate_set_property`}
		return
	} else if dupe && deleteOK {
		srcUUID, _ := uuid.FromString(prop.GetSourceInstance())
		switch prop.GetType() {
		case `custom`:
			teb.deletePropertyInherited(&PropertyCustom{
				SourceId: srcUUID,
			})
		case `service`:
			teb.deletePropertyInherited(&PropertyService{
				SourceId: srcUUID,
			})
		case `system`:
			teb.deletePropertyInherited(&PropertySystem{
				SourceId: srcUUID,
			})
		case `oncall`:
			teb.deletePropertyInherited(&PropertyOncall{
				SourceId: srcUUID,
			})
		}
	}
	p.SetId(p.GetInstanceId(teb.Type, teb.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
	log.Printf("SetProperty(Bucket) created source instance: %s", p.GetID())
	// this property is the source instance
	p.SetInheritedFrom(teb.Id)
	p.SetInherited(false)
	p.SetSourceType(teb.Type)
	if i, e := uuid.FromString(p.GetID()); e == nil {
		p.SetSourceId(i)
	}
	// send a scrubbed copy down
	f := p.Clone()
	f.SetInherited(true)
	f.SetId(uuid.UUID{})
	teb.setPropertyOnChildren(f)
	// scrub instance startup information prior to storing
	p.clearInstances()
	teb.addProperty(p)
	teb.actionPropertyNew(p.MakeAction())
}

func (teb *Bucket) setPropertyInherited(p Property) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(teb.Type, teb.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
		log.Printf("Inherit (Bucket) Generated: %s", f.GetID())
	}
	f.clearInstances()

	if !f.GetIsInherited() {
		teb.Fault.Error <- &Error{
			Action: `bucket.setPropertyInherited on inherited=false`}
		return
	}
	teb.addProperty(f)
	p.SetId(uuid.UUID{})
	teb.setPropertyOnChildren(p)
	teb.actionPropertyNew(f.MakeAction())
}

func (teb *Bucket) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	log.Printf("InheritDeep Sending down: %s", p.GetID())
	for child, _ := range teb.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teb.Children[c].setPropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teb *Bucket) addProperty(p Property) {
	switch p.GetType() {
	case `custom`:
		teb.PropertyCustom[p.GetID()] = p
	case `system`:
		teb.PropertySystem[p.GetID()] = p
	case `service`:
		teb.PropertyService[p.GetID()] = p
	case `oncall`:
		teb.PropertyOncall[p.GetID()] = p
	default:
		teb.Fault.Error <- &Error{Action: `bucket.addProperty unknown type`}
	}
}

//
// Propertier:> Update Property

func (teb *Bucket) UpdateProperty(p Property) {
	if !teb.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		teb.Fault.Error <- &Error{Action: `update_property_on_non_source`}
		return
	}

	// keep a copy for ourselves, no shared pointers
	p.SetInheritedFrom(teb.Id)
	p.SetSourceType(teb.Type)
	p.SetInherited(true)
	f := p.Clone()
	f.SetInherited(false)
	teb.switchProperty(f)
	teb.updatePropertyOnChildren(p)
}

func (teb *Bucket) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	if !f.GetIsInherited() {
		teb.Fault.Error <- &Error{
			Action: `bucket.updatePropertyInherited on inherited=false`}
		return
	}
	teb.switchProperty(f)
	teb.updatePropertyOnChildren(p)
}

func (teb *Bucket) updatePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teb.Children[c].updatePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teb *Bucket) switchProperty(p Property) {
	updId, _ := uuid.FromString(teb.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	))
	p.SetId(updId)
	teb.addProperty(p)
	teb.actionPropertyUpdate(p.MakeAction())
}

//
// Propertier:> Delete Property

func (teb *Bucket) DeleteProperty(p Property) {
	if !teb.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		teb.Fault.Error <- &Error{Action: `bucket.DeleteProperty on !source`}
		return
	}

	teb.rmProperty(p)
	teb.deletePropertyOnChildren(p)
}

func (teb *Bucket) deletePropertyInherited(p Property) {
	teb.rmProperty(p)
	teb.deletePropertyOnChildren(p)
}

func (teb *Bucket) deletePropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			teb.Children[c].deletePropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (teb *Bucket) rmProperty(p Property) {
	delId := teb.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delId == `` {
		teb.Fault.Error <- &Error{
			Action: `bucket.rmProperty property not found`}
		return
	}

	switch p.GetType() {
	case `custom`:
		teb.actionPropertyDelete(
			teb.PropertyCustom[delId].MakeAction(),
		)
		delete(teb.PropertyCustom, delId)
	case `service`:
		teb.actionPropertyDelete(
			teb.PropertyService[delId].MakeAction(),
		)
		delete(teb.PropertyService, delId)
	case `system`:
		teb.actionPropertyDelete(
			teb.PropertySystem[delId].MakeAction(),
		)
		delete(teb.PropertySystem, delId)
	case `oncall`:
		teb.actionPropertyDelete(
			teb.PropertyOncall[delId].MakeAction(),
		)
		delete(teb.PropertyOncall, delId)
	default:
		teb.Fault.Error <- &Error{Action: `bucket.rmProperty unknown type`}
	}
}

//
// Propertier:> Utility

// used to verify this is a source instance
func (teb *Bucket) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := teb.PropertyCustom[id]; !ok {
			return false
		}
		return teb.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := teb.PropertyService[id]; !ok {
			return false
		}
		return teb.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := teb.PropertySystem[id]; !ok {
			return false
		}
		return teb.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := teb.PropertyOncall[id]; !ok {
			return false
		}
		return teb.PropertyOncall[id].GetSourceInstance() == id
	default:
		return false
	}
}

//
func (teb *Bucket) findIdForSource(source, prop string) string {
	switch prop {
	case `custom`:
		for id, _ := range teb.PropertyCustom {
			if teb.PropertyCustom[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `system`:
		for id, _ := range teb.PropertySystem {
			if teb.PropertySystem[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `service`:
		for id, _ := range teb.PropertyService {
			if teb.PropertyService[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	case `oncall`:
		for id, _ := range teb.PropertyOncall {
			if teb.PropertyOncall[id].GetSourceInstance() != source {
				continue
			}
			return id
		}
	}
	return ``
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (teb *Bucket) syncProperty(childId string) {
customloop:
	for prop, _ := range teb.PropertyCustom {
		if !teb.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := teb.PropertyCustom[prop].(*PropertyCustom).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		teb.Children[childId].setPropertyInherited(f)
	}
oncallloop:
	for prop, _ := range teb.PropertyOncall {
		if !teb.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := teb.PropertyOncall[prop].(*PropertyOncall).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		teb.Children[childId].setPropertyInherited(f)
	}
serviceloop:
	for prop, _ := range teb.PropertyService {
		if !teb.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := teb.PropertyService[prop].(*PropertyService).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		teb.Children[childId].setPropertyInherited(f)
	}
systemloop:
	for prop, _ := range teb.PropertySystem {
		if !teb.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := teb.PropertySystem[prop].(*PropertySystem).Clone()
		f.SetInherited(true)
		f.SetId(uuid.UUID{})
		f.clearInstances()
		teb.Children[childId].setPropertyInherited(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (teb *Bucket) checkProperty(propType string, propId string) bool {
	switch propType {
	case "custom":
		if _, ok := teb.PropertyCustom[propId]; ok {
			return true
		}
	case "service":
		if _, ok := teb.PropertyService[propId]; ok {
			return true
		}
	case "system":
		if _, ok := teb.PropertySystem[propId]; ok {
			return true
		}
	case "oncall":
		if _, ok := teb.PropertyOncall[propId]; ok {
			return true
		}
	}
	return false
}

// Checks if this property is already defined on this node, and
// whether it was inherited, ie. can be deleted so it can be
// overwritten
func (teb *Bucket) checkDuplicate(p Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

propswitch:
	switch p.GetType() {
	case "custom":
		for _, pVal := range teb.PropertyCustom {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "service":
		for _, pVal := range teb.PropertyService {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "oncall":
		for _, pVal := range teb.PropertyOncall {
			dupe, deleteOK, prop = isDupe(pVal, p)
			if dupe {
				break propswitch
			}
		}
	case "system":
		for _, pVal := range teb.PropertySystem {
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
