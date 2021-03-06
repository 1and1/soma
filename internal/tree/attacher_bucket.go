/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"fmt"
	"reflect"
	"sync"
)

//
// Interface: Attacher
func (teb *Bucket) Attach(a AttachRequest) {
	if teb.Parent != nil {
		panic(`Bucket.Attach: already attached`)
	}
	switch {
	case a.ParentType == "repository":
		teb.attachToRepository(a)
	default:
		panic(`Bucket.Attach`)
	}

	teb.Parent.(Propertier).syncProperty(teb.Id.String())
	teb.Parent.(Checker).syncCheck(teb.Id.String())
}

func (teb *Bucket) Destroy() {
	if teb.Parent == nil {
		panic(`Bucket.Destroy called without Parent to unlink from`)
	}
	teb.deletePropertyAllLocal()
	teb.deletePropertyAllInherited()
	// TODO delete all checks + check instances
	// TODO delete all inherited checks + check instances

	wg := new(sync.WaitGroup)
	for child, _ := range teb.Children {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			teb.Children[c].Destroy()
		}(child)
	}
	wg.Wait()

	teb.Parent.Unlink(UnlinkRequest{
		ParentType: teb.Parent.(Builder).GetType(),
		ParentId:   teb.Parent.(Builder).GetID(),
		ParentName: teb.Parent.(Builder).GetName(),
		ChildType:  teb.GetType(),
		ChildName:  teb.GetName(),
		ChildId:    teb.GetID(),
	},
	)

	teb.setFault(nil)
	teb.actionDelete()
	teb.setAction(nil)
}

func (teb *Bucket) Detach() {
	teb.Destroy()
}

func (teb *Bucket) clearParent() {
	teb.Parent = nil
	teb.State = "floating"
}

func (teb *Bucket) setFault(f *Fault) {
	teb.Fault = f
}

func (teb *Bucket) updateParentRecursive(p Receiver) {
	teb.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func(str Receiver) {
			defer wg.Done()
			teb.Children[c].updateParentRecursive(str)
		}(teb)
	}
	wg.Wait()
}

func (teb *Bucket) updateFaultRecursive(f *Fault) {
	teb.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func(ptr *Fault) {
			defer wg.Done()
			teb.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

func (teb *Bucket) setParent(p Receiver) {
	switch p.(type) {
	case BucketReceiver:
		teb.setBucketParent(p.(BucketReceiver))
		teb.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`Bucket.setParent`)
	}
}

func (teb *Bucket) setBucketParent(p BucketReceiver) {
	teb.Parent = p
}

//
// Interface: RepositoryAttacher
func (teb *Bucket) attachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teb.Type,
		Bucket:     teb,
	})

	if teb.Parent == nil {
		a.Root.(*Tree).AttachError(Error{Action: `attach_bucket`})
		return
	}
	teb.actionCreate()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
