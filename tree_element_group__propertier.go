package somatree

import "sync"

//
// Interface: SomaTreePropertier
func (teg *SomaTreeElemGroup) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		p.(*SomaTreePropertyCustom).InheritedFrom = teg.Id
		p.(*SomaTreePropertyCustom).Inherited = false
		teg.setCustomProperty(p)
		f := new(SomaTreePropertyCustom)
		*f = *p.(*SomaTreePropertyCustom)
		f.Inherited = true
		teg.inheritPropertyDeep(f)
	case "service":
		p.(*SomaTreePropertyService).InheritedFrom = teg.Id
		p.(*SomaTreePropertyService).Inherited = false
		teg.setServiceProperty(p)
		f := new(SomaTreePropertyService)
		*f = *p.(*SomaTreePropertyService)
		f.Inherited = true
		teg.inheritPropertyDeep(f)
	case "system":
		p.(*SomaTreePropertySystem).InheritedFrom = teg.Id
		p.(*SomaTreePropertySystem).Inherited = false
		teg.setSystemProperty(p)
		f := new(SomaTreePropertySystem)
		*f = *p.(*SomaTreePropertySystem)
		f.Inherited = true
		teg.inheritPropertyDeep(f)
	case "oncall":
		p.(*SomaTreePropertyOncall).InheritedFrom = teg.Id
		p.(*SomaTreePropertyOncall).Inherited = false
		teg.setOncallProperty(p)
		f := new(SomaTreePropertyOncall)
		*f = *p.(*SomaTreePropertyOncall)
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
		f := new(SomaTreePropertyCustom)
		*f = *teg.PropertyCustom[prop].(*SomaTreePropertyCustom)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range teg.PropertyOncall {
		if !teg.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(SomaTreePropertyOncall)
		*f = *teg.PropertyOncall[prop].(*SomaTreePropertyOncall)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range teg.PropertyService {
		if !teg.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(SomaTreePropertyService)
		*f = *teg.PropertyService[prop].(*SomaTreePropertyService)
		f.Inherited = true
		teg.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range teg.PropertySystem {
		if !teg.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(SomaTreePropertySystem)
		*f = *teg.PropertySystem[prop].(*SomaTreePropertySystem)
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
