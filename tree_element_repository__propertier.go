package somatree

import "sync"

//
// Interface: SomaTreePropertier
func (ter *SomaTreeElemRepository) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		p.(*SomaTreePropertyCustom).InheritedFrom = ter.Id
		p.(*SomaTreePropertyCustom).Inherited = false
		ter.setCustomProperty(p)
		f := new(SomaTreePropertyCustom)
		*f = *p.(*SomaTreePropertyCustom)
		f.Inherited = true
		ter.inheritPropertyDeep(f)
	case "service":
		p.(*SomaTreePropertyService).InheritedFrom = ter.Id
		p.(*SomaTreePropertyService).Inherited = false
		ter.setServiceProperty(p)
		f := new(SomaTreePropertyService)
		*f = *p.(*SomaTreePropertyService)
		f.Inherited = true
		ter.inheritPropertyDeep(f)
	case "system":
		p.(*SomaTreePropertySystem).InheritedFrom = ter.Id
		p.(*SomaTreePropertySystem).Inherited = false
		ter.setSystemProperty(p)
		f := new(SomaTreePropertySystem)
		*f = *p.(*SomaTreePropertySystem)
		f.Inherited = true
		ter.inheritPropertyDeep(f)
	case "oncall":
		p.(*SomaTreePropertyOncall).InheritedFrom = ter.Id
		p.(*SomaTreePropertyOncall).Inherited = false
		ter.setOncallProperty(p)
		f := new(SomaTreePropertyOncall)
		*f = *p.(*SomaTreePropertyOncall)
		f.Inherited = true
		ter.inheritPropertyDeep(f)
	}
	ter.Action <- &Action{
		Action:         "property_new",
		Type:           "repository",
		Id:             ter.Id.String(),
		Name:           ter.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}
}

func (ter *SomaTreeElemRepository) inheritProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		ter.setCustomProperty(p)
	case "service":
		ter.setServiceProperty(p)
	case "system":
		ter.setSystemProperty(p)
	case "oncall":
		ter.setOncallProperty(p)
	}
	ter.Action <- &Action{
		Action:         "property_new",
		Type:           "repository",
		Id:             ter.Id.String(),
		Name:           ter.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}
	ter.inheritPropertyDeep(p)
}

func (ter *SomaTreeElemRepository) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func(stp SomaTreeProperty) {
			defer wg.Done()
			ter.Children[c].inheritProperty(stp)
		}(p)
	}
	wg.Wait()
}

func (ter *SomaTreeElemRepository) setCustomProperty(
	p SomaTreeProperty) {
	ter.PropertyCustom[p.GetID()] = p
}

func (ter *SomaTreeElemRepository) setServiceProperty(
	p SomaTreeProperty) {
	ter.PropertyService[p.GetID()] = p
}

func (ter *SomaTreeElemRepository) setSystemProperty(
	p SomaTreeProperty) {
	ter.PropertySystem[p.GetID()] = p
}

func (ter *SomaTreeElemRepository) setOncallProperty(
	p SomaTreeProperty) {
	ter.PropertyOncall[p.GetID()] = p
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (ter *SomaTreeElemRepository) syncProperty(
	childId string) {
customloop:
	for prop, _ := range ter.PropertyCustom {
		if !ter.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := new(SomaTreePropertyCustom)
		*f = *ter.PropertyCustom[prop].(*SomaTreePropertyCustom)
		f.Inherited = true
		ter.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range ter.PropertyOncall {
		if !ter.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(SomaTreePropertyOncall)
		*f = *ter.PropertyOncall[prop].(*SomaTreePropertyOncall)
		f.Inherited = true
		ter.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range ter.PropertyService {
		if !ter.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(SomaTreePropertyService)
		*f = *ter.PropertyService[prop].(*SomaTreePropertyService)
		f.Inherited = true
		ter.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range ter.PropertySystem {
		if !ter.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(SomaTreePropertySystem)
		*f = *ter.PropertySystem[prop].(*SomaTreePropertySystem)
		f.Inherited = true
		ter.Children[childId].inheritProperty(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (ter *SomaTreeElemRepository) checkProperty(
	propType string, propId string) bool {
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
