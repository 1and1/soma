package somatree

import (
	"fmt"
	"log"
	"reflect"

	"github.com/satori/go.uuid"
)

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

//
// Interface: SomaTreeAttacher
func (tec *SomaTreeElemCluster) Attach(a AttachRequest) {
	switch {
	case a.ParentType == "bucket":
		tec.attachToBucket(a)
	case a.ParentType == "group":
		tec.attachToGroup(a)
	default:
		panic(`SomaTreeElemCluster.Attach`)
	}
}

func (tec *SomaTreeElemCluster) ReAttach(a AttachRequest) {
	log.Fatal("Not implemented")
}

func (tec *SomaTreeElemCluster) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case *SomaTreeElemBucket:
		tec.setClusterParent(p.(SomaTreeClusterReceiver))
	case *SomaTreeElemGroup:
		tec.setClusterParent(p.(SomaTreeClusterReceiver))
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemCluster.setParent`)
	}
}

// SomaTreeClusterReceiver == can receive Clusters as children
func (tec *SomaTreeElemCluster) setClusterParent(p SomaTreeClusterReceiver) {
	tec.Parent = p
}

func (tec *SomaTreeElemCluster) Destroy() {
	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   tec.Parent.(SomaTreeBuilder).GetID(),
		ParentName: tec.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)
}

//
// Interface: SomaTreeBucketAttacher
func (tec *SomaTreeElemCluster) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tec.Type,
		Cluster:    tec,
	})
}

//
// Interface: SomaTreeGroupAttacher
func (tec *SomaTreeElemCluster) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tec.Type,
		Cluster:    tec,
	})
}

//
// Interface: SomaTreeReceiver
func (tec *SomaTreeElemCluster) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, tec) {
		switch r.ChildType {
		case "node":
			tec.receiveNode(r)
		default:
			panic(`SomaTreeElemCluster.Receive`)
		}
	}
	// no passing along since only nodes are a SomeTreeClusterAttacher
	// and nodes can have no children
	return
}

//
// Interface: SomaTreeUnlinker
func (tec *SomaTreeElemCluster) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, tec) {
		switch u.ChildType {
		case "node":
			tec.unlinkNode(u)
		default:
			panic(`SomaTreeElemCluster.Unlink`)
		}
	}
	// no passing along since only nodes are a SomeTreeClusterAttacher
	// and nodes can have no children
	return
}

//
// Interface: SomaTreeNodeReceiver
func (tec *SomaTreeElemCluster) receiveNode(r ReceiveRequest) {
	if receiveRequestCheck(r, tec) {
		switch r.ChildType {
		case "node":
			tec.Children[r.Node.GetID()] = r.Node
			r.Node.setParent(tec)
		default:
			panic(`SomaTreeElemCluster.receiveNode`)
		}
		return
	}
	panic(`SomaTreeElemCluster.receiveNode`)
}

//
// Interface: SomaTreeNodeUnlinker
func (tec *SomaTreeElemCluster) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, tec) {
		switch u.ChildType {
		case "node":
			if _, ok := tec.Children[u.ChildId]; ok {
				if u.ChildName == tec.Children[u.ChildId].GetName() {
					delete(tec.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemCluster.unlinkNode`)
		}
		return
	}
	panic(`SomaTreeElemCluster.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
