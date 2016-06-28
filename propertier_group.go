package tree

import (
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
			cstUUID, _ := uuid.FromString(prop.GetKey())
			teg.deletePropertyInherited(&PropertyCustom{
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
			teg.deletePropertyInherited(&PropertyService{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Service:   prop.GetKey(),
			})
		case `system`:
			teg.deletePropertyInherited(&PropertySystem{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				Key:       prop.GetKey(),
				Value:     prop.GetValue(),
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(prop.GetKey())
			teg.deletePropertyInherited(&PropertyOncall{
				SourceId:  srcUUID,
				View:      prop.GetView(),
				Inherited: true,
				OncallId:  oncUUID,
				Name:      prop.(*PropertyOncall).GetName(),
				Number:    prop.(*PropertyOncall).GetNumber(),
			})
		}
	}
	p.SetId(p.GetInstanceId(teg.Type, teg.Id))
	if p.Equal(uuid.Nil) {
		p.SetId(uuid.NewV4())
	}
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
	if f.hasInheritance() {
		teg.setPropertyOnChildren(f)
	}
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
	}
	f.clearInstances()

	if !f.GetIsInherited() {
		teg.Fault.Error <- &Error{
			Action: `group.setPropertyInherited on inherited=false`}
		return
	}
	if dupe, deleteOK, _ := teg.checkDuplicate(p); dupe && deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate, but we are not the
		// source of the duplicate -> corrupt tree
		teg.Fault.Error <- &Error{
			Action: `group.setPropertyInherited corruption detected`}
		return
	} else if dupe && !deleteOK {
		// we received an inherited SetProperty from above us in the
		// tree for a property that is duplicate; we have a locally
		// set property -> stop inheritance, no error
		return
	}
	teg.addProperty(f)
	p.SetId(uuid.UUID{})
	teg.setPropertyOnChildren(p)
	teg.actionPropertyNew(f.MakeAction())
}

func (teg *Group) setPropertyOnChildren(p Property) {
	var wg sync.WaitGroup
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
	default:
		teg.Fault.Error <- &Error{Action: `group.addProperty unknown type`}
	}
}

//
// Propertier:> Update Property

func (teg *Group) UpdateProperty(p Property) {
	if !teg.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		teg.Fault.Error <- &Error{Action: `update_property_on_non_source`}
		return
	}

	// keep a copy for ourselves, no shared pointers
	p.SetInheritedFrom(teg.Id)
	p.SetSourceType(teg.Type)
	p.SetInherited(true)
	f := p.Clone()
	f.SetInherited(false)
	if teg.switchProperty(f) {
		teg.updatePropertyOnChildren(p)
	}
}

