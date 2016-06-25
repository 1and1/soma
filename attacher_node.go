package tree

//
// Interface: Attacher
func (ten *Node) Attach(a AttachRequest) {
	if ten.Parent != nil {
		panic(`Node.Attach: already attached`)
	}
	switch {
	case a.ParentType == "bucket":
		ten.attachToBucket(a)
	case a.ParentType == "group":
		ten.attachToGroup(a)
	case a.ParentType == "cluster":
		ten.attachToCluster(a)
	default:
		panic(`Node.Attach`)
	}

	if ten.Parent == nil {
		panic(`Node.Attach: failed`)
	}
	ten.Parent.(Propertier).syncProperty(ten.Id.String())
}

func (ten *Node) ReAttach(a AttachRequest) {
	if ten.Parent == nil {
		panic(`Node.ReAttach: not attached`)
	}
	ten.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(Builder).GetType(),
		ParentName: ten.Parent.(Builder).GetName(),
		ParentId:   ten.Parent.(Builder).GetID(),
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

	if ten.Parent == nil {
		panic(`Node.ReAttach: not reattached`)
	}
	ten.actionUpdate()
	ten.Parent.(Propertier).syncProperty(ten.Id.String())
}

func (ten *Node) Destroy() {
	if ten.Parent == nil {
		panic(`Node.Destroy called without Parent to unlink from`)
	}
	// call before unlink since it requires tec.Parent.*
	ten.actionDelete()
	ten.deletePropertyAllLocal()
	ten.deletePropertyAllInherited()
	// TODO delete all checks + check instances
	// TODO delete all inherited checks + check instances

	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(Builder).GetType(),
		ParentId:   ten.Parent.(Builder).GetID(),
		ParentName: ten.Parent.(Builder).GetName(),
		ChildType:  ten.GetType(),
		ChildName:  ten.GetName(),
		ChildId:    ten.GetID(),
	},
	)

	ten.setFault(nil)
	ten.setAction(nil)
}

func (ten *Node) Detach() {
	if ten.Parent == nil {
		panic(`Node.Detach called without Parent to detach from`)
	}
	bucket := ten.Parent.(Bucketeer).GetBucket()

	ten.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	ten.Parent.Unlink(UnlinkRequest{
		ParentType: ten.Parent.(Builder).GetType(),
		ParentId:   ten.Parent.(Builder).GetID(),
		ParentName: ten.Parent.(Builder).GetName(),
		ChildType:  ten.GetType(),
		ChildName:  ten.GetName(),
		ChildId:    ten.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(Builder).GetType(),
		ParentId:   bucket.(Builder).GetID(),
		ParentName: bucket.(Builder).GetName(),
		ChildType:  ten.Type,
		Node:       ten,
	},
	)

	ten.actionUpdate()
	ten.Parent.(Propertier).syncProperty(ten.Id.String())
}

//
// Interface: BucketAttacher
func (ten *Node) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})

	if ten.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_node`})
		return
	}
	ten.actionUpdate()
}

//
// Interface: GroupAttacher
func (ten *Node) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})

	if ten.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_node`})
		return
	}
	ten.actionUpdate()
}

//
// Interface: ClusterAttacher
func (ten *Node) attachToCluster(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  ten.Type,
		Node:       ten,
	})

	if ten.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_node`})
		return
	}
	ten.actionUpdate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
