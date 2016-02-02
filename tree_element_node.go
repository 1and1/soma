package somatree

import "github.com/satori/go.uuid"

type SomaTreeElemNode struct {
	Id              uuid.UUID
	Name            string
	Parent          SomaTreeNodeReceiver
	PropertyOncall  map[string]*SomaTreePropertyOncall
	PropertyService map[string]*SomaTreePropertyService
	PropertySystem  map[string]*SomaTreePropertySystem
	PropertyCustom  map[string]*SomaTreePropertyCustom
	Checks          map[string]*SomaTreeCheck
}

func NewNode() *SomaTreeElemNode {
	return new(SomaTreeElemNode)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