func (teg *Group) updatePropertyInherited(p Property) {
	// keep a copy for ourselves, no shared pointers
	f := p.Clone()
	if !f.GetIsInherited() {
		teg.Fault.Error <- &Error{
			Action: `group.updatePropertyInherited on inherited=false`}
		return
	}
	if teg.switchProperty(f) {
		teg.updatePropertyOnChildren(p)
	}
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

func (teg *Group) switchProperty(p Property) bool {
	uid := teg.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if uid == `` {
		// we do not have the property for which we received an update
		if dupe, deleteOK, _ := teg.checkDuplicate(p); dupe && !deleteOK {
			// the update is duplicate to an property for which we
			// have the source instance, ie we just received an update
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}
		teg.Fault.Error <- &Error{
			Action: `group.switchProperty property not found`}
		return false
	}
	updId, _ := uuid.FromString(uid)
	p.SetId(updId)
	curr := teg.getCurrentProperty(p)
	if curr == nil {
		return false
	}
	teg.addProperty(p)
	teg.actionPropertyUpdate(p.MakeAction())

	if !p.hasInheritance() && curr.hasInheritance() {
		// replacing inheritance with !inheritance:
		// call deletePropertyOnChildren(curr) to clean up
		srcUUID, _ := uuid.FromString(curr.GetSourceInstance())
		switch curr.GetType() {
		case `custom`:
			cstUUID, _ := uuid.FromString(curr.GetKey())
			teg.deletePropertyOnChildren(&PropertyCustom{
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
			teg.deletePropertyOnChildren(&PropertyService{
				SourceId:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Service:     curr.GetKey(),
				Inheritance: true,
			})
		case `system`:
			teg.deletePropertyOnChildren(&PropertySystem{
				SourceId:    srcUUID,
				View:        curr.GetView(),
				Inherited:   true,
				Key:         curr.GetKey(),
				Value:       curr.GetValue(),
				Inheritance: true,
			})
		case `oncall`:
			oncUUID, _ := uuid.FromString(curr.GetKey())
			teg.deletePropertyOnChildren(&PropertyOncall{
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
		teg.setPropertyOnChildren(t)
	}
	return p.hasInheritance() && curr.hasInheritance()
}

func (teg *Group) getCurrentProperty(p Property) Property {
	switch p.GetType() {
	case `custom`:
		return teg.PropertyCustom[p.GetID()].Clone()
	case `system`:
		return teg.PropertySystem[p.GetID()].Clone()
	case `service`:
		return teg.PropertyService[p.GetID()].Clone()
	case `oncall`:
		return teg.PropertyOncall[p.GetID()].Clone()
	}
	teg.Fault.Error <- &Error{
		Action: `group.getCurrentProperty unknown type`}
	return nil
}

//
// Propertier:> Delete Property

func (teg *Group) DeleteProperty(p Property) {
	if !teg.verifySourceInstance(
		p.GetSourceInstance(),
		p.GetType(),
	) {
		teg.Fault.Error <- &Error{Action: `group.DeleteProperty on !source`}
		return
	}

	var flow Property
	resync := false
	delId := teg.findIdForSource(
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
		resync, _, flow = teg.Parent.(Propertier).checkDuplicate(p)
	}

	p.SetInherited(false)
	if teg.rmProperty(p) {
		p.SetInherited(true)
		teg.deletePropertyOnChildren(p)
	}

	// now that the property is deleted from us and our children,
	// request resync if required
	if resync {
		teg.Parent.resyncProperty(flow.GetSourceInstance(),
			p.GetType(),
			teg.Id.String(),
		)
	}
}

func (teg *Group) deletePropertyInherited(p Property) {
	if teg.rmProperty(p) {
		teg.deletePropertyOnChildren(p)
	}
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

func (teg *Group) deletePropertyAllInherited() {
	for _, p := range teg.PropertyCustom {
		if !p.GetIsInherited() {
			continue
		}
		teg.deletePropertyInherited(p.Clone())
	}
	for _, p := range teg.PropertySystem {
		if !p.GetIsInherited() {
			continue
		}
		teg.deletePropertyInherited(p.Clone())
	}
	for _, p := range teg.PropertyService {
		if !p.GetIsInherited() {
			continue
		}
		teg.deletePropertyInherited(p.Clone())
	}
	for _, p := range teg.PropertyOncall {
		if !p.GetIsInherited() {
			continue
		}
		teg.deletePropertyInherited(p.Clone())
	}
}

func (teg *Group) deletePropertyAllLocal() {
	for _, p := range teg.PropertyCustom {
		if p.GetIsInherited() {
			continue
		}
		teg.DeleteProperty(p.Clone())
	}
	for _, p := range teg.PropertySystem {
		if p.GetIsInherited() {
			continue
		}
		teg.DeleteProperty(p.Clone())
	}
	for _, p := range teg.PropertyService {
		if p.GetIsInherited() {
			continue
		}
		teg.DeleteProperty(p.Clone())
	}
	for _, p := range teg.PropertyOncall {
		if p.GetIsInherited() {
			continue
		}
		teg.DeleteProperty(p.Clone())
	}
}

func (teg *Group) rmProperty(p Property) bool {
	delId := teg.findIdForSource(
		p.GetSourceInstance(),
		p.GetType(),
	)
	if delId == `` {
		// we do not have the property for which we received a delete
		if dupe, deleteOK, _ := teg.checkDuplicate(p); dupe && !deleteOK {
			// the delete is duplicate to a property for which we
			// have the source instance, ie we just received a delete
			// for which we have an overwrite. Ignore it and do not
			// inherit it further down
			return false
		}

		teg.Fault.Error <- &Error{
			Action: `group.rmProperty property not found`}
		return false
	}

	hasInheritance := false
	switch p.GetType() {
	case `custom`:
		teg.actionPropertyDelete(
			teg.PropertyCustom[delId].MakeAction(),
		)
		hasInheritance = teg.PropertyCustom[delId].hasInheritance()
		delete(teg.PropertyCustom, delId)
	case `service`:
		teg.actionPropertyDelete(
			teg.PropertyService[delId].MakeAction(),
		)
		hasInheritance = teg.PropertyService[delId].hasInheritance()
		delete(teg.PropertyService, delId)
	case `system`:
		teg.actionPropertyDelete(
			teg.PropertySystem[delId].MakeAction(),
		)
		hasInheritance = teg.PropertySystem[delId].hasInheritance()
		delete(teg.PropertySystem, delId)
	case `oncall`:
		teg.actionPropertyDelete(
			teg.PropertyOncall[delId].MakeAction(),
		)
		hasInheritance = teg.PropertyOncall[delId].hasInheritance()
		delete(teg.PropertyOncall, delId)
	default:
		teg.Fault.Error <- &Error{Action: `group.rmProperty unknown type`}
		return false
	}
	return hasInheritance
}

//
// Propertier:> Utility

//
func (teg *Group) verifySourceInstance(id, prop string) bool {
	switch prop {
	case `custom`:
		if _, ok := teg.PropertyCustom[id]; !ok {
			goto bailout
		}
		return teg.PropertyCustom[id].GetSourceInstance() == id
	case `service`:
		if _, ok := teg.PropertyService[id]; !ok {
			goto bailout
		}
		return teg.PropertyService[id].GetSourceInstance() == id
	case `system`:
		if _, ok := teg.PropertySystem[id]; !ok {
			goto bailout
		}
		return teg.PropertySystem[id].GetSourceInstance() == id
	case `oncall`:
		if _, ok := teg.PropertyOncall[id]; !ok {
			goto bailout
		}
		return teg.PropertyOncall[id].GetSourceInstance() == id
	}

bailout:
	teg.Fault.Error <- &Error{
		Action: `group.verifySourceInstance not found`}
	return false
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

//
func (teg *Group) resyncProperty(srcId, pType, childId string) {
	pId := teg.findIdForSource(srcId, pType)
	if pId == `` {
		return
	}

	var f Property
	switch pType {
	case `custom`:
		f = teg.PropertyCustom[pId].(*PropertyCustom).Clone()
	case `oncall`:
		f = teg.PropertyOncall[pId].(*PropertyOncall).Clone()
	case `service`:
		f = teg.PropertyService[pId].(*PropertyService).Clone()
	case `system`:
		f = teg.PropertySystem[pId].(*PropertySystem).Clone()
	}
	if !f.hasInheritance() {
		return
	}
	f.SetInherited(true)
	f.SetId(uuid.UUID{})
	f.clearInstances()
	teg.Children[childId].setPropertyInherited(f)
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
		teg.Fault.Error <- &Error{Action: `group.checkDuplicate unknown type`}
		dupe = true
		deleteOK = false
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
