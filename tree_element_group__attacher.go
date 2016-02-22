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

	teg.actionUpdate()
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
		ParentType: teg.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   teg.Parent.(SomaTreeBuilder).GetID(),
		ParentName: teg.Parent.(SomaTreeBuilder).GetName(),
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

	teg.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
