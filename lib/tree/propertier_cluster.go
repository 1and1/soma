/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"sync"

	"github.com/satori/go.uuid"
)

// Implementation of the `Propertier` interface

//
// Propertier:> Add Property

func (tec *Cluster) SetProperty(p Property) {
	// if deleteOK is true, then prop is the property that can be
	// deleted
	if dupe, deleteOK, prop := tec.checkDuplicate(p); dupe && !deleteOK {
		tec.Fault.Error <- &Error{Action: `duplicate_set_property`}
		return
	} else if dupe && deleteOK {
		srcUUID, _ := uuid.FromString(prop.GetSourceInstance())
		switch prop.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(prop.GetKey())
			tec.deletePropertyInherited(&PropertyCustom{
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
			tec.deletePropertyInherited(&PropertyService{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Service:   prop.GetKey(),
			})
		case `system`:
			tec.deletePropertyInherited(&PropertySystem{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Key:       prop.GetKey(),
				Value:     prop.GetValue(),
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(prop.GetKey())
			tec.deletePropertyInherited(&PropertyOncall{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				OncallId:  oncUUID,
				Name:      prop.(*PropertyOncall).GetName(),
				Number:    prop.(*PropertyOncall).GetNumber(),
			})
		}
	}
	p.SetId(p.GetInstanceId(tec.Type, tec.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
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
	if f.hasInheritance() {
		tec.setPropertyOnChildren(f)
	}
	// scrub instance startup information prior to storing
	p.clearInstances()
	tec.addProperty(p)
	tec.actionPropertyNew(p.MakeAction())
}

func (tec *Cluster) setPropertyInherited(p Property) {
	f := p.Clone()
	f.SetId(f.GetInstanceId(tec.Type, tec.Id))
	if f.Equal(uuid.Nil) {
		f.SetId(uuid.NewV4())
	}
	f.clearInstances()

	if !f.GetIsInherited() {
		tec.Fault.Error <- &Error{
			Action: `cluster.setPropertyInherited on inherited=false`}
		return
	}
	if dupe, deleteOK, _ := tec.checkDuplicate(p); dupe && deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate, but we are not the
		// source of the duplicate -> corrupt tree
		tec.Fault.Error <- &Error{
			Action: `cluster.setPropertyInherited corruption detected`}
		return
	} else if dupe && !deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate; we have a locally
		// set property -> stop inheritance, no error
		return
	}
	tec.addProperty(f)
	p.SetId(uuid.UUID{})
	tec.setPropertyOnChildren(p)
	tec.actionPropertyNew(f.MakeAction())
}

func (tec *Cluster) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		go func(stp Property, c string) {
			defer wg.Done()
			tec.Children[c].setPropertyInherited(stp)
		}(p, child)
	}
	wg.Wait()
}

func (tec *Cluster) addProperty(p Property) {
	switch p.GetType() {
	case `custom`:
		tec.PropertyCustom[p.GetID()] = p
	case `system`:
		tec.PropertySystem[p.GetID()] = p
	case `service`:
		tec.PropertyService[p.GetID()] = p
	case `oncall`:
		tec.PropertyOncall[p.GetID()] = p
	default:
		tec.Fault.Error <- &Error{Action: `cluster.addProperty unknown type`}
	}
}

//
// Propertier:> Update Property

func (tec *Cluster) UpdateProperty(p Property) {
	if !tec.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		tec.Fault.Error <- &Error{Action: `update_property_on_non_source`}
		return
	}

	// keep a copy for ourselves, no shared pointers
	p.SetInheritedFrom(tec.Id)
	p.SetSourceType(tec.Type)
	p.SetInherited(true)
	f := p.Clone()
	f.SetInherited(false)
	if tec.switchProperty(f) {
		tec.updatePropertyOnChildren(p)
	}
}

func (tec *Cluster) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	if !f.GetIsInherited() {
		tec.Fault.Error <- &Error{
			Action: `cluster.updatePropertyInherited on inherited=false`}
		return
	}
	if tec.switchProperty(f) {
		tec.updatePropertyOnChildren(p)
	}
}

