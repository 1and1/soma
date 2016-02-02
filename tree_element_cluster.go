package somatree

import "github.com/satori/go.uuid"

type SomaTreeElemCluster struct {
	Id       uuid.UUID
	Name     string
	State    string
	Team     uuid.UUID
	Type     string
	Parent   SomaTreeClusterReceiver `json:"-"`
	Children map[string]SomaTreeClusterAttacher
	//PropertyOncall  map[string]*SomaTreePropertyOncall
	//PropertyService map[string]*SomaTreePropertyService
	//PropertySystem  map[string]*SomaTreePropertySystem
	//PropertyCustom  map[string]*SomaTreePropertyCustom
	//Checks          map[string]*SomaTreeCheck
}

type ClusterSpec struct {
	Id   uuid.UUID
	Name string
	Team uuid.UUID
}

//
// NEW
func NewCluster(name string) *SomaTreeElemCluster {
	tec := new(SomaTreeElemCluster)
	tec.Id = uuid.NewV4()
	tec.Name = name
	tec.Type = "cluster"
	tec.Children = make(map[string]SomaTreeClusterAttacher)
	//tec.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//tec.PropertyService = make(map[string]*SomaTreePropertyService)
	//tec.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//tec.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//tec.Checks = make(map[string]*SomaTreeCheck)

	return tec
}

//
// Interface: SomaTreeBuilder
func (tec *SomaTreeElemCluster) GetID() string {
	return tec.Id.String()
}

func (tec *SomaTreeElemCluster) GetName() string {
	return tec.Name
}

func (tec *SomaTreeElemCluster) GetType() string {
	return tec.Type
}

func (tec *SomaTreeElemCluster) Attach(a AttachRequest) {
}

func (tec *SomaTreeElemCluster) ReAttach(a AttachRequest) {
}

func (tec *SomaTreeElemCluster) setParent(p SomaTreeReceiver) {
}

func (tec *SomaTreeElemCluster) attachToBucket(a AttachRequest) {
}

func (tec *SomaTreeElemCluster) attachToGroup(a AttachRequest) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
