package somatree

//
// Interface: SomaTreeUnlinker
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
// Interface: SomaTreeNodeUnlinker
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
