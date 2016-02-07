package somatree

import "github.com/satori/go.uuid"

type SomaTreeProperty interface {
	GetType() string
	GetID() string
	GetSource() string
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
	value         string
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
