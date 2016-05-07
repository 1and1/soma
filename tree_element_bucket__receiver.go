package somatree


//
// Interface: SomaTreeReceiver
func (teb *SomaTreeElemBucket) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "group":
			teb.receiveGroup(r)
		case "cluster":
			teb.receiveCluster(r)
		case "node":
			teb.receiveNode(r)
		default:
			panic(`SomaTreeElemBucket.Receive`)
		}
		return
	}
loop:
	for child, _ := range teb.Children {
		if teb.Children[child].(Builder).GetType() == "node" {
			continue loop
		}
		teb.Children[child].(SomaTreeReceiver).Receive(r)
	}
}

//
// Interface: SomaTreeGroupReceiver
func (teb *SomaTreeElemBucket) receiveGroup(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "group":
			teb.Children[r.Group.GetID()] = r.Group
			r.Group.setParent(teb)
			r.Group.setAction(teb.Action)
			r.Group.setFault(teb.Fault)
		default:
			panic(`SomaTreeElemBucket.receiveGroup`)
		}
		return
	}
	panic(`SomaTreeElemBucket.receiveGroup`)
}

//
// Interface: SomaTreeClusterReceiver
func (teb *SomaTreeElemBucket) receiveCluster(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "cluster":
			teb.Children[r.Cluster.GetID()] = r.Cluster
			r.Cluster.setParent(teb)
			r.Cluster.setAction(teb.Action)
			r.Cluster.setFault(teb.Fault)
		default:
			panic(`SomaTreeElemBucket.receiveCluster`)
		}
		return
	}
	panic(`SomaTreeElemBucket.receiveCluster`)
}

//
// Interface: SomaTreeNodeReceiver
func (teb *SomaTreeElemBucket) receiveNode(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "node":
			teb.Children[r.Node.GetID()] = r.Node
			r.Node.setParent(teb)
			r.Node.setAction(teb.Action)
			r.Node.setFault(teb.Fault)

			teb.actionAssignNode(Action{
				ChildType: "node",
				ChildNode: proto.Node{
					Id: r.Node.GetID(),
				},
			})
		default:
			panic(`SomaTreeElemBucket.receiveNode`)
		}
		return
	}
	panic(`SomaTreeElemBucket.receiveNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
