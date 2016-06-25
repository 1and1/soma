package somatree

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

	if teb.Parent == nil {
		panic(`Bucket.Attach: failed`)
	}
	teb.Parent.(Propertier).syncProperty(teb.Id.String())
}

func (teb *Bucket) Destroy() {
	if teb.Parent == nil {
		panic(`Bucket.Destroy called without Parent to unlink from`)
	}
	// XXX: destroy all inherited properties before unlinking
	// teb.(SomaTreePropertier).destroyInheritedProperties()

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

func (teb *Bucket) updateParentRecursive(p SomaTreeReceiver) {
	teb.setParent(p)
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func(str SomaTreeReceiver) {
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

func (teb *Bucket) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeBucketReceiver:
		teb.setBucketParent(p.(SomaTreeBucketReceiver))
		teb.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`Bucket.setParent`)
	}
}

func (teb *Bucket) setBucketParent(p SomaTreeBucketReceiver) {
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
