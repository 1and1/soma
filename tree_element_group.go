package somatree

import (
	"fmt"
	"log"
	"reflect"

	"github.com/satori/go.uuid"
)

type SomaTreeElemGroup struct {
	Id       uuid.UUID
	Name     string
	State    string
	Team     uuid.UUID
	Type     string
	Parent   SomaTreeGroupReceiver `json:"-"`
	Children map[string]SomaTreeGroupAttacher
	//PropertyOncall  map[string]*SomaTreePropertyOncall
	//PropertyService map[string]*SomaTreePropertyService
	//PropertySystem  map[string]*SomaTreePropertySystem
	//PropertyCustom  map[string]*SomaTreePropertyCustom
	//Checks          map[string]*SomaTreeCheck
}

type GroupSepc struct {
	Id   uuid.UUID
	Name string
	Team uuid.UUID
}

//
// NEW
func NewGroup(name string) *SomaTreeElemGroup {
	teg := new(SomaTreeElemGroup)
	teg.Id = uuid.NewV4()
	teg.Name = name
	teg.Type = "group"
	teg.Children = make(map[string]SomaTreeGroupAttacher)
	//teg.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//teg.PropertyService = make(map[string]*SomaTreePropertyService)
	//teg.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//teg.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//teg.Checks = make(map[string]*SomaTreeCheck)

	return teg
}

//
// Interface: SomaTreeBuilder
func (teg *SomaTreeElemGroup) GetID() string {
	return teg.Id.String()
}

func (teg *SomaTreeElemGroup) GetName() string {
	return teg.Name
}

func (teg *SomaTreeElemGroup) GetType() string {
	return teg.Type
}

//
// Interface: SomaTreeAttacher
func (teg *SomaTreeElemGroup) Attach(a AttachRequest) {
	switch {
	case a.ParentType == "bucket":
		teg.attachToBucket(a)
	case a.ParentType == "group":
		teg.attachToGroup(a)
	default:
		panic(`SomaTreeElemGroup.Attach`)
	}
}

func (teg *SomaTreeElemGroup) ReAttach(a AttachRequest) {
	log.Fatal("Not implemented")
}

func (teg *SomaTreeElemGroup) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case *SomaTreeElemBucket:
		teg.setGroupParent(p.(SomaTreeGroupReceiver))
	case *SomaTreeElemGroup:
		teg.setGroupParent(p.(SomaTreeGroupReceiver))
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemGroup.setParent`)
	}
}

// SomaTreeGroupReceiver == can receive Groups as children
func (teg *SomaTreeElemGroup) setGroupParent(p SomaTreeGroupReceiver) {
	teg.Parent = p
}

func (teg *SomaTreeElemGroup) Destroy() {
	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   teg.Parent.(SomaTreeBuilder).GetID(),
		ParentName: teg.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildId:    teg.GetID(),
	},
	)
}

//
// Interface: SomaTreeBucketAttacher
func (teg *SomaTreeElemGroup) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teg.Type,
		Group:      teg,
	})
}

//
// Interface: SomaTreeGroupAttacher
func (teg *SomaTreeElemGroup) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teg.Type,
		Group:      teg,
	})
}

//
// Interface: SomaTreeReceiver
func (teg *SomaTreeElemGroup) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "group":
			teg.receiveGroup(r)
		case "cluster":
			teg.receiveCluster(r)
		case "node":
			teg.receiveNode(r)
		default:
			panic(`SomaTreeElemGroup.Receive`)
		}
		return
	}
loop:
	for _, child := range teg.Children {
		if child.(SomaTreeBuilder).GetType() == "node" {
			continue loop
		}
		child.(SomaTreeReceiver).Receive(r)
	}
}

//
// Interface: SomaTreeBucketeer
func (teg *SomaTreeElemGroup) GetBucket() SomaTreeReceiver {
	return teg.Parent.(SomaTreeBucketeer).GetBucket()
}

//
// Interface: SomaTreeUnlinker
func (teg *SomaTreeElemGroup) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "group":
			teg.unlinkGroup(u)
		case "cluster":
			teg.unlinkCluster(u)
		case "node":
			teg.unlinkNode(u)
		default:
			panic(`SomaTreeElemGroup.Unlink`)
		}
		return
	}
loop:
	for _, child := range teg.Children {
		if child.(SomaTreeBuilder).GetType() == "node" {
			continue loop
		}
		child.(SomaTreeUnlinker).Unlink(u)
	}
}

//
// Interface: SomaTreeGroupReceiver
func (teg *SomaTreeElemGroup) receiveGroup(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "group":
			teg.Children[r.Group.GetID()] = r.Group
			r.Group.setParent(teg)
		default:
			panic(`SomaTreeElemGroup.receiveGroup`)
		}
		return
	}
	panic(`SomaTreeElemGroup.receiveGroup`)
}

//
// Interface: SomaTreeGroupUnlinker
func (teg *SomaTreeElemGroup) unlinkGroup(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "group":
			if _, ok := teg.Children[u.ChildId]; ok {
				if u.ChildName == teg.Children[u.ChildId].GetName() {
					delete(teg.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemGroup.unlinkGroup`)
		}
		return
	}
	panic(`SomaTreeElemGroup.unlinkGroup`)
}

//
// Interface: SomaTreeClusterReceiver
func (teg *SomaTreeElemGroup) receiveCluster(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "cluster":
			teg.Children[r.Cluster.GetID()] = r.Cluster
			r.Cluster.setParent(teg)
		default:
			panic(`SomaTreeElemGroup.receiveCluster`)
		}
		return
	}
	panic(`SomaTreeElemGroup.receiveCluster`)
}

//
// Interface: SomaTreeClusterUnlinker
func (teg *SomaTreeElemGroup) unlinkCluster(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "cluster":
			if _, ok := teg.Children[u.ChildId]; ok {
				if u.ChildName == teg.Children[u.ChildId].GetName() {
					delete(teg.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemGroup.unlinkCluster`)
		}
		return
	}
	panic(`SomaTreeElemGroup.unlinkCluster`)
}

//
// Interface: SomaTreeNodeReceiver
func (teg *SomaTreeElemGroup) receiveNode(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "node":
			teg.Children[r.Node.GetID()] = r.Node
			r.Node.setParent(teg)
		default:
			panic(`SomaTreeElemGroup.receiveNode`)
		}
		return
	}
	panic(`SomaTreeElemGroup.receiveNode`)
}

//
// Interface: SomaTreeNodeUnlinker
func (teg *SomaTreeElemGroup) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "node":
			if _, ok := teg.Children[u.ChildId]; ok {
				if u.ChildName == teg.Children[u.ChildId].GetName() {
					delete(teg.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemGroup.unlinkNode`)
		}
		return
	}
	panic(`SomaTreeElemGroup.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
