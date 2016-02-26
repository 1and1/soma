package somatree

import (
	"github.com/satori/go.uuid"
)

type SomaTreeProperty interface {
	GetType() string
	GetID() string
	GetSource() string
	GetInstanceId(objType string, objId uuid.UUID) uuid.UUID
	Clone() SomaTreeProperty
	hasInheritance() bool
	isChildrenOnly() bool
	GetSourceInstance() string
	GetSourceType() string
	GetIsInherited() bool
	GetView() string
}

type PropertyInstance struct {
	ObjectId   uuid.UUID
	ObjectType string
	InstanceId uuid.UUID
}

//
// Custom
type PropertyCustom struct {
	Id            uuid.UUID
	SourceId      uuid.UUID
	SourceType    string
	CustomId      uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Key           string
	Value         string
	Instances     []PropertyInstance
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

func (p *PropertyCustom) GetSourceInstance() string {
	return p.SourceId.String()
}

func (p *PropertyCustom) GetSourceType() string {
	return p.SourceType
}

func (p *PropertyCustom) GetIsInherited() bool {
	return p.Inherited
}

func (p *PropertyCustom) GetView() string {
	return p.View
}

func (p *PropertyCustom) GetInstanceId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(p.Id, uuid.Nil) {
		return p.Id
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectId, objId) {
			return instance.ObjectId
		}
	}
	return uuid.Nil
}

func (p PropertyCustom) Clone() SomaTreeProperty {
	cl := PropertyCustom{
		SourceType:   p.SourceType,
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Key:          p.Key,
		Value:        p.Value,
	}
	cl.Id, _ = uuid.FromString(p.Id.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())
	cl.SourceId, _ = uuid.FromString(p.SourceId.String())
	cl.CustomId, _ = uuid.FromString(p.CustomId.String())

	return &cl
}

//
// Service
type PropertyService struct {
	Id            uuid.UUID
	SourceId      uuid.UUID
	SourceType    string
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Service       string
	Attributes    []somaproto.TreeServiceAttribute
	Instances     []PropertyInstance
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

func (p *PropertyService) GetSourceInstance() string {
	return p.SourceId.String()
}

func (p *PropertyService) GetSourceType() string {
	return p.SourceType
}

func (p *PropertyService) GetIsInherited() bool {
	return p.Inherited
}

func (p *PropertyService) GetView() string {
	return p.View
}

func (p *PropertyService) GetInstanceId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(p.Id, uuid.Nil) {
		return p.Id
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectId, objId) {
			return instance.ObjectId
		}
	}
	return uuid.Nil
}

func (p PropertyService) Clone() SomaTreeProperty {
	cl := PropertyService{
		SourceType:   p.SourceType,
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Service:      p.Service,
	}
	cl.Id, _ = uuid.FromString(p.Id.String())
	cl.SourceId, _ = uuid.FromString(p.SourceId.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())
	cl.Attributes = make([]somaproto.TreeServiceAttribute, 0)
	for _, attr := range p.Attributes {
		a := somaproto.TreeServiceAttribute{
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
	SourceId      uuid.UUID
	SourceType    string
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Key           string
	Value         string
	Instances     []PropertyInstance
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

func (p *PropertySystem) GetSourceInstance() string {
	return p.SourceId.String()
}

func (p *PropertySystem) GetSourceType() string {
	return p.SourceType
}

func (p *PropertySystem) GetIsInherited() bool {
	return p.Inherited
}

func (p *PropertySystem) GetView() string {
	return p.View
}

func (p *PropertySystem) GetInstanceId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(p.Id, uuid.Nil) {
		return p.Id
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectId, objId) {
			return instance.ObjectId
		}
	}
	return uuid.Nil
}

func (p PropertySystem) Clone() SomaTreeProperty {
	cl := PropertySystem{
		SourceType:   p.SourceType,
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Key:          p.Key,
		Value:        p.Value,
	}
	cl.Id, _ = uuid.FromString(p.Id.String())
	cl.SourceId, _ = uuid.FromString(p.SourceId.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())

	return &cl
}

//
// Oncall
type PropertyOncall struct {
	Id            uuid.UUID
	SourceId      uuid.UUID
	SourceType    string
	OncallId      uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
	Inheritance   bool
	ChildrenOnly  bool
	View          string
	Name          string
	Number        string
	Instances     []PropertyInstance
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

func (p *PropertyOncall) GetSourceInstance() string {
	return p.SourceId.String()
}

func (p *PropertyOncall) GetSourceType() string {
	return p.SourceType
}

func (p *PropertyOncall) GetIsInherited() bool {
	return p.Inherited
}

func (p *PropertyOncall) GetView() string {
	return p.View
}

func (p *PropertyOncall) GetInstanceId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(p.Id, uuid.Nil) {
		return p.Id
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectId, objId) {
			return instance.ObjectId
		}
	}
	return uuid.Nil
}

func (p PropertyOncall) Clone() SomaTreeProperty {
	cl := PropertyOncall{
		SourceType:   p.SourceType,
		Inherited:    p.Inherited,
		Inheritance:  p.Inheritance,
		ChildrenOnly: p.ChildrenOnly,
		View:         p.View,
		Name:         p.Name,
		Number:       p.Number,
	}
	cl.Id, _ = uuid.FromString(p.Id.String())
	cl.SourceId, _ = uuid.FromString(p.SourceId.String())
	cl.OncallId, _ = uuid.FromString(p.OncallId.String())
	cl.InheritedFrom, _ = uuid.FromString(p.InheritedFrom.String())

	return &cl
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
