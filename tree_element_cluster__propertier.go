package somatree

import "sync"

//
// Interface: SomaTreePropertier
func (tec *SomaTreeElemCluster) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		p.(*SomaTreePropertyCustom).InheritedFrom = tec.Id
		p.(*SomaTreePropertyCustom).Inherited = false
		tec.setCustomProperty(p)
		f := new(SomaTreePropertyCustom)
		*f = *p.(*SomaTreePropertyCustom)
		f.Inherited = true
		tec.inheritPropertyDeep(f)
	case "service":
		p.(*SomaTreePropertyService).InheritedFrom = tec.Id
		p.(*SomaTreePropertyService).Inherited = false
		tec.setServiceProperty(p)
		f := new(SomaTreePropertyService)
		*f = *p.(*SomaTreePropertyService)
		f.Inherited = true
		tec.inheritPropertyDeep(f)
	case "system":
		p.(*SomaTreePropertySystem).InheritedFrom = tec.Id
		p.(*SomaTreePropertySystem).Inherited = false
		tec.setSystemProperty(p)
		f := new(SomaTreePropertySystem)
		*f = *p.(*SomaTreePropertySystem)
		f.Inherited = true
		tec.inheritPropertyDeep(f)
	case "oncall":
		p.(*SomaTreePropertyOncall).InheritedFrom = tec.Id
		p.(*SomaTreePropertyOncall).Inherited = false
		tec.setOncallProperty(p)
		f := new(SomaTreePropertyOncall)
		*f = *p.(*SomaTreePropertyOncall)
		f.Inherited = true
		tec.inheritPropertyDeep(f)
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
		f := new(SomaTreePropertyCustom)
		*f = *tec.PropertyCustom[prop].(*SomaTreePropertyCustom)
		f.Inherited = true
		tec.Children[childId].inheritProperty(f)
	}
oncallloop:
	for prop, _ := range tec.PropertyOncall {
		if !tec.PropertyOncall[prop].hasInheritance() {
			continue oncallloop
		}
		f := new(SomaTreePropertyOncall)
		*f = *tec.PropertyOncall[prop].(*SomaTreePropertyOncall)
		f.Inherited = true
		tec.Children[childId].inheritProperty(f)
	}
serviceloop:
	for prop, _ := range tec.PropertyService {
		if !tec.PropertyService[prop].hasInheritance() {
			continue serviceloop
		}
		f := new(SomaTreePropertyService)
		*f = *tec.PropertyService[prop].(*SomaTreePropertyService)
		f.Inherited = true
		tec.Children[childId].inheritProperty(f)
	}
systemloop:
	for prop, _ := range tec.PropertySystem {
		if !tec.PropertySystem[prop].hasInheritance() {
			continue systemloop
		}
		f := new(SomaTreePropertySystem)
		*f = *tec.PropertySystem[prop].(*SomaTreePropertySystem)
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
