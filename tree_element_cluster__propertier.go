package somatree

import "sync"

//
// Interface: SomaTreePropertier
func (tec *SomaTreeElemCluster) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		p.(*PropertyCustom).InheritedFrom = tec.Id
		p.(*PropertyCustom).Inherited = false
		tec.setCustomProperty(p)
		f := new(PropertyCustom)
		*f = *p.(*PropertyCustom)
		f.Inherited = true
		tec.inheritPropertyDeep(f)
	case "service":
		p.(*PropertyService).InheritedFrom = tec.Id
		p.(*PropertyService).Inherited = false
		tec.setServiceProperty(p)
		f := new(PropertyService)
		*f = *p.(*PropertyService)
		f.Inherited = true
		tec.inheritPropertyDeep(f)
	case "system":
		p.(*PropertySystem).InheritedFrom = tec.Id
		p.(*PropertySystem).Inherited = false
		tec.setSystemProperty(p)
		f := new(PropertySystem)
		*f = *p.(*PropertySystem)
		f.Inherited = true
		tec.inheritPropertyDeep(f)
	case "oncall":
		p.(*PropertyOncall).InheritedFrom = tec.Id
		p.(*PropertyOncall).Inherited = false
		tec.setOncallProperty(p)
		f := new(PropertyOncall)
		*f = *p.(*PropertyOncall)
		f.Inherited = true
		tec.inheritPropertyDeep(f)
	}
	tec.Action <- &Action{
		Action:         "property_new",
		Type:           "cluster",
		Id:             tec.Id.String(),
		Name:           tec.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}
}

func (tec *SomaTreeElemCluster) inheritProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		tec.setCustomProperty(p)
	case "service":
		tec.setServiceProperty(p)
	case "system":
		tec.setSystemProperty(p)
	case "oncall":
		tec.setOncallProperty(p)
	}
	tec.Action <- &Action{
		Action:         "property_new",
		Type:           "cluster",
		Id:             tec.Id.String(),
		Name:           tec.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}
	tec.inheritPropertyDeep(p)
}

func (tec *SomaTreeElemCluster) inheritPropertyDeep(
	p SomaTreeProperty) {
	var wg sync.WaitGroup
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(stp SomaTreeProperty) {
			defer wg.Done()
			tec.Children[c].inheritProperty(stp)
		}(p)
	}
	wg.Wait()
}

func (tec *SomaTreeElemCluster) setCustomProperty(
	p SomaTreeProperty) {
	tec.PropertyCustom[p.GetID()] = p
}

func (tec *SomaTreeElemCluster) setServiceProperty(
	p SomaTreeProperty) {
	tec.PropertyService[p.GetID()] = p
}

func (tec *SomaTreeElemCluster) setSystemProperty(
	p SomaTreeProperty) {
	tec.PropertySystem[p.GetID()] = p
}

func (tec *SomaTreeElemCluster) setOncallProperty(
	p SomaTreeProperty) {
	tec.PropertyOncall[p.GetID()] = p
}

// when a child attaches, it calls self.Parent.syncProperty(self.Id)
// to get get all properties of that part of the tree
func (tec *SomaTreeElemCluster) syncProperty(
	childId string) {
customloop:
	for prop, _ := range tec.PropertyCustom {
		if !tec.PropertyCustom[prop].hasInheritance() {
			continue customloop
		}
		f := new(PropertyCustom)
		*f = *tec.PropertyCustom[prop].(*PropertyCustom)
		f.Inherited = true
		tec.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range tec.PropertyOncall {
		if !tec.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(PropertyOncall)
		*f = *tec.PropertyOncall[prop].(*PropertyOncall)
		f.Inherited = true
		tec.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range tec.PropertyService {
		if !tec.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(PropertyService)
		*f = *tec.PropertyService[prop].(*PropertyService)
		f.Inherited = true
		tec.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range tec.PropertySystem {
		if !tec.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(PropertySystem)
		*f = *tec.PropertySystem[prop].(*PropertySystem)
		f.Inherited = true
		tec.Children[childId].inheritProperty(f)
	}
}

// function to be used by a child to check if the parent has a
// specific Property
func (tec *SomaTreeElemCluster) checkProperty(
	propType string, propId string) bool {
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
