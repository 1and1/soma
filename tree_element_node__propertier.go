package somatree

//
// Interface: SomaTreePropertier
func (ten *SomaTreeElemNode) SetProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		p.(*SomaTreePropertyCustom).InheritedFrom = ten.Id
		p.(*SomaTreePropertyCustom).Inherited = false
		ten.setCustomProperty(p)
		f := new(SomaTreePropertyCustom)
		*f = *p.(*SomaTreePropertyCustom)
		f.Inherited = true
		ten.inheritPropertyDeep(f)
	case "service":
		p.(*SomaTreePropertyService).InheritedFrom = ten.Id
		p.(*SomaTreePropertyService).Inherited = false
		ten.setServiceProperty(p)
		f := new(SomaTreePropertyService)
		*f = *p.(*SomaTreePropertyService)
		f.Inherited = true
		ten.inheritPropertyDeep(f)
	case "system":
		p.(*SomaTreePropertySystem).InheritedFrom = ten.Id
		p.(*SomaTreePropertySystem).Inherited = false
		ten.setSystemProperty(p)
		f := new(SomaTreePropertySystem)
		*f = *p.(*SomaTreePropertySystem)
		f.Inherited = true
		ten.inheritPropertyDeep(f)
	case "oncall":
		p.(*SomaTreePropertyOncall).InheritedFrom = ten.Id
		p.(*SomaTreePropertyOncall).Inherited = false
		ten.setOncallProperty(p)
		f := new(SomaTreePropertyOncall)
		*f = *p.(*SomaTreePropertyOncall)
		f.Inherited = true
		ten.inheritPropertyDeep(f)
	}
	ten.Action <- &Action{
		Action:         "property_new",
		Type:           "node",
		Id:             ten.Id.String(),
		Name:           ten.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}
}

func (ten *SomaTreeElemNode) inheritProperty(
	p SomaTreeProperty) {
	switch p.GetType() {
	case "custom":
		ten.setCustomProperty(p)
	case "service":
		ten.setServiceProperty(p)
	case "system":
		ten.setSystemProperty(p)
	case "oncall":
		ten.setOncallProperty(p)
	}
	ten.Action <- &Action{
		Action:         "property_new",
		Type:           "node",
		Id:             ten.Id.String(),
		Name:           ten.Name,
		PropertyType:   p.GetType(),
		PropertyId:     p.GetID(),
		PropertySource: p.GetSource(),
	}
	// no inheritPropertyDeep(), nodes have no children
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) inheritPropertyDeep(
	p SomaTreeProperty) {
}

func (ten *SomaTreeElemNode) setCustomProperty(
	p SomaTreeProperty) {
	ten.PropertyCustom[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setServiceProperty(
	p SomaTreeProperty) {
	ten.PropertyService[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setSystemProperty(
	p SomaTreeProperty) {
	ten.PropertySystem[p.GetID()] = p
}

func (ten *SomaTreeElemNode) setOncallProperty(
	p SomaTreeProperty) {
	ten.PropertyOncall[p.GetID()] = p
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) syncProperty(
	childId string) {
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) checkProperty(
	propType string, propId string) bool {
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
