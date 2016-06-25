package somatree

import "sync"

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

	if teg.Parent == nil {
		panic(`SomaTreeElemGroup.Attach: failed`)
	}
	teg.Parent.(Propertier).syncProperty(teg.Id.String())
}

func (teg *SomaTreeElemGroup) ReAttach(a AttachRequest) {
	if teg.Parent == nil {
		panic(`SomaTreeElemGroup.ReAttach: not attached`)
	}
	// XXX: destroy all inherited properties before unlinking
	// teg.(SomaTreePropertier).destroyInheritedProperties()

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
		panic(`SomaTreeElemGroup.ReAttach: not reattached`)
	}
	teg.actionUpdate()
	teg.Parent.(Propertier).syncProperty(teg.Id.String())
}

func (teg *SomaTreeElemGroup) Destroy() {
	if teg.Parent == nil {
		panic(`SomaTreeElemGroup.Destroy called without Parent to unlink from`)
	}

	wg := new(sync.WaitGroup)
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			teg.Children[c].Destroy()
		}()
	}
	wg.Wait()

	// call before unlink since it requires teg.Parent.*
	teg.actionDelete()

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

func (teg *SomaTreeElemGroup) Detach() {
	if teg.Parent == nil {
		panic(`SomaTreeElemGroup.Destroy called without Parent to detach from`)
	}
	bucket := teg.Parent.(Bucketeer).GetBucket()

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

	if teg.Parent == nil {
		a.Root.(*SomaTree).AttachError(Error{Action: `attach_group`})
		return
	}
	teg.actionCreate()
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

	if teg.Parent == nil {
		a.Root.(*SomaTree).AttachError(Error{Action: `attach_group`})
		return
	}
	teg.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
