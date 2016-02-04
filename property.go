package somatree

import "github.com/satori/go.uuid"

type SomaTreeProperty interface {
	GetType() string
	GetID() string
}

type SomaTreePropertyOncall struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
}

type SomaTreePropertyService struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
}

type SomaTreePropertySystem struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
}

type SomaTreePropertyCustom struct {
	Id            uuid.UUID
	Inherited     bool
	InheritedFrom uuid.UUID
}

func (p *SomaTreePropertyCustom) GetType() string {
	return "custom"
}

func (p *SomaTreePropertyCustom) GetID() string {
	return p.Id.String()
}

func (p *SomaTreePropertyService) GetType() string {
	return "custom"
}

func (p *SomaTreePropertyService) GetID() string {
	return p.Id.String()
}

func (p *SomaTreePropertySystem) GetType() string {
	return "custom"
}

func (p *SomaTreePropertySystem) GetID() string {
	return p.Id.String()
}

func (p *SomaTreePropertyOncall) GetType() string {
	return "custom"
}

func (p *SomaTreePropertyOncall) GetID() string {
	return p.Id.String()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
