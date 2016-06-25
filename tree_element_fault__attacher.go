package somatree

//
// Interface: SomaTreeAttacher
func (tef *Fault) Attach(a AttachRequest) {
	if tef.Parent != nil {
		panic(`Fault.Attach: already attached`)
	}
	switch {
	case a.ParentType == "repository":
		tef.attachToRepository(a)
	}
}

func (tef *Fault) Destroy() {
	if tef.Parent == nil {
		panic(`Fault.Destroy called without Parent to unlink from`)
	}

	tef.Parent.(SomaTreeAttacher).updateFaultRecursive(nil)

	tef.Parent.Unlink(UnlinkRequest{
		ParentType: tef.Parent.(Builder).GetType(),
		ParentId:   tef.Parent.(Builder).GetID(),
		ParentName: tef.Parent.(Builder).GetName(),
		ChildType:  tef.GetType(),
		ChildName:  tef.GetName(),
		ChildId:    tef.GetID(),
	},
	)

	tef.setAction(nil)
	tef.Error = nil
}

func (tef *Fault) Detach() {
	tef.Destroy()
}

func (tef *Fault) ReAttach(a AttachRequest) {
}

//
// Interface: SomaTreeRepositoryAttacher
func (tef *Fault) attachToRepository(a AttachRequest) {
	a.Root.Receive(ReceiveRequest{
		ParentType: a.ParentType,
		ParentId:   a.ParentId,
		ParentName: a.ParentName,
		ChildType:  tef.Type,
		Fault:      tef,
	})
}

//
// Interface: SomaTreeBucketAttacher
func (tef *Fault) attachToBucket(a AttachRequest) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
