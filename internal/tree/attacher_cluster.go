/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import "sync"

//
// Interface: Attacher
func (tec *Cluster) Attach(a AttachRequest) {
	if tec.Parent != nil {
		panic(`Cluster.Attach: already attached`)
	}
	switch {
	case a.ParentType == "bucket":
		tec.attachToBucket(a)
	case a.ParentType == "group":
		tec.attachToGroup(a)
	default:
		panic(`Cluster.Attach`)
	}

	tec.Parent.(Propertier).syncProperty(tec.Id.String())
	tec.Parent.(Checker).syncCheck(tec.Id.String())
}

func (tec *Cluster) ReAttach(a AttachRequest) {
	if tec.Parent == nil {
		panic(`Cluster.ReAttach: not attached`)
	}
	tec.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(Builder).GetType(),
		ParentName: tec.Parent.(Builder).GetName(),
		ParentId:   tec.Parent.(Builder).GetID(),
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

	if tec.Parent == nil {
		panic(`Group.ReAttach: not reattached`)
	}
	tec.actionUpdate()
	tec.Parent.(Propertier).syncProperty(tec.Id.String())
	tec.Parent.(Checker).syncCheck(tec.Id.String())
}

func (tec *Cluster) Destroy() {
	if tec.Parent == nil {
		panic(`Cluster.Destroy called without Parent to unlink from`)
	}

	// call before unlink since it requires tec.Parent.*
	tec.actionDelete()
	tec.deletePropertyAllLocal()
	tec.deletePropertyAllInherited()
	// TODO delete all checks + check instances
	// TODO delete all inherited checks + check instances

	wg := new(sync.WaitGroup)
	for child, _ := range tec.Children {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			tec.Children[c].Destroy()
		}(child)
	}
	wg.Wait()

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(Builder).GetType(),
		ParentId:   tec.Parent.(Builder).GetID(),
		ParentName: tec.Parent.(Builder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)

	tec.setFault(nil)
	tec.setAction(nil)
}

func (tec *Cluster) Detach() {
	if tec.Parent == nil {
		panic(`Cluster.Detach called without Parent to detach from`)
	}
	bucket := tec.Parent.(Bucketeer).GetBucket()

	tec.deletePropertyAllInherited()
	// TODO delete all inherited checks + check instances

	tec.Parent.Unlink(UnlinkRequest{
		ParentType: tec.Parent.(Builder).GetType(),
		ParentId:   tec.Parent.(Builder).GetID(),
		ParentName: tec.Parent.(Builder).GetName(),
		ChildType:  tec.GetType(),
		ChildName:  tec.GetName(),
		ChildId:    tec.GetID(),
	},
	)

	bucket.Receive(ReceiveRequest{
		ParentType: bucket.(Builder).GetType(),
		ParentId:   bucket.(Builder).GetID(),
		ParentName: bucket.(Builder).GetName(),
		ChildType:  tec.Type,
		Cluster:    tec,
	},
	)

	tec.actionUpdate()
	tec.Parent.(Propertier).syncProperty(tec.Id.String())
}

//
// Interface: BucketAttacher
func (tec *Cluster) attachToBucket(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tec.Type,
		Cluster:    tec,
	})

	if tec.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_cluster`})
		return
	}
	tec.actionCreate()
}

//
// Interface: GroupAttacher
func (tec *Cluster) attachToGroup(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tec.Type,
		Cluster:    tec,
	})

	if tec.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_cluster`})
		return
	}
	tec.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
