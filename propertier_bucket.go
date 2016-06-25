package tree

import (
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
			cstUUID, _ := uuid.FromString(prop.GetKey())
			teb.deletePropertyInherited(&PropertyCustom{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				CustomId:  cstUUID,
				Key:       prop.(*PropertyCustom).GetKeyField(),
				Value:     prop.(*PropertyCustom).GetValueField(),
			})
		case `service`:
			// GetValue for serviceproperty returns the uuid to never
			// match, we do not set it
			teb.deletePropertyInherited(&PropertyService{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Service:   prop.GetKey(),
			})
		case `system`:
			teb.deletePropertyInherited(&PropertySystem{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Key:       prop.GetKey(),
				Value:     prop.GetValue(),
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(prop.GetKey())
			teb.deletePropertyInherited(&PropertyOncall{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				OncallId:  oncUUID,
				Name:      prop.(*PropertyOncall).GetName(),
				Number:    prop.(*PropertyOncall).GetNumber(),
			})
		}
	}
	p.SetId(p.GetInstanceId(teb.Type, teb.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
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
	if f.hasInheritance() {
		teb.setPropertyOnChildren(f)
	}
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
	}
	f.clearInstances()

	if !f.GetIsInherited() {
		teb.Fault.Error <- &Error{
			Action: `bucket.setPropertyInherited on inherited=false`}
		return
	}
	if dupe, deleteOK, _ := teb.checkDuplicate(p); dupe && deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate, but we are not the
		// source of the duplicate -> corrupt tree
		teb.Fault.Error <- &Error{
			Action: `bucket.setPropertyInherited corruption detected`}
		return
	} else if dupe && !deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate; we have a locally
		// set property -> stop inheritance, no error
		return
	}
	teb.addProperty(f)
	p.SetId(uuid.UUID{})
	teb.setPropertyOnChildren(p)
	teb.actionPropertyNew(f.MakeAction())
}

func (teb *Bucket) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
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
	if teb.switchProperty(f) {
		teb.updatePropertyOnChildren(p)
	}
}

func (teb *Bucket) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	if !f.GetIsInherited() {
		teb.Fault.Error <- &Error{
			Action: `bucket.updatePropertyInherited on inherited=false`}
		return
	}
	if teb.switchProperty(f) {
		teb.updatePropertyOnChildren(p)
	}
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

