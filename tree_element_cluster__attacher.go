package somatree

import "sync"

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
	if tec.Parent == nil {
		panic(`SomaTreeElemGroup.ReAttach: not attached`)
	}
	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(SomaTreeBuilder).GetType(),
		ParentName: tec.Parent.(SomaTreeBuilder).GetName(),
		ParentId:   tec.Parent.(SomaTreeBuilder).GetID(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)

	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tec.GetType(),
		Cluster:    tec,
	},
	)

	tec.actionUpdate()
}

func (tec *SomaTreeElemCluster) Destroy() {
	if tec.Parent == nil {
		panic(`SomaTreeElemCluster.Destroy called without Parent to unlink from`)
	}

	wg := new(sync.WaitGroup)
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func() {
			defer wg.Done()
			tec.Children[c].Destroy()
		}()
	}
	wg.Wait()

	// call before unlink since it requires tec.Parent.*
	tec.actionDelete()

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   tec.Parent.(SomaTreeBuilder).GetID(),
		ParentName: tec.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)

	tec.setFault(nil)
	tec.setAction(nil)
}

func (tec *SomaTreeElemCluster) Detach() {
	if tec.Parent == nil {
		panic(`SomaTreeElemCluster.Detach called without Parent to detach from`)
	}
	bucket := tec.Parent.(SomaTreeBucketeer).GetBucket()

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   tec.Parent.(SomaTreeBuilder).GetID(),
		ParentName: tec.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(SomaTreeBuilder).GetType(),
		ParentId:   bucket.(SomaTreeBuilder).GetID(),
		ParentName: bucket.(SomaTreeBuilder).GetName(),
		ChildType:  tec.Type,
		Cluster:    tec,
	},
	)

	tec.actionUpdate()
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

	tec.actionCreate()
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

	tec.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
