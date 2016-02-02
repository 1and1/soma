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

func (ten *SomaTreeElemNode) GetID() string {
	return ten.Id.String()
}

func (ten *SomaTreeElemNode) GetName() string {
	return ten.Name
}

func (ten *SomaTreeElemNode) Attach(a AttachRequest) {
}

func (ten *SomaTreeElemNode) ReAttach(a AttachRequest) {
}

func (ten *SomaTreeElemNode) setParent(p SomaTreeReceiver) {
}

func (ten *SomaTreeElemNode) attachToBucket(a AttachRequest) {
}

func (ten *SomaTreeElemNode) attachToGroup(a AttachRequest) {
}

func (ten *SomaTreeElemNode) attachToCluster(a AttachRequest) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