func (tec *Cluster) updatePropertyOnChildren(p Property) {
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

func (tec *Cluster) switchProperty(p Property) bool {
	uid := tec.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if uid == `` {
		// we do not have the property for which we received an update
		if dupe, deleteOK, _ := tec.checkDuplicate(p); dupe && !deleteOK {
			// the update is duplicate to an property for which we
			// have the source instance, ie we just received an update
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}
		tec.Fault.Error <- &Error{
			Action: `cluster.switchProperty property not found`}
		return false
	}
	updId, _ := uuid.FromString(uid)
	p.SetId(updId)
	curr := tec.getCurrentProperty(p)
	if curr == nil {
		return false
	}
	tec.addProperty(p)
	tec.actionPropertyUpdate(p.MakeAction())

	if !p.hasInheritance() && curr.hasInheritance() {
		// replacing inheritance with !inheritance:
		// call deletePropertyOnChildren(curr) to clean up
		srcUUID, _ := uuid.FromString(curr.GetSourceInstance())
		switch curr.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(curr.GetKey())
			tec.deletePropertyOnChildren(&PropertyCustom{
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
			tec.deletePropertyOnChildren(&PropertyService{
				SourceId:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Service:     curr.GetKey(),
				Inheritance: true,
			})
		case `system`:
			tec.deletePropertyOnChildren(&PropertySystem{
				SourceId:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Key:         curr.GetKey(),
				Value:       curr.GetValue(),
				Inheritance: true,
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(curr.GetKey())
			tec.deletePropertyOnChildren(&PropertyOncall{
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
		tec.setPropertyOnChildren(t)
	}
	return p.hasInheritance() && curr.hasInheritance()
}

func (tec *Cluster) getCurrentProperty(p Property) Property {
	switch p.GetType() {
	case `custom`:
		return tec.PropertyCustom[p.GetID()].Clone()
	case `system`:
		return tec.PropertySystem[p.GetID()].Clone()
	case `service`:
		return tec.PropertyService[p.GetID()].Clone()
	case `oncall`:
		return tec.PropertyOncall[p.GetID()].Clone()
	}
	tec.Fault.Error <- &Error{
		Action: `cluster.getCurrentProperty unknown type`}
	return nil
}

//
// Propertier:> Delete Property

func (tec *Cluster) DeleteProperty(p Property) {
	if !tec.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		tec.Fault.Error <- &Error{Action: `cluster.DeleteProperty on !source`}
		return
	}

	var flow Property
	resync := false
	delId := tec.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delId != `` {
		// this is a delete for a locally set property. It might be a
		// delete for an overwrite property, in which case we need to
		// ask the parent to sync it to us again.
		// If it was an overwrite, the parent should have a property
		// we would consider a dupe if it were to be passed down to
		// us.
		// If p is considered a dupe, then flow is set to the prop we
		// need to inherit.
		var delProp Property
		switch p.GetType() {
		case `custom`:
			delProp = tec.PropertyCustom[delId]
		case `system`:
			delProp = tec.PropertySystem[delId]
		case `service`:
			delProp = tec.PropertyService[delId]
		case `oncall`:
			delProp = tec.PropertyOncall[delId]
		}
		resync, _, flow = tec.Parent.(Propertier).checkDuplicate(
			delProp,
		)
	}

	p.SetInherited(false)
	if tec.rmProperty(p) {
		p.SetInherited(true)
		tec.deletePropertyOnChildren(p)
	}

	// now that the property is deleted from us and our children,
	// request resync if required
	if resync {
		tec.Parent.(Propertier).resyncProperty(
			flow.GetSourceInstance(),
			p.GetType(),
			tec.Id.String(),
		)
	}
}

func (tec *Cluster) deletePropertyInherited(p Property) {
	if tec.rmProperty(p) {
		tec.deletePropertyOnChildren(p)
	}
}

func (tec *Cluster) deletePropertyOnChildren(p Property) {
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

func (tec *Cluster) deletePropertyAllInherited() {
	for _, p := range tec.PropertyCustom {
		if !p.GetIsInherited() {
			continue
		}
		tec.deletePropertyInherited(p.Clone())
	}
	for _, p := range tec.PropertySystem {
		if !p.GetIsInherited() {
			continue
		}
		tec.deletePropertyInherited(p.Clone())
	}
	for _, p := range tec.PropertyService {
		if !p.GetIsInherited() {
			continue
		}
		tec.deletePropertyInherited(p.Clone())
	}
	for _, p := range tec.PropertyOncall {
		if !p.GetIsInherited() {
			continue
		}
		tec.deletePropertyInherited(p.Clone())
	}
}

func (tec *Cluster) deletePropertyAllLocal() {
	for _, p := range tec.PropertyCustom {
		if p.GetIsInherited() {
			continue
		}
		tec.DeleteProperty(p.Clone())
	}
	for _, p := range tec.PropertySystem {
		if p.GetIsInherited() {
			continue
		}
		tec.DeleteProperty(p.Clone())
	}
	for _, p := range tec.PropertyService {
		if p.GetIsInherited() {
			continue
		}
		tec.DeleteProperty(p.Clone())
	}
	for _, p := range tec.PropertyOncall {
		if p.GetIsInherited() {
			continue
		}
		tec.DeleteProperty(p.Clone())
	}
}

func (tec *Cluster) rmProperty(p Property) bool {
	delId := tec.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delId == `` {
		// we do not have the property for which we received a delete
		if dupe, deleteOK, _ := tec.checkDuplicate(p); dupe && !deleteOK {
			// the delete is duplicate to a property for which we
			// have the source instance, ie we just received a delete
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}

		tec.Fault.Error <- &Error{
			Action: `cluster.rmProperty property not found`}
		return false
	}

	hasInheritance := false
	switch p.GetType() {
	case `custom`:
		tec.actionPropertyDelete(
			tec.PropertyCustom[delId].MakeAction(),
		)
		hasInheritance = tec.PropertyCustom[delId].hasInheritance()
		delete(tec.PropertyCustom, delId)
	case `service`:
		tec.actionPropertyDelete(
			tec.PropertyService[delId].MakeAction(),
		)
		hasInheritance = tec.PropertyService[delId].hasInheritance()
		delete(tec.PropertyService, delId)
	case `system`:
		tec.actionPropertyDelete(
			tec.PropertySystem[delId].MakeAction(),
		)
		hasInheritance = tec.PropertySystem[delId].hasInheritance()
		delete(tec.PropertySystem, delId)
	case `oncall`:
		tec.actionPropertyDelete(
			tec.PropertyOncall[delId].MakeAction(),
		)
		hasInheritance = tec.PropertyOncall[delId].hasInheritance()
		delete(tec.PropertyOncall, delId)
	default:
		tec.Fault.Error <- &Error{Action: `cluster.rmProperty unknown type`}
		return false
	}
	return hasInheritance
}

//
// Propertier:> Utility

//
func (tec *Cluster) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := tec.PropertyCustom[id]; !ok {
			goto bailout
		}
		return tec.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := tec.PropertyService[id]; !ok {
			goto bailout
		}
		return tec.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := tec.PropertySystem[id]; !ok {
			goto bailout
		}
		return tec.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := tec.PropertyOncall[id]; !ok {
			goto bailout
		}
		return tec.PropertyOncall[id].GetSourceInstance() == id
	}

bailout:
	tec.Fault.Error <- &Error{
		Action: `cluster.verifySourceInstance not found`}
	return false
}

//
func (tec *Cluster) findIdForSource(source, prop string) string {
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

//
func (tec *Cluster) resyncProperty(srcId, pType, childId string) {
	pId := tec.findIdForSource(srcId, pType)
	if pId == `` {
		return
	}

	var f Property
	switch pType {
	case `custom`:
		f = tec.PropertyCustom[pId].(*PropertyCustom).Clone()
	case `oncall`:
		f = tec.PropertyOncall[pId].(*PropertyOncall).Clone()
	case `service`:
		f = tec.PropertyService[pId].(*PropertyService).Clone()
	case `system`:
		f = tec.PropertySystem[pId].(*PropertySystem).Clone()
	}
	if !f.hasInheritance() {
		return
	}
	f.SetInherited(true)
	f.SetId(uuid.UUID{})
	f.clearInstances()
	tec.Children[childId].setPropertyInherited(f)
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (tec *Cluster) syncProperty(childId string) {
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
func (tec *Cluster) checkProperty(propType string, propId string) bool {
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
func (tec *Cluster) checkDuplicate(p Property) (bool, bool, Property) {
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
		tec.Fault.Error <- &Error{Action: `cluster.checkDuplicate unknown type`}
		dupe = true
		deleteOK = false
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
