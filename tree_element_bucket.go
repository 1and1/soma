package somatree

import (
	"fmt"
	"log"
	"reflect"

	"github.com/satori/go.uuid"
)

type SomaTreeElemBucket struct {
	Id          uuid.UUID
	Name        string
	Environment string
	Type        string
	State       string
	Parent      SomaTreeBucketReceiver `json:"-"`
	Children    map[string]SomaTreeBucketAttacher
	//PropertyOncall  map[string]*SomaTreePropertyOncall
	//PropertyService map[string]*SomaTreePropertyService
	//PropertySystem  map[string]*SomaTreePropertySystem
	//PropertyCustom  map[string]*SomaTreePropertyCustom
	//Checks          map[string]*SomaTreeCheck
}

//
// NEW
func NewBucket(name string, environment string, id string) *SomaTreeElemBucket {
	teb := new(SomaTreeElemBucket)
	if id == "" {
		teb.Id = uuid.NewV4()
	} else {
		teb.Id, _ = uuid.FromString(id)
	}
	teb.Name = name
	teb.Environment = environment
	teb.Type = "bucket"
	teb.State = "floating"
	teb.Children = make(map[string]SomaTreeBucketAttacher)
	//teb.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//teb.PropertyService = make(map[string]*SomaTreePropertyService)
	//teb.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//teb.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//teb.Checks = make(map[string]*SomaTreeCheck)

	return teb
}

func (teb SomaTreeElemBucket) CloneRepository() SomaTreeRepositoryAttacher {
	f := make(map[string]SomaTreeBucketAttacher)
	for k, child := range teb.Children {
		f[k] = child.CloneBucket()
	}
	teb.Children = f
	return &teb
}

//
// Interface: SomaTreeBuilder
func (teb *SomaTreeElemBucket) GetID() string {
	return teb.Id.String()
}

func (teb *SomaTreeElemBucket) GetName() string {
	return teb.Name
}

func (teb *SomaTreeElemBucket) GetType() string {
	return teb.Type
}

//
// Interface: SomaTreeAttacher
func (teb *SomaTreeElemBucket) Attach(a AttachRequest) {
	switch {
	case a.ParentType == "repository":
		teb.attachToRepository(a)
	}
}

func (teb *SomaTreeElemBucket) ReAttach(a AttachRequest) {
	log.Fatal("Not implemented")
}

func (teb *SomaTreeElemBucket) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeBucketReceiver:
		teb.setBucketParent(p.(SomaTreeBucketReceiver))
		teb.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemBucket.setParent`)
	}
}

func (teb *SomaTreeElemBucket) setBucketParent(p SomaTreeBucketReceiver) {
	teb.Parent = p
}

func (teb *SomaTreeElemBucket) clearParent() {
	teb.Parent = nil
	teb.State = "floating"
}

func (teb *SomaTreeElemBucket) Destroy() {
	if teb.Parent == nil {
		panic(`SomaTreeElemBucket.Destroy called without Parent to unlink from`)
	}

	teb.Parent.Unlink(UnlinkRequest{
		ParentType: teb.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   teb.Parent.(SomaTreeBuilder).GetID(),
		ParentName: teb.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  teb.GetType(),
		ChildName:  teb.GetName(),
		ChildId:    teb.GetID(),
	},
	)
}

func (teb *SomaTreeElemBucket) Detach() {
	teb.Destroy()
}

//
// Interface: SomaTreeRepositoryAttacher
func (teb *SomaTreeElemBucket) attachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teb.Type,
		Bucket:     teb,
	})
}

//
// Interface: SomaTreeReceiver
func (teb *SomaTreeElemBucket) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "group":
			teb.receiveGroup(r)
		case "cluster":
			teb.receiveCluster(r)
		case "node":
			teb.receiveNode(r)
		default:
			panic(`SomaTreeElemBucket.Receive`)
		}
		return
	}
loop:
	for _, child := range teb.Children {
		if child.(SomaTreeBuilder).GetType() == "node" {
			continue loop
		}
		child.(SomaTreeReceiver).Receive(r)
	}
}

//
// Interface: SomaTreeBucketeer
func (teb *SomaTreeElemBucket) GetBucket() SomaTreeReceiver {
	return teb
}

//
// Interface: SomaTreeUnlinker
func (teb *SomaTreeElemBucket) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "group":
			teb.unlinkGroup(u)
		case "cluster":
			teb.unlinkCluster(u)
		case "node":
			teb.unlinkNode(u)
		default:
			panic(`SomaTreeElemBucket.Unlink`)
		}
		return
	}
loop:
	for _, child := range teb.Children {
		if child.(SomaTreeBuilder).GetType() == "node" {
			continue loop
		}
		child.(SomaTreeUnlinker).Unlink(u)
	}
}

//
// Interface: SomaTreeGroupReceiver
func (teb *SomaTreeElemBucket) receiveGroup(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "group":
			teb.Children[r.Group.GetID()] = r.Group
			r.Group.setParent(teb)
		default:
			panic(`SomaTreeElemBucket.receiveGroup`)
		}
		return
	}
	panic(`SomaTreeElemBucket.receiveGroup`)
}

//
// Interface: SomaTreeGroupUnlinker
func (teb *SomaTreeElemBucket) unlinkGroup(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "group":
			if _, ok := teb.Children[u.ChildId]; ok {
				if u.ChildName == teb.Children[u.ChildId].GetName() {
					teb.Children[u.ChildId].clearParent()
					delete(teb.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemBucket.unlinkGroup`)
		}
		return
	}
	panic(`SomaTreeElemBucket.unlinkGroup`)
}

//
// Interface: SomaTreeClusterReceiver
func (teb *SomaTreeElemBucket) receiveCluster(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "cluster":
			teb.Children[r.Cluster.GetID()] = r.Cluster
			r.Cluster.setParent(teb)
		default:
			panic(`SomaTreeElemBucket.receiveCluster`)
		}
		return
	}
	panic(`SomaTreeElemBucket.receiveCluster`)
}

//
// Interface: SomaTreeClusterUnlinker
func (teb *SomaTreeElemBucket) unlinkCluster(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "cluster":
			if _, ok := teb.Children[u.ChildId]; ok {
				if u.ChildName == teb.Children[u.ChildId].GetName() {
					teb.Children[u.ChildId].clearParent()
					delete(teb.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemBucket.unlinkCluster`)
		}
		return
	}
	panic(`SomaTreeElemBucket.unlinkCluster`)
}

//
// Interface: SomaTreeNodeReceiver
func (teb *SomaTreeElemBucket) receiveNode(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "node":
			teb.Children[r.Node.GetID()] = r.Node
			r.Node.setParent(teb)
		default:
			panic(`SomaTreeElemBucket.receiveNote`)
		}
		return
	}
	panic(`SomaTreeElemBucket.receiveNote`)
}

//
// Interface: SomaTreeNodeUnlinker
func (teb *SomaTreeElemBucket) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "node":
			if _, ok := teb.Children[u.ChildId]; ok {
				if u.ChildName == teb.Children[u.ChildId].GetName() {
					teb.Children[u.ChildId].clearParent()
					delete(teb.Children, u.ChildId)
				}
			}
		default:
			panic(`SomaTreeElemBucket.unlinkNode`)
		}
		return
	}
	panic(`SomaTreeElemBucket.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
