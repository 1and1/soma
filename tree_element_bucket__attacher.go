package somatree

import (
	"fmt"
	"reflect"
	"sync"
)

//
// Interface: SomaTreeAttacher
func (teb *SomaTreeElemBucket) Attach(a AttachRequest) {
	if teb.Parent != nil {
		panic(`SomaTreeElemBucket.Attach: already attached`)
	}
	switch {
	case a.ParentType == "repository":
		teb.attachToRepository(a)
	}
}

func (teb *SomaTreeElemBucket) Destroy() {
	if teb.Parent == nil {
		panic(`SomaTreeElemBucket.Destroy called without Parent to unlink from`)
	}

	teb.Parent.Unlink(UnlinkRequest{
		ParentType: teb.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   teb.Parent.(SomaTreeBuilder).GetID(),
		ParentName: teb.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  teb.GetType(),
		ChildName:  teb.GetName(),
		ChildId:    teb.GetID(),
	},
	)

	teb.setFault(nil)

	teb.Action <- &Action{
		Action: "delete",
		Type:   "bucket",
		Id:     teb.Id.String(),
		Name:   teb.Name,
	}
	teb.setAction(nil)
}

func (teb *SomaTreeElemBucket) Detach() {
	teb.Destroy()
}

func (teb *SomaTreeElemBucket) clearParent() {
	teb.Parent = nil
	teb.State = "floating"
}

func (teb *SomaTreeElemBucket) setFault(f *SomaTreeElemFault) {
	teb.Fault = f
}

func (teb *SomaTreeElemBucket) updateParentRecursive(p SomaTreeReceiver) {
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

func (teb *SomaTreeElemBucket) updateFaultRecursive(f *SomaTreeElemFault) {
	teb.setFault(f)
	var wg sync.WaitGroup
	for child, _ := range teb.Children {
		wg.Add(1)
		c := child
		go func(ptr *SomaTreeElemFault) {
			defer wg.Done()
			teb.Children[c].updateFaultRecursive(ptr)
		}(f)
	}
	wg.Wait()
}

func (teb *SomaTreeElemBucket) setParent(p SomaTreeReceiver) {
	switch p.(type) {
	case SomaTreeBucketReceiver:
		teb.setBucketParent(p.(SomaTreeBucketReceiver))
		teb.State = "attached"
	default:
		fmt.Printf("Type: %s\n", reflect.TypeOf(p))
		panic(`SomaTreeElemBucket.setParent`)
	}
}

func (teb *SomaTreeElemBucket) setBucketParent(p SomaTreeBucketReceiver) {
	teb.Parent = p
}

//
// Interface: SomaTreeRepositoryAttacher
func (teb *SomaTreeElemBucket) attachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  teb.Type,
		Bucket:     teb,
	})

	teb.Action <- &Action{
		Action: "create",
		Type:   "bucket",
		Id:     teb.Id.String(),
		Name:   teb.Name,
	}
}
