package tree

import "sync"

//
// Interface: Attacher
func (ter *Repository) Attach(a AttachRequest) {
	if ter.Parent != nil {
		panic(`Repository.Attach: already attached`)
	}
	switch {
	case a.ParentType == "root":
		ter.attachToRoot(a)
	default:
		panic(`Repository.Attach`)
	}

	if ter.Parent == nil {
		panic(`Repository.Attach: failed`)
	}
	// no need to sync properties, as top level element the repo can't
	// inherit
}

func (ter *Repository) Destroy() {
	if ter.Parent == nil {
		panic(`Repository.Destroy called without Parent to unlink from`)
	}
	// call before unlink since it requires tec.Parent.*
	ter.actionDelete()
	ter.deletePropertyAllLocal()
	ter.deletePropertyAllInherited()
	// TODO delete all checks + check instances
	// TODO delete all inherited checks + check instances

	wg := new(sync.WaitGroup)
	for child, _ := range ter.Children {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			ter.Children[c].Destroy()
		}(child)
	}
	wg.Wait()

	// the Destroy handler of Fault calls
	// updateFaultRecursive(nil) on us
	ter.Fault.Destroy()

	ter.Parent.Unlink(UnlinkRequest{
		ParentType: ter.Parent.(Builder).GetType(),
		ParentId:   ter.Parent.(Builder).GetID(),
		ParentName: ter.Parent.(Builder).GetName(),
		ChildType:  ter.GetType(),
		ChildName:  ter.GetName(),
		ChildId:    ter.GetID(),
	},
	)

	ter.setAction(nil)
}

func (ter *Repository) Detach() {
	ter.Destroy()
}

// Interface: RootAttacher
func (ter *Repository) attachToRoot(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  "repository",
		Repository: ter,
	})

	ter.actionCreate()
	ter.Fault.setAction(ter.Action)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
