/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"

	"github.com/1and1/soma/lib/proto"
)

type Property interface {
	GetID() string
	GetInstanceId(objType string, objId uuid.UUID) uuid.UUID
	GetIsInherited() bool
	GetKey() string
	GetSource() string
	GetSourceInstance() string
	GetSourceType() string
	GetType() string
	GetValue() string
	GetView() string

	SetId(id uuid.UUID)
	SetInherited(inherited bool)
	SetInheritedFrom(id uuid.UUID)
	SetSourceId(id uuid.UUID)
	SetSourceType(s string)

	Clone() Property
	Equal(id uuid.UUID) bool
	MakeAction() Action

	hasInheritance() bool
	isChildrenOnly() bool
	clearInstances()
}

type PropertyInstance struct {
	ObjectId   uuid.UUID
	ObjectType string
	InstanceId uuid.UUID
}

//
// Custom
type PropertyCustom struct {
	// Id of the custom property
	Id uuid.UUID
	// Id of the source custom property this was inherited from
	SourceId uuid.UUID
	// ObjectType the source property was attached to
	SourceType string
	// Id of the custom property type
	CustomId uuid.UUID
	// Indicator if this was inherited
	Inherited bool
	// Id of the object the SourceId property is on
	InheritedFrom uuid.UUID
	// Inheritance is enabled/disabled
	Inheritance bool
	// ChildrenOnly is enabled/disabled
	ChildrenOnly bool
	// View this property is attached in
	View string
	// Property Key
	Key string
	// Property Value
	Value string
	// Filled with IDs during from-DB load to restore with same IDs
	Instances []PropertyInstance
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

func (p *PropertyCustom) GetKey() string {
	return p.CustomId.String()
}

func (p *PropertyCustom) GetValue() string {
	return p.Key + p.Value
}

func (p *PropertyCustom) GetKeyField() string {
	return p.Key
}

func (p *PropertyCustom) GetValueField() string {
	return p.Value
}

func (p *PropertyCustom) GetInstanceId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(p.Id, uuid.Nil) {
		return p.Id
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectId, objId) {
			log.Printf("tree.Property.GetInstanceId() found existing instance: %s\n", instance.InstanceId)
			return instance.InstanceId
		}
	}
	return uuid.Nil
}

func (p *PropertyCustom) SetId(id uuid.UUID) {
	p.Id, _ = uuid.FromString(id.String())
}

func (p *PropertyCustom) Equal(id uuid.UUID) bool {
	return uuid.Equal(p.Id, id)
}

func (p *PropertyCustom) clearInstances() {
	p.Instances = nil
}

func (p *PropertyCustom) SetInheritedFrom(id uuid.UUID) {
	p.InheritedFrom, _ = uuid.FromString(id.String())
}

func (p *PropertyCustom) SetInherited(inherited bool) {
	p.Inherited = inherited
}

func (p *PropertyCustom) SetSourceId(id uuid.UUID) {
	p.SourceId, _ = uuid.FromString(id.String())
}

func (p *PropertyCustom) SetSourceType(s string) {
	p.SourceType = s
}

func (p PropertyCustom) Clone() Property {
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
	cl.Instances = make([]PropertyInstance, len(p.Instances))
	copy(cl.Instances, p.Instances)

	return &cl
}

func (p *PropertyCustom) MakeAction() Action {
	return Action{
		Property: proto.Property{
			InstanceId:       p.GetID(),
			SourceInstanceId: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			Type:             p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			Custom: &proto.PropertyCustom{
				Id:    p.CustomId.String(),
				Name:  p.Key,
				Value: p.Value,
			},
		},
	}
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
	Attributes    []proto.ServiceAttribute
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

func (p *PropertyService) GetKey() string {
	return p.Service
}

// service has no Value per se, so ensure comparing values never
// succeeds, but Interface is fulfilled
func (p *PropertyService) GetValue() string {
	return p.Id.String()
}

func (p *PropertyService) GetInstanceId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(p.Id, uuid.Nil) {
		return p.Id
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectId, objId) {
			log.Printf("tree.Property.GetInstanceId() found existing instance: %s\n", instance.InstanceId)
			return instance.InstanceId
		}
	}
	return uuid.Nil
}

func (p *PropertyService) SetId(id uuid.UUID) {
	p.Id, _ = uuid.FromString(id.String())
}

func (p *PropertyService) Equal(id uuid.UUID) bool {
	return uuid.Equal(p.Id, id)
}

func (p *PropertyService) clearInstances() {
	p.Instances = nil
}

func (p *PropertyService) SetInheritedFrom(id uuid.UUID) {
	p.InheritedFrom, _ = uuid.FromString(id.String())
}

func (p *PropertyService) SetInherited(inherited bool) {
	p.Inherited = inherited
}

func (p *PropertyService) SetSourceId(id uuid.UUID) {
	p.SourceId, _ = uuid.FromString(id.String())
}

func (p *PropertyService) SetSourceType(s string) {
	p.SourceType = s
}

func (p PropertyService) Clone() Property {
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
	cl.Attributes = make([]proto.ServiceAttribute, 0)
	for _, attr := range p.Attributes {
		a := proto.ServiceAttribute{
			Name:  attr.Name,
			Value: attr.Value,
		}
		cl.Attributes = append(cl.Attributes, a)
	}
	cl.Instances = make([]PropertyInstance, len(p.Instances))
	copy(cl.Instances, p.Instances)

	return &cl
}

