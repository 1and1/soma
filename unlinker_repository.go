package tree

// Interface: Unlinker
func (ter *Repository) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "bucket":
			ter.unlinkBucket(u)
		case "fault":
			ter.unlinkFault(u)
		default:
			panic(`Repository.Unlink`)
		}
		return
	}
	for child, _ := range ter.Children {
		ter.Children[child].(Unlinker).Unlink(u)
	}
}

// Interface: BucketUnlinker
func (ter *Repository) unlinkBucket(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "bucket":
			if _, ok := ter.Children[u.ChildId]; ok {
				if u.ChildName == ter.Children[u.ChildId].GetName() {
					ter.Children[u.ChildId].clearParent()
					delete(ter.Children, u.ChildId)
				}
			}
		default:
			panic(`Repository.unlinkBucket`)
		}
		return
	}
	panic(`Repository.unlinkBucket`)
}

// Interface: FaultUnlinker
func (ter *Repository) unlinkFault(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "fault":
			ter.Fault = nil
			ter.updateFaultRecursive(ter.Fault)
		default:
			panic(`Repository.unlinkFault`)
		}
		return
	}
	panic(`Repository.unlinkFault`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
