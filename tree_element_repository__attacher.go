package somatree

//
// Interface: SomaTreeAttacher
func (ter *SomaTreeElemRepository) Attach(a AttachRequest) {
	if ter.Parent != nil {
		panic(`SomaTreeElemRepository.Attach: already attached`)
	}
	switch {
	case a.ParentType == "root":
		ter.attachToRoot(a)
	default:
		panic(`SomaTreeElemRepository.Attach`)
	}
}

func (ter *SomaTreeElemRepository) Destroy() {
	if ter.Parent == nil {
		panic(`SomaTreeElemRepository.Destroy called without Parent to unlink from`)
	}

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

func (ter *SomaTreeElemRepository) Detach() {
	ter.Destroy()
}

// Interface: SomaTreeRootAttacher
func (ter *SomaTreeElemRepository) attachToRoot(a AttachRequest) {
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
