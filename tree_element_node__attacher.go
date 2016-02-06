package somatree

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
	ten.setAction(nil)
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

	ten.Action <- &Action{
		Action: "create",
		Type:   "node",
		Id:     ten.Id.String(),
		Name:   ten.Name,
		Team:   ten.Team.String(),
	}
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

	ten.Action <- &Action{
		Action: "create",
		Type:   "node",
		Id:     ten.Id.String(),
		Name:   ten.Name,
		Team:   ten.Team.String(),
	}
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

	ten.Action <- &Action{
		Action: "create",
		Type:   "node",
		Id:     ten.Id.String(),
		Name:   ten.Name,
		Team:   ten.Team.String(),
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
