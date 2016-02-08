package somatree


//
// Interface: SomaTreeAttacher
func (ter *SomaTreeElemRepository) Attach(a AttachRequest) {
	if ter.Parent != nil {
		panic(`SomaTreeElemRepository.Attach: already attached`)
	}
	switch {
	case a.ParentType == "root" &&
		a.ChildType == "repository" &&
		a.ChildName == ter.Name:
		ter.attachToRoot(a)
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
		ParentType: ter.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   ter.Parent.(SomaTreeBuilder).GetID(),
		ParentName: ter.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  ter.GetType(),
		ChildName:  ter.GetName(),
		ChildId:    ter.GetID(),
	},
	)

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

	ter.Action <- &Action{
		Action: "create",
		Type:   "repository",
		Repository: somaproto.ProtoRepository{
			Id:        ter.Id.String(),
			Name:      ter.Name,
			Team:      ter.Team.String(),
			IsDeleted: ter.Deleted,
			IsActive:  ter.Active,
		},
	}
	ter.Fault.setAction(ter.Action)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
