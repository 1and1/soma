package somatree

//
// Interface: SomaTreeAttacher
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
	// XXX: destroy all properties before unlinking
	// ter.(SomaTreePropertier).nukeAllProperties()

	// the Destroy handler of SomaTreeElemFault calls
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

	ter.actionDelete()
	ter.setAction(nil)
	ter.PropertyOncall = nil
	ter.PropertyService = nil
	ter.PropertySystem = nil
	ter.PropertyCustom = nil
}

func (ter *Repository) Detach() {
	ter.Destroy()
}

// Interface: SomaTreeRootAttacher
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
