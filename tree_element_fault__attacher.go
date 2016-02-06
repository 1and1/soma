package somatree

//
// Interface: SomaTreeAttacher
func (tef *SomaTreeElemFault) Attach(a AttachRequest) {
	if tef.Parent != nil {
		panic(`SomaTreeElemFault.Attach: already attached`)
	}
	switch {
	case a.ParentType == "repository":
		tef.attachToRepository(a)
	}
}

func (tef *SomaTreeElemFault) Destroy() {
	if tef.Parent == nil {
		panic(`SomaTreeElemFault.Destroy called without Parent to unlink from`)
	}

	tef.Parent.(SomaTreeAttacher).updateFaultRecursive(nil)

	tef.Parent.Unlink(UnlinkRequest{
		ParentType: tef.Parent.(SomaTreeBuilder).GetType(),
		ParentId:   tef.Parent.(SomaTreeBuilder).GetID(),
		ParentName: tef.Parent.(SomaTreeBuilder).GetName(),
		ChildType:  tef.GetType(),
		ChildName:  tef.GetName(),
		ChildId:    tef.GetID(),
	},
	)

	tef.setAction(nil)
	tef.Error = nil
}

func (tef *SomaTreeElemFault) Detach() {
	tef.Destroy()
}

//
// Interface: SomaTreeRepositoryAttacher
func (tef *SomaTreeElemFault) attachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tef.Type,
		Fault:      tef,
	})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
