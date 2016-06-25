package tree

//
// Interface: Unlinker
func (teg *Group) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "group":
			teg.unlinkGroup(u)
		case "cluster":
			teg.unlinkCluster(u)
		case "node":
			teg.unlinkNode(u)
		default:
			panic(`Group.Unlink`)
		}
		return
	}
loop:
	for child, _ := range teg.Children {
		if teg.Children[child].(Builder).GetType() == "node" {
			continue loop
		}
		teg.Children[child].(Unlinker).Unlink(u)
	}
}

//
// Interface: GroupUnlinker
func (teg *Group) unlinkGroup(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "group":
			if _, ok := teg.Children[u.ChildId]; ok {
				if u.ChildName == teg.Children[u.ChildId].GetName() {
					a := Action{
						ChildType:  "group",
						ChildGroup: teg.Children[u.ChildId].(*Group).export(),
					}

					teg.Children[u.ChildId].clearParent()
					delete(teg.Children, u.ChildId)

					teg.actionMemberRemoved(a)
				}
			}
		default:
			panic(`Group.unlinkGroup`)
		}
		return
	}
	panic(`Group.unlinkGroup`)
}

//
// Interface: ClusterUnlinker
func (teg *Group) unlinkCluster(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "cluster":
			if _, ok := teg.Children[u.ChildId]; ok {
				if u.ChildName == teg.Children[u.ChildId].GetName() {
					a := Action{
						ChildType:    "cluster",
						ChildCluster: teg.Children[u.ChildId].(*Cluster).export(),
					}

					teg.Children[u.ChildId].clearParent()
					delete(teg.Children, u.ChildId)

					teg.actionMemberRemoved(a)
				}
			}
		default:
			panic(`Group.unlinkCluster`)
		}
		return
	}
	panic(`Group.unlinkCluster`)
}

//
// Interface: NodeUnlinker
func (teg *Group) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, teg) {
		switch u.ChildType {
		case "node":
			if _, ok := teg.Children[u.ChildId]; ok {
				if u.ChildName == teg.Children[u.ChildId].GetName() {
					a := Action{
						ChildType: "node",
						ChildNode: teg.Children[u.ChildId].(*Node).export(),
					}

					teg.Children[u.ChildId].clearParent()
					delete(teg.Children, u.ChildId)

					teg.actionMemberRemoved(a)
				}
			}
		default:
			panic(`Group.unlinkNode`)
		}
		return
	}
	panic(`Group.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
