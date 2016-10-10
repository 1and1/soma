/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

//
// Interface: Unlinker
func (teb *Bucket) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "group":
			teb.unlinkGroup(u)
		case "cluster":
			teb.unlinkCluster(u)
		case "node":
			teb.unlinkNode(u)
		default:
			panic(`Bucket.Unlink`)
		}
		return
	}
loop:
	for child, _ := range teb.Children {
		if teb.Children[child].(Builder).GetType() == "node" {
			continue loop
		}
		teb.Children[child].(Unlinker).Unlink(u)
	}
}

//
// Interface: GroupUnlinker
func (teb *Bucket) unlinkGroup(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "group":
			if _, ok := teb.Children[u.ChildId]; ok {
				if u.ChildName == teb.Children[u.ChildId].GetName() {
					teb.Children[u.ChildId].clearParent()
					delete(teb.Children, u.ChildId)
					for i, grp := range teb.ordChildrenGrp {
						if grp == u.ChildId {
							delete(teb.ordChildrenGrp, i)
						}
					}
				}
			}
		default:
			panic(`Bucket.unlinkGroup`)
		}
		return
	}
	panic(`Bucket.unlinkGroup`)
}

//
// Interface: ClusterUnlinker
func (teb *Bucket) unlinkCluster(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "cluster":
			if _, ok := teb.Children[u.ChildId]; ok {
				if u.ChildName == teb.Children[u.ChildId].GetName() {
					teb.Children[u.ChildId].clearParent()
					delete(teb.Children, u.ChildId)
					for i, clr := range teb.ordChildrenClr {
						if clr == u.ChildId {
							delete(teb.ordChildrenClr, i)
						}
					}
				}
			}
		default:
			panic(`Bucket.unlinkCluster`)
		}
		return
	}
	panic(`Bucket.unlinkCluster`)
}

//
// Interface: NodeUnlinker
func (teb *Bucket) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "node":
			if _, ok := teb.Children[u.ChildId]; ok {
				if u.ChildName == teb.Children[u.ChildId].GetName() {
					teb.Children[u.ChildId].clearParent()
					delete(teb.Children, u.ChildId)

					// no action here, the node itself will either
					// update its state from standalone->grouped|clustered
					// or delete the bucket_assignment on Destroy(),
					// which can not be differentiated here
					for i, nod := range teb.ordChildrenNod {
						if nod == u.ChildId {
							delete(teb.ordChildrenNod, i)
						}
					}
				}
			}
		default:
			panic(`Bucket.unlinkNode`)
		}
		return
	}
	panic(`Bucket.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
