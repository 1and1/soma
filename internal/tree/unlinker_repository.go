/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

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
					for i, bck := range ter.ordChildrenBck {
						if bck == u.ChildId {
							delete(ter.ordChildrenBck, i)
						}
					}
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
