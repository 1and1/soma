package somatree

//
// Interface: SomaTreeUnlinker
func (tec *SomaTreeElemCluster) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, tec) {
		switch u.ChildType {
		case "node":
			tec.unlinkNode(u)
		default:
			panic(`SomaTreeElemCluster.Unlink`)
		}
	}
	// no passing along since only nodes are a SomeTreeClusterAttacher
	// and nodes can have no children
	return
}

//
// Interface: SomaTreeNodeUnlinker
func (tec *SomaTreeElemCluster) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, tec) {
		switch u.ChildType {
		case "node":
			if _, ok := tec.Children[u.ChildId]; ok {
				if u.ChildName == tec.Children[u.ChildId].GetName() {
					tec.Children[u.ChildId].clearParent()
					delete(tec.Children, u.ChildId)

					tec.Action <- &Action{
						Action:    "member_removed",
						Type:      "cluster",
						Id:        tec.Id.String(),
						Name:      tec.Name,
						Team:      tec.Team.String(),
						ChildType: "node",
						ChildId:   u.ChildId,
					}
				}
			}
		default:
			panic(`SomaTreeElemCluster.unlinkNode`)
		}
		return
	}
	panic(`SomaTreeElemCluster.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
