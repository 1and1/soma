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
func (tec *Cluster) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, tec) {
		switch u.ChildType {
		case "node":
			tec.unlinkNode(u)
		default:
			panic(`Cluster.Unlink`)
		}
	}
	// no passing along since only nodes are a SomeTreeClusterAttacher
	// and nodes can have no children
	return
}

//
// Interface: NodeUnlinker
func (tec *Cluster) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, tec) {
		switch u.ChildType {
		case "node":
			if _, ok := tec.Children[u.ChildId]; ok {
				if u.ChildName == tec.Children[u.ChildId].GetName() {
					a := Action{
						ChildType: "node",
						ChildNode: tec.Children[u.ChildId].(*Node).export(),
					}

					tec.Children[u.ChildId].clearParent()
					delete(tec.Children, u.ChildId)
					for i, nod := range tec.ordChildrenNod {
						if nod == u.ChildId {
							delete(tec.ordChildrenNod, i)
						}
					}

					tec.actionMemberRemoved(a)
				}
			}
		default:
			panic(`Cluster.unlinkNode`)
		}
		return
	}
	panic(`Cluster.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
