package tree

import "sync"

//
// Interface: Attacher
func (teg *Group) Attach(a AttachRequest) {
	if teg.Parent != nil {
		panic(`Group.Attach: already attached`)
	}
	switch {
	case a.ParentType == "bucket":
		teg.attachToBucket(a)
	case a.ParentType == "group":
		teg.attachToGroup(a)
	default:
		panic(`Group.Attach`)
	}

	teg.Parent.(Propertier).syncProperty(teg.Id.String())
}

func (teg *Group) ReAttach(a AttachRequest) {
	if teg.Parent == nil {
		panic(`Group.ReAttach: not attached`)
	}
	teg.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(Builder).GetType(),
		ParentName: teg.Parent.(Builder).GetName(),
		ParentId:   teg.Parent.(Builder).GetID(),
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

	if teg.Parent == nil {
		panic(`Group.ReAttach: not reattached`)
	}
	teg.actionUpdate()
	teg.Parent.(Propertier).syncProperty(teg.Id.String())
}

func (teg *Group) Destroy() {
	if teg.Parent == nil {
		panic(`Group.Destroy called without Parent to unlink from`)
	}

	// call before unlink since it requires teg.Parent.*
	teg.actionDelete()
	teg.deletePropertyAllLocal()
	teg.deletePropertyAllInherited()
	// TODO delete all checks + check instances
	// TODO delete all inherited checks + check instances

	wg := new(sync.WaitGroup)
	for child, _ := range teg.Children {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			teg.Children[c].Destroy()
		}(child)
	}
	wg.Wait()

	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(Builder).GetType(),
		ParentId:   teg.Parent.(Builder).GetID(),
		ParentName: teg.Parent.(Builder).GetName(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildId:    teg.GetID(),
	},
	)

	teg.setFault(nil)
	teg.setAction(nil)
}

func (teg *Group) Detach() {
	if teg.Parent == nil {
		panic(`Group.Destroy called without Parent to detach from`)
	}
	bucket := teg.Parent.(Bucketeer).GetBucket()

	teg.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	teg.Parent.Unlink(UnlinkRequest{
		ParentType: teg.Parent.(Builder).GetType(),
		ParentId:   teg.Parent.(Builder).GetID(),
		ParentName: teg.Parent.(Builder).GetName(),
		ChildType:  teg.GetType(),
		ChildName:  teg.GetName(),
		ChildId:    teg.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(Builder).GetType(),
		ParentId:   bucket.(Builder).GetID(),
		ParentName: bucket.(Builder).GetName(),
		ChildType:  teg.Type,
		Group:      teg,
	},
	)

	teg.actionUpdate()
	teg.Parent.(Propertier).syncProperty(teg.Id.String())
}

//
// Interface: BucketAttacher
func (teg *Group) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teg.Type,
		Group:      teg,
	})

	if teg.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_group`})
		return
	}
	teg.actionCreate()
}

//
// Interface: GroupAttacher
func (teg *Group) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teg.Type,
		Group:      teg,
	})

	if teg.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_group`})
		return
	}
	teg.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