func (p *PropertyService) MakeAction() Action {
	a := Action{
		Property: proto.Property{
			InstanceId:       p.GetID(),
			SourceInstanceId: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			Type:             p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			Service: &proto.PropertyService{
				Name: p.Service,
			},
		},
	}
	a.Property.Service.Attributes = make([]proto.ServiceAttribute, len(p.Attributes))
	for i, attr := range p.Attributes {
		t := proto.ServiceAttribute{
			Name:  attr.Name,
			Value: attr.Value,
		}
		a.Property.Service.Attributes[i] = t
	}
	return a
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

func (p *PropertySystem) GetKey() string {
	return p.Key
}

func (p *PropertySystem) GetValue() string {
	return p.Value
}

func (p *PropertySystem) GetInstanceId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(p.Id, uuid.Nil) {
		return p.Id
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectId, objId) {
			log.Printf("tree.Property.GetInstanceId() found existing instance: %s\n", instance.InstanceId)
			return instance.InstanceId
		}
	}
	return uuid.Nil
}

func (p *PropertySystem) SetId(id uuid.UUID) {
	p.Id, _ = uuid.FromString(id.String())
}

func (p *PropertySystem) Equal(id uuid.UUID) bool {
	return uuid.Equal(p.Id, id)
}

func (p *PropertySystem) clearInstances() {
	p.Instances = nil
}

func (p *PropertySystem) SetInheritedFrom(id uuid.UUID) {
	p.InheritedFrom, _ = uuid.FromString(id.String())
}

func (p *PropertySystem) SetInherited(inherited bool) {
	p.Inherited = inherited
}

func (p *PropertySystem) SetSourceId(id uuid.UUID) {
	p.SourceId, _ = uuid.FromString(id.String())
}

func (p *PropertySystem) SetSourceType(s string) {
	p.SourceType = s
}

func (p PropertySystem) Clone() Property {
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
	cl.Instances = make([]PropertyInstance, len(p.Instances))
	copy(cl.Instances, p.Instances)

	return &cl
}

func (p *PropertySystem) MakeAction() Action {
	return Action{
		Property: proto.Property{
			InstanceId:       p.GetID(),
			SourceInstanceId: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			Type:             p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			System: &proto.PropertySystem{
				Name:  p.Key,
				Value: p.Value,
			},
		},
	}
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
	return "oncall"
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

func (p *PropertyOncall) GetKey() string {
	return p.OncallId.String()
}

func (p *PropertyOncall) GetValue() string {
	return p.Name + p.Number
}

func (p *PropertyOncall) GetName() string {
	return p.Name
}

func (p *PropertyOncall) GetNumber() string {
	return p.Number
}

func (p *PropertyOncall) GetInstanceId(objType string, objId uuid.UUID) uuid.UUID {
	if !uuid.Equal(p.Id, uuid.Nil) {
		return p.Id
	}
	for _, instance := range p.Instances {
		if objType == instance.ObjectType && uuid.Equal(instance.ObjectId, objId) {
			log.Printf("tree.Property.GetInstanceId() found existing instance: %s\n", instance.InstanceId)
			return instance.InstanceId
		}
	}
	return uuid.Nil
}

func (p *PropertyOncall) SetId(id uuid.UUID) {
	p.Id, _ = uuid.FromString(id.String())
}

func (p *PropertyOncall) Equal(id uuid.UUID) bool {
	return uuid.Equal(p.Id, id)
}

func (p *PropertyOncall) clearInstances() {
	p.Instances = nil
}

func (p *PropertyOncall) SetInheritedFrom(id uuid.UUID) {
	p.InheritedFrom, _ = uuid.FromString(id.String())
}

func (p *PropertyOncall) SetInherited(inherited bool) {
	p.Inherited = inherited
}

func (p *PropertyOncall) SetSourceId(id uuid.UUID) {
	p.SourceId, _ = uuid.FromString(id.String())
}

func (p *PropertyOncall) SetSourceType(s string) {
	p.SourceType = s
}

func (p PropertyOncall) Clone() Property {
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
	cl.Instances = make([]PropertyInstance, len(p.Instances))
	copy(cl.Instances, p.Instances)

	return &cl
}

func (p *PropertyOncall) MakeAction() Action {
	return Action{
		Property: proto.Property{
			InstanceId:       p.GetID(),
			SourceInstanceId: p.GetSourceInstance(),
			SourceType:       p.GetSourceType(),
			IsInherited:      p.GetIsInherited(),
			InheritedFrom:    p.GetSource(),
			Type:             p.GetType(),
			Inheritance:      p.hasInheritance(),
			ChildrenOnly:     p.isChildrenOnly(),
			View:             p.GetView(),
			Oncall: &proto.PropertyOncall{
				Id:     p.OncallId.String(),
				Name:   p.Name,
				Number: p.Number,
			},
		},
	}
}

func isDupe(o, n Property) (bool, bool, Property) {
	var dupe, deleteOK bool
	var prop Property

	if o.GetKey() == n.GetKey() {
		// not allowed to replace view any with a more
		// specific view or vice versa. Replacing any with any
		// is fine
		if (o.GetView() == `any` && n.GetView() != `any`) ||
			(o.GetView() != `any` && n.GetView() == `any`) {
			// not actually a dupe, but trigger error path
			dupe = true
			deleteOK = false
		}
		// same view means we have a duplicate
		if o.GetView() == n.GetView() {
			dupe = true
			prop = o.Clone()
			// inherited properties can be deleted and replaced
			if o.GetIsInherited() {
				deleteOK = true
			}
		}
	}
	return dupe, deleteOK, prop
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
