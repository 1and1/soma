package somatree

// Interface: SomaTreeUnlinker
func (ter *SomaTreeElemRepository) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "bucket":
			ter.unlinkBucket(u)
		case "fault":
			ter.unlinkFault(u)
		default:
			panic(`SomaTreeElemRepository.Unlink`)
		}
		return
	}
	for child, _ := range ter.Children {
		ter.Children[child].(SomaTreeUnlinker).Unlink(u)
	}
}

// Interface: SomaTreeBucketUnlinker
func (ter *SomaTreeElemRepository) unlinkBucket(u UnlinkRequest) {
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
			panic(`SomaTreeElemRepository.unlinkBucket`)
		}
		return
	}
	panic(`SomaTreeElemRepository.unlinkBucket`)
}

// Interface: SomaTreeFaultUnlinker
func (ter *SomaTreeElemRepository) unlinkFault(u UnlinkRequest) {
	if unlinkRequestCheck(u, ter) {
		switch u.ChildType {
		case "fault":
			ter.Fault = nil
			ter.updateFaultRecursive(ter.Fault)
		default:
			panic(`SomaTreeElemRepository.unlinkFault`)
		}
		return
	}
	panic(`SomaTreeElemRepository.unlinkFault`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
