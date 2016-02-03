package somatree

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/satori/go.uuid"
)

type SomaTreeElemGroup struct {
	Id       uuid.UUID
	Name     string
	State    string
	Team     uuid.UUID
	Type     string
	Parent   SomaTreeGroupReceiver `json:"-"`
	Fault    *SomaTreeElemFault    `json:"-"`
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
	teg.State = "floating"
	teg.Parent = nil
	teg.Children = make(map[string]SomaTreeGroupAttacher)
	//teg.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//teg.PropertyService = make(map[string]*SomaTreePropertyService)
	//teg.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//teg.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//teg.Checks = make(map[string]*SomaTreeCheck)

	return teg
}

func (teg SomaTreeElemGroup) CloneBucket() SomaTreeBucketAttacher {
	for k, child := range teg.Children {
		teg.Children[k] = child.CloneGroup()
	}
	return &teg
}

func (teg SomaTreeElemGroup) CloneGroup() SomaTreeGroupAttacher {
	f := make(map[string]SomaTreeGroupAttacher)
	for k, child := range teg.Children {
		f[k] = child.CloneGroup()
	}
	teg.Children = f
	return &teg
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
	if teg.Parent != nil {
		panic(`SomaTreeElemGroup.Attach: already attached`)
	}
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
	if teg.Parent == nil {
		panic(`SomaTreeElemGroup.ReAttach: not attached`)
	}
	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(SomaTreeBuilder).GetType(),
		ParentName: teg.Parent.(SomaTreeBuilder).GetName(),
		ParentId:   teg.Parent.(SomaTreeBuilder).GetID(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildId:    teg.GetID(),
	},
	)

	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teg.GetType(),
		Group:      teg,
	},
	)
}

func (teg *SomaTreeElemGroup) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case *SomaTreeElemBucket:
		teg.setGroupParent(p.(SomaTreeGroupReceiver))
		teg.State = "standalone"
	case *SomaTreeElemGroup:
		teg.setGroupParent(p.(SomaTreeGroupReceiver))
		teg.State = "grouped"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemGroup.setParent`)
	}
}

// SomaTreeGroupReceiver == can receive Groups as children
func (teg *SomaTreeElemGroup) setGroupParent(p SomaTreeGroupReceiver) {
	teg.Parent = p
}

func (teg *SomaTreeElemGroup) updateParentRecursive(p SomaTreeReceiver) {
	teg.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		go func(str SomaTreeReceiver) {
			defer wg.Done()
			teg.Children[child].updateParentRecursive(str)
		}(teg)
	}
	wg.Wait()
}

func (teg *SomaTreeElemGroup) clearParent() {
	teg.Parent = nil
	teg.State = "floating"
}

func (teg *SomaTreeElemGroup) setFault(f *SomaTreeElemFault) {
	teg.Fault = f
}

func (teg *SomaTreeElemGroup) updateFaultRecursive(f *SomaTreeElemFault) {
	teg.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range teg.Children {
		wg.Add(1)
		go func(ptr *SomaTreeElemFault) {
			defer wg.Done()
			teg.Children[child].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

func (teg *SomaTreeElemGroup) Destroy() {
	if teg.Parent == nil {
		panic(`SomaTreeElemGroup.Destroy called without Parent to unlink from`)
	}

	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   teg.Parent.(SomaTreeBuilder).GetID(),
		ParentName: teg.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildId:    teg.GetID(),
	},
	)

	teg.setFault(nil)
}

func (teg *SomaTreeElemGroup) Detach() {
	if teg.Parent == nil {
		panic(`SomaTreeElemGroup.Destroy called without Parent to detach from`)
	}

	bucket := teg.Parent.(SomaTreeBucketeer).GetBucket()

	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   teg.Parent.(SomaTreeBuilder).GetID(),
		ParentName: teg.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildId:    teg.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(SomaTreeBuilder).GetType(),
		ParentId:   bucket.(SomaTreeBuilder).GetID(),
		ParentName: bucket.(SomaTreeBuilder).GetName(),
		ChildType:  teg.Type,
		Group:      teg,
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
	for child, _ := range teg.Children {
		if teg.Children[child].(SomaTreeBuilder).GetType() == "node" {
			continue loop
		}
		teg.Children[child].(SomaTreeReceiver).Receive(r)
	}
}

//
// Interface: SomaTreeBucketeer
func (teg *SomaTreeElemGroup) GetBucket() SomaTreeReceiver {
	if teg.Parent == nil {
		if teg.Fault == nil {
			panic(`SomaTreeElemGroup.GetBucket called without Parent`)
		} else {
			return teg.Fault
		}
	}
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
	for child, _ := range teg.Children {
		if teg.Children[child].(SomaTreeBuilder).GetType() == "node" {
			continue loop
		}
		teg.Children[child].(SomaTreeUnlinker).Unlink(u)
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
			r.Group.setFault(teg.Fault)
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
					teg.Children[u.ChildId].clearParent()
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
			r.Cluster.setFault(teg.Fault)
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
					teg.Children[u.ChildId].clearParent()
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
			r.Node.setFault(teg.Fault)
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
					teg.Children[u.ChildId].clearParent()
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
