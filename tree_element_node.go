package somatree

import (
	"fmt"
	"reflect"

	"github.com/satori/go.uuid"
)

type SomaTreeElemNode struct {
	Id       uuid.UUID
	Name     string
	AssetId  uuid.UUID
	Team     uuid.UUID
	ServerId uuid.UUID
	State    string
	Online   bool
	Deleted  bool
	Type     string
	Parent   SomaTreeNodeReceiver `json:"-"`
	Fault    *SomaTreeElemFault   `json:"-"`
	//PropertyOncall  map[string]*SomaTreePropertyOncall
	//PropertyService map[string]*SomaTreePropertyService
	//PropertySystem  map[string]*SomaTreePropertySystem
	//PropertyCustom  map[string]*SomaTreePropertyCustom
	//Checks          map[string]*SomaTreeCheck
}

type NodeSpec struct {
	Id       uuid.UUID
	Name     string
	AssetId  uuid.UUID
	Team     uuid.UUID
	ServerId uuid.UUID
	State    string
}

//
// NEW
func NewNode(name string) *SomaTreeElemNode {
	ten := new(SomaTreeElemNode)
	ten.Id = uuid.NewV4()
	ten.Name = name
	ten.Type = "node"
	ten.State = "floating"
	ten.Parent = nil
	//ten.PropertyOncall = make(map[string]*SomaTreePropertyOncall)
	//ten.PropertyService = make(map[string]*SomaTreePropertyService)
	//ten.PropertySystem = make(map[string]*SomaTreePropertySystem)
	//ten.PropertyCustom = make(map[string]*SomaTreePropertyCustom)
	//ten.Checks = make(map[string]*SomaTreeCheck)

	return ten
}

func (ten SomaTreeElemNode) CloneBucket() SomaTreeBucketAttacher {
	return &ten
}

func (ten SomaTreeElemNode) CloneGroup() SomaTreeGroupAttacher {
	return &ten
}

func (ten SomaTreeElemNode) CloneCluster() SomaTreeClusterAttacher {
	return &ten
}

//
// Interface:
func (ten *SomaTreeElemNode) GetID() string {
	return ten.Id.String()
}

func (ten *SomaTreeElemNode) GetName() string {
	return ten.Name
}

func (ten *SomaTreeElemNode) GetType() string {
	return ten.Type
}

//
// Interface: SomaTreeAttacher
func (ten *SomaTreeElemNode) Attach(a AttachRequest) {
	if ten.Parent != nil {
		panic(`SomaTreeElemNode.Attach: already attached`)
	}
	switch {
	case a.ParentType == "bucket":
		ten.attachToBucket(a)
	case a.ParentType == "group":
		ten.attachToGroup(a)
	case a.ParentType == "cluster":
		ten.attachToCluster(a)
	default:
		panic(`SomaTreeElemNode.Attach`)
	}
}

func (ten *SomaTreeElemNode) ReAttach(a AttachRequest) {
	if ten.Parent == nil {
		panic(`SomaTreeElemNode.ReAttach: not attached`)
	}
	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(SomaTreeBuilder).GetType(),
		ParentName: ten.Parent.(SomaTreeBuilder).GetName(),
		ParentId:   ten.Parent.(SomaTreeBuilder).GetID(),
		ChildType:  ten.GetType(),
		ChildName:  ten.GetName(),
		ChildId:    ten.GetID(),
	},
	)

	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  ten.GetType(),
		Node:       ten,
	},
	)
}

func (ten *SomaTreeElemNode) Destroy() {
	if ten.Parent == nil {
		panic(`SomaTreeElemNode.Destroy called without Parent to unlink from`)
	}

	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   ten.Parent.(SomaTreeBuilder).GetID(),
		ParentName: ten.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  ten.GetType(),
		ChildName:  ten.GetName(),
		ChildId:    ten.GetID(),
	},
	)

	ten.setFault(nil)
}

func (ten *SomaTreeElemNode) Detach() {
	if ten.Parent == nil {
		panic(`SomaTreeElemNode.Detach called without Parent to detach from`)
	}

	bucket := ten.Parent.(SomaTreeBucketeer).GetBucket()

	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   ten.Parent.(SomaTreeBuilder).GetID(),
		ParentName: ten.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  ten.GetType(),
		ChildName:  ten.GetName(),
		ChildId:    ten.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(SomaTreeBuilder).GetType(),
		ParentId:   bucket.(SomaTreeBuilder).GetID(),
		ParentName: bucket.(SomaTreeBuilder).GetName(),
		ChildType:  ten.Type,
		Node:       ten,
	},
	)
}

func (ten *SomaTreeElemNode) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case *SomaTreeElemBucket:
		ten.setNodeParent(p.(SomaTreeNodeReceiver))
		ten.State = "standalone"
	case *SomaTreeElemGroup:
		ten.setNodeParent(p.(SomaTreeNodeReceiver))
		ten.State = "grouped"
	case *SomaTreeElemCluster:
		ten.setNodeParent(p.(SomaTreeNodeReceiver))
		ten.State = "clustered"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemNode.setParent`)
	}
}

func (ten *SomaTreeElemNode) updateParentRecursive(p SomaTreeReceiver) {
	ten.setParent(p)
}

func (ten *SomaTreeElemNode) setNodeParent(p SomaTreeNodeReceiver) {
	ten.Parent = p
}

func (ten *SomaTreeElemNode) clearParent() {
	ten.Parent = nil
	ten.State = "floating"
}

func (ten *SomaTreeElemNode) setFault(f *SomaTreeElemFault) {
	ten.Fault = f
}

func (ten *SomaTreeElemNode) updateFaultRecursive(f *SomaTreeElemFault) {
	ten.setFault(f)
}

//
// Interface: SomaTreeBucketAttacher
func (ten *SomaTreeElemNode) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})
}

//
// Interface: SomaTreeGroupAttacher
func (ten *SomaTreeElemNode) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})
}

//
// Interface: SomaTreeClusterAttacher
func (ten *SomaTreeElemNode) attachToCluster(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
