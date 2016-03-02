package somatree

//
// Interface: SomaTreeReceiver
func (teg *SomaTreeElemGroup) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, teg) {
		switch r.ChildType {
		case "group":
			teg.receiveGroup(r)
		case "cluster":
			teg.receiveCluster(r)
		case "node":
			teg.receiveNode(r)
		default:
			panic(`SomaTreeElemGroup.Receive`)
		}
		return
	}
loop:
	for child, _ := range teg.Children {
		if teg.Children[child].(Builder).GetType() == "node" {
			continue loop
		}
		teg.Children[child].(SomaTreeReceiver).Receive(r)
	}
}

//
// Interface: SomaTreeGroupReceiver
func (teg *SomaTreeElemGroup) receiveGroup(r ReceiveRequest) {
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
			panic(`SomaTreeElemGroup.receiveGroup`)
		}
		return
	}
	panic(`SomaTreeElemGroup.receiveGroup`)
}

//
// Interface: SomaTreeClusterReceiver
func (teg *SomaTreeElemGroup) receiveCluster(r ReceiveRequest) {
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
			panic(`SomaTreeElemGroup.receiveCluster`)
		}
		return
	}
	panic(`SomaTreeElemGroup.receiveCluster`)
}

//
// Interface: SomaTreeNodeReceiver
func (teg *SomaTreeElemGroup) receiveNode(r ReceiveRequest) {
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
			panic(`SomaTreeElemGroup.receiveNode`)
		}
		return
	}
	panic(`SomaTreeElemGroup.receiveNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