func (teb *Bucket) switchProperty(p Property) bool {
	uid := teb.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if uid == `` {
		// we do not have the property for which we received an update
		if dupe, deleteOK, _ := teb.checkDuplicate(p); dupe && !deleteOK {
			// the update is duplicate to an property for which we
			// have the source instance, ie we just received an update
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}
		teb.Fault.Error <- &Error{
			Action: `bucket.switchProperty property not found`}
		return false
	}
	updId, _ := uuid.FromString(uid)
	p.SetId(updId)
	curr := teb.getCurrentProperty(p)
	if curr == nil {
		return false
	}
	teb.addProperty(p)
	teb.actionPropertyUpdate(p.MakeAction())

	if !p.hasInheritance() && curr.hasInheritance() {
		// replacing inheritance with !inheritance:
		// call deletePropertyOnChildren(curr) to clean up
		srcUUID, _ := uuid.FromString(curr.GetSourceInstance())
		switch curr.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(curr.GetKey())
			teb.deletePropertyOnChildren(&PropertyCustom{
				SourceId:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				CustomId:    cstUUID,
				Key:         curr.(*PropertyCustom).GetKeyField(),
				Value:       curr.(*PropertyCustom).GetValueField(),
				Inheritance: true,
			})
		case `service`:
			// GetValue for serviceproperty returns the uuid to never
			// match, we do not set it
			teb.deletePropertyOnChildren(&PropertyService{
				SourceId:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Service:     curr.GetKey(),
				Inheritance: true,
			})
		case `system`:
			teb.deletePropertyOnChildren(&PropertySystem{
				SourceId:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Key:         curr.GetKey(),
				Value:       curr.GetValue(),
				Inheritance: true,
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(curr.GetKey())
			teb.deletePropertyOnChildren(&PropertyOncall{
				SourceId:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				OncallId:    oncUUID,
				Name:        curr.(*PropertyOncall).GetName(),
				Number:      curr.(*PropertyOncall).GetNumber(),
				Inheritance: true,
			})
		}
	}
	if p.hasInheritance() && !curr.hasInheritance() {
		// replacing !inheritance with inheritance:
		// call setPropertyonChildren(p) to propagate
		t := p.Clone()
		t.SetInherited(true)
		teb.setPropertyOnChildren(t)
	}
	return p.hasInheritance() && curr.hasInheritance()
}

func (teb *Bucket) getCurrentProperty(p Property) Property {
	switch p.GetType() {
	case `custom`:
		return teb.PropertyCustom[p.GetID()].Clone()
	case `system`:
		return teb.PropertySystem[p.GetID()].Clone()
	case `service`:
		return teb.PropertyService[p.GetID()].Clone()
	case `oncall`:
		return teb.PropertyOncall[p.GetID()].Clone()
	}
	teb.Fault.Error <- &Error{
		Action: `bucket.getCurrentProperty unknown type`}
	return nil
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

	p.SetInherited(false)
	if teb.rmProperty(p) {
		p.SetInherited(true)
		teb.deletePropertyOnChildren(p)
	}
}

func (teb *Bucket) deletePropertyInherited(p Property) {
	if teb.rmProperty(p) {
		teb.deletePropertyOnChildren(p)
	}
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

func (teb *Bucket) deletePropertyAllInherited() {
	for _, p := range teb.PropertyCustom {
		if !p.GetIsInherited() {
			continue
		}
		teb.deletePropertyInherited(p.Clone())
	}
	for _, p := range teb.PropertySystem {
		if !p.GetIsInherited() {
			continue
		}
		teb.deletePropertyInherited(p.Clone())
	}
	for _, p := range teb.PropertyService {
		if !p.GetIsInherited() {
			continue
		}
		teb.deletePropertyInherited(p.Clone())
	}
	for _, p := range teb.PropertyOncall {
		if !p.GetIsInherited() {
			continue
		}
		teb.deletePropertyInherited(p.Clone())
	}
}

func (teb *Bucket) deletePropertyAllLocal() {
	for _, p := range teb.PropertyCustom {
		if p.GetIsInherited() {
			continue
		}
		teb.DeleteProperty(p.Clone())
	}
	for _, p := range teb.PropertySystem {
		if p.GetIsInherited() {
			continue
		}
		teb.DeleteProperty(p.Clone())
	}
	for _, p := range teb.PropertyService {
		if p.GetIsInherited() {
			continue
		}
		teb.DeleteProperty(p.Clone())
	}
	for _, p := range teb.PropertyOncall {
		if p.GetIsInherited() {
			continue
		}
		teb.DeleteProperty(p.Clone())
	}
}

func (teb *Bucket) rmProperty(p Property) bool {
	delId := teb.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delId == `` {
		// we do not have the property for which we received a delete
		if dupe, deleteOK, _ := teb.checkDuplicate(p); dupe && !deleteOK {
			// the delete is duplicate to a property for which we
			// have the source instance, ie we just received a delete
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}

		teb.Fault.Error <- &Error{
			Action: `bucket.rmProperty property not found`}
		return false
	}

	hasInheritance := false
	switch p.GetType() {
	case `custom`:
		teb.actionPropertyDelete(
			teb.PropertyCustom[delId].MakeAction(),
		)
		hasInheritance = teb.PropertyCustom[delId].hasInheritance()
		delete(teb.PropertyCustom, delId)
	case `service`:
		teb.actionPropertyDelete(
			teb.PropertyService[delId].MakeAction(),
		)
		hasInheritance = teb.PropertyService[delId].hasInheritance()
		delete(teb.PropertyService, delId)
	case `system`:
		teb.actionPropertyDelete(
			teb.PropertySystem[delId].MakeAction(),
		)
		hasInheritance = teb.PropertySystem[delId].hasInheritance()
		delete(teb.PropertySystem, delId)
	case `oncall`:
		teb.actionPropertyDelete(
			teb.PropertyOncall[delId].MakeAction(),
		)
		hasInheritance = teb.PropertyOncall[delId].hasInheritance()
		delete(teb.PropertyOncall, delId)
	default:
		teb.Fault.Error <- &Error{Action: `bucket.rmProperty unknown type`}
		return false
	}
	return hasInheritance
}

//
// Propertier:> Utility

// used to verify this is a source instance
func (teb *Bucket) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := teb.PropertyCustom[id]; !ok {
			goto bailout
		}
		return teb.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := teb.PropertyService[id]; !ok {
			goto bailout
		}
		return teb.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := teb.PropertySystem[id]; !ok {
			goto bailout
		}
		return teb.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := teb.PropertyOncall[id]; !ok {
			goto bailout
		}
		return teb.PropertyOncall[id].GetSourceInstance() == id
	}

bailout:
	teb.Fault.Error <- &Error{
		Action: `bucket.verifySourceInstance not found`}
	return false
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
		teb.Fault.Error <- &Error{Action: `bucket.checkDuplicate unknown type`}
		dupe = true
		deleteOK = false
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
