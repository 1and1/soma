package somatree

import "github.com/satori/go.uuid"

type SomaTreeProperty interface {
	GetType() string
	GetID() string
	GetSource() string
	Clone() SomaTreeProperty
	hasInheritance() bool
	isChildrenOnly() bool
}

//
// Custom
type PropertyCustom struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Key           string
	Value         string
}

func (p *PropertyCustom) GetType() string {
	return "custom"
}

func (p *PropertyCustom) GetID() string {
	return p.Id.String()
}

func (p *PropertyCustom) GetSource() string {
	return p.InheritedFrom.String()
}

func (p *PropertyCustom) hasInheritance() bool {
	return p.Inheritance
}

func (p *PropertyCustom) isChildrenOnly() bool {
	return p.ChildrenOnly
}

func (p PropertyCustom) Clone() SomaTreeProperty {
	cl := PropertyCustom{
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Key:          p.Key,
		Value:        p.Value,
	}
	cl.Id, _ = uuid.FromString(p.Id.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())

	return &cl
}

//
// Service
type PropertyService struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Service       string
	Attributes    []PropertyServiceAttribute
}

type PropertyServiceAttribute struct {
	Attribute string
	Value     string
}

func (p *PropertyService) GetType() string {
	return "service"
}

func (p *PropertyService) GetID() string {
	return p.Id.String()
}

func (p *PropertyService) GetSource() string {
	return p.InheritedFrom.String()
}

func (p *PropertyService) hasInheritance() bool {
	return p.Inheritance
}

func (p *PropertyService) isChildrenOnly() bool {
	return p.ChildrenOnly
}

func (p PropertyService) Clone() SomaTreeProperty {
	cl := PropertyService{
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Service:      p.Service,
	}
	cl.Id, _ = uuid.FromString(p.Id.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())
	cl.Attributes = make([]PropertyServiceAttribute, 0)
	for _, attr := range p.Attributes {
		a := PropertyServiceAttribute{
			Attribute: attr.Attribute,
			Value:     attr.Value,
		}
		cl.Attributes = append(cl.Attributes, a)
	}

	return &cl
}

//
// System
type PropertySystem struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Key           string
	Value         string
}

func (p *PropertySystem) GetType() string {
	return "system"
}

func (p *PropertySystem) GetID() string {
	return p.Id.String()
}

func (p *PropertySystem) GetSource() string {
	return p.InheritedFrom.String()
}

func (p *PropertySystem) hasInheritance() bool {
	return p.Inheritance
}

func (p *PropertySystem) isChildrenOnly() bool {
	return p.ChildrenOnly
}

func (p PropertySystem) Clone() SomaTreeProperty {
	cl := PropertySystem{
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Key:          p.Key,
		Value:        p.Value,
	}
	cl.Id, _ = uuid.FromString(p.Id.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())

	return &cl
}

//
// Oncall
type PropertyOncall struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Oncall        string
}

func (p *PropertyOncall) GetType() string {
	return "return oncall"
}

func (p *PropertyOncall) GetID() string {
	return p.Id.String()
}

func (p *PropertyOncall) GetSource() string {
	return p.InheritedFrom.String()
}

func (p *PropertyOncall) hasInheritance() bool {
	return p.Inheritance
}

func (p *PropertyOncall) isChildrenOnly() bool {
	return p.ChildrenOnly
}

func (p PropertyOncall) Clone() SomaTreeProperty {
	cl := PropertyOncall{
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Oncall:       p.Oncall,
	}
	cl.Id, _ = uuid.FromString(p.Id.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())

	return &cl
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
