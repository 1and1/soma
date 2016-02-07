package somatree

import "sync"

//
// Interface: SomaTreePropertier
func (teg *SomaTreeElemGroup) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		p.(*PropertyCustom).InheritedFrom = teg.Id
		p.(*PropertyCustom).Inherited = false
		teg.setCustomProperty(p)
		f := new(PropertyCustom)
		*f = *p.(*PropertyCustom)
		f.Inherited = true
		teg.inheritPropertyDeep(f)
	case "service":
		p.(*PropertyService).InheritedFrom = teg.Id
		p.(*PropertyService).Inherited = false
		teg.setServiceProperty(p)
		f := new(PropertyService)
		*f = *p.(*PropertyService)
		f.Inherited = true
		teg.inheritPropertyDeep(f)
	case "system":
		p.(*PropertySystem).InheritedFrom = teg.Id
		p.(*PropertySystem).Inherited = false
		teg.setSystemProperty(p)
		f := new(PropertySystem)
		*f = *p.(*PropertySystem)
		f.Inherited = true
		teg.inheritPropertyDeep(f)
	case "oncall":
		p.(*PropertyOncall).InheritedFrom = teg.Id
		p.(*PropertyOncall).Inherited = false
		teg.setOncallProperty(p)
		f := new(PropertyOncall)
		*f = *p.(*PropertyOncall)
		f.Inherited = true
		teg.inheritPropertyDeep(f)
	}
	teg.Action <- &Action{
		Action:         "property_new",
		Type:           "group",
		Id:             teg.Id.String(),
		Name:           teg.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}
}

func (teg *SomaTreeElemGroup) inheritProperty(
	p SomaTreeProperty) {
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
	teg.Action <- &Action{
		Action:         "property_new",
		Type:           "group",
		Id:             teg.Id.String(),
		Name:           teg.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}

	teg.inheritPropertyDeep(p)
}

func (teg *SomaTreeElemGroup) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
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
		f := new(PropertyCustom)
		*f = *teg.PropertyCustom[prop].(*PropertyCustom)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range teg.PropertyOncall {
		if !teg.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(PropertyOncall)
		*f = *teg.PropertyOncall[prop].(*PropertyOncall)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range teg.PropertyService {
		if !teg.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(PropertyService)
		*f = *teg.PropertyService[prop].(*PropertyService)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range teg.PropertySystem {
		if !teg.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(PropertySystem)
		*f = *teg.PropertySystem[prop].(*PropertySystem)
		f.Inherited = true
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
