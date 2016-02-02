package somatree

import "github.com/satori/go.uuid"

type SomaTreeElemGroup struct {
	Id              uuid.UUID
	Name            string
	Parent          SomaTreeGroupReceiver
	Children        map[string]SomaTreeGroupAttacher
	PropertyOncall  map[string]*SomaTreePropertyOncall
	PropertyService map[string]*SomaTreePropertyService
	PropertySystem  map[string]*SomaTreePropertySystem
	PropertyCustom  map[string]*SomaTreePropertyCustom
	Checks          map[string]*SomaTreeCheck
}

func NewGroup() *SomaTreeElemGroup {
	return new(SomaTreeElemGroup)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
