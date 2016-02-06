package somatree

import "sync"

//
// Interface: SomaTreePropertier
func (teb *SomaTreeElemBucket) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		p.(*SomaTreePropertyCustom).InheritedFrom = teb.Id
		p.(*SomaTreePropertyCustom).Inherited = false
		teb.setCustomProperty(p)
		f := new(SomaTreePropertyCustom)
		*f = *p.(*SomaTreePropertyCustom)
		f.Inherited = true
		teb.inheritPropertyDeep(f)
	case "service":
		p.(*SomaTreePropertyService).InheritedFrom = teb.Id
		p.(*SomaTreePropertyService).Inherited = false
		teb.setServiceProperty(p)
		f := new(SomaTreePropertyService)
		*f = *p.(*SomaTreePropertyService)
		f.Inherited = true
		teb.inheritPropertyDeep(f)
	case "system":
		p.(*SomaTreePropertySystem).InheritedFrom = teb.Id
		p.(*SomaTreePropertySystem).Inherited = false
		teb.setSystemProperty(p)
		f := new(SomaTreePropertySystem)
		*f = *p.(*SomaTreePropertySystem)
		f.Inherited = true
		teb.inheritPropertyDeep(f)
	case "oncall":
		p.(*SomaTreePropertyOncall).InheritedFrom = teb.Id
		p.(*SomaTreePropertyOncall).Inherited = false
		teb.setOncallProperty(p)
		f := new(SomaTreePropertyOncall)
		*f = *p.(*SomaTreePropertyOncall)
		f.Inherited = true
		teb.inheritPropertyDeep(f)
	}
	teb.Action <- &Action{
		Action:         "property_new",
		Type:           "bucket",
		Id:             teb.Id.String(),
		Name:           teb.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}
}

func (teb *SomaTreeElemBucket) inheritProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		teb.setCustomProperty(p)
	case "service":
		teb.setServiceProperty(p)
	case "system":
		teb.setSystemProperty(p)
	case "oncall":
		teb.setOncallProperty(p)
	}
	teb.Action <- &Action{
		Action:         "property_new",
		Type:           "bucket",
		Id:             teb.Id.String(),
		Name:           teb.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}

	teb.inheritPropertyDeep(p)
}

func (teb *SomaTreeElemBucket) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func(stp SomaTreeProperty) {
			defer wg.Done()
			teb.Children[c].inheritProperty(stp)
		}(p)
	}
	wg.Wait()
}

func (teb *SomaTreeElemBucket) setCustomProperty(
	p SomaTreeProperty) {
	teb.PropertyCustom[p.GetID()] = p
}

func (teb *SomaTreeElemBucket) setServiceProperty(
	p SomaTreeProperty) {
	teb.PropertyService[p.GetID()] = p
}

func (teb *SomaTreeElemBucket) setSystemProperty(
	p SomaTreeProperty) {
	teb.PropertySystem[p.GetID()] = p
}

func (teb *SomaTreeElemBucket) setOncallProperty(
	p SomaTreeProperty) {
	teb.PropertyOncall[p.GetID()] = p
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (teb *SomaTreeElemBucket) syncProperty(
	childId string) {
customloop:
	for prop, _ := range teb.PropertyCustom {
		if !teb.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := new(SomaTreePropertyCustom)
		*f = *teb.PropertyCustom[prop].(*SomaTreePropertyCustom)
		f.Inherited = true
		teb.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range teb.PropertyOncall {
		if !teb.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(SomaTreePropertyOncall)
		*f = *teb.PropertyOncall[prop].(*SomaTreePropertyOncall)
		f.Inherited = true
		teb.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range teb.PropertyService {
		if !teb.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(SomaTreePropertyService)
		*f = *teb.PropertyService[prop].(*SomaTreePropertyService)
		f.Inherited = true
		teb.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range teb.PropertySystem {
		if !teb.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(SomaTreePropertySystem)
		*f = *teb.PropertySystem[prop].(*SomaTreePropertySystem)
		f.Inherited = true
		teb.Children[childId].inheritProperty(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (teb *SomaTreeElemBucket) checkProperty(
	propType string, propId string) bool {
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
