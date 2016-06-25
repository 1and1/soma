package tree

//
// Interface: Receiver
func (teg *Group) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "group":
			teg.receiveGroup(r)
		case "cluster":
			teg.receiveCluster(r)
		case "node":
			teg.receiveNode(r)
		default:
			panic(`Group.Receive`)
		}
		return
	}
loop:
	for child, _ := range teg.Children {
		if teg.Children[child].(Builder).GetType() == "node" {
			continue loop
		}
		teg.Children[child].(Receiver).Receive(r)
	}
}

//
// Interface: GroupReceiver
func (teg *Group) receiveGroup(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "group":
			teg.Children[r.Group.GetID()] = r.Group
			r.Group.setParent(teg)
			r.Group.setAction(teg.Action)
			r.Group.setFault(teg.Fault)

			teg.actionMemberNew(Action{
				ChildType:  "group",
				ChildGroup: r.Group.export(),
			})
		default:
			panic(`Group.receiveGroup`)
		}
		return
	}
	panic(`Group.receiveGroup`)
}

//
// Interface: ClusterReceiver
func (teg *Group) receiveCluster(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "cluster":
			teg.Children[r.Cluster.GetID()] = r.Cluster
			r.Cluster.setParent(teg)
			r.Cluster.setAction(teg.Action)
			r.Cluster.setFault(teg.Fault)

			teg.actionMemberNew(Action{
				ChildType:    "cluster",
				ChildCluster: r.Cluster.export(),
			})
		default:
			panic(`Group.receiveCluster`)
		}
		return
	}
	panic(`Group.receiveCluster`)
}

//
// Interface: NodeReceiver
func (teg *Group) receiveNode(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "node":
			teg.Children[r.Node.GetID()] = r.Node
			r.Node.setParent(teg)
			r.Node.setAction(teg.Action)
			r.Node.setFault(teg.Fault)

			teg.actionMemberNew(Action{
				ChildType: "node",
				ChildNode: r.Node.export(),
			})
		default:
			panic(`Group.receiveNode`)
		}
		return
	}
	panic(`Group.receiveNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
