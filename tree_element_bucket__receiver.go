package somatree


//
// Interface: SomaTreeReceiver
func (teb *Bucket) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "group":
			teb.receiveGroup(r)
		case "cluster":
			teb.receiveCluster(r)
		case "node":
			teb.receiveNode(r)
		default:
			panic(`Bucket.Receive`)
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
func (teb *Bucket) receiveGroup(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "group":
			teb.Children[r.Group.GetID()] = r.Group
			r.Group.setParent(teb)
			r.Group.setAction(teb.Action)
			r.Group.setFault(teb.Fault)
		default:
			panic(`Bucket.receiveGroup`)
		}
		return
	}
	panic(`Bucket.receiveGroup`)
}

//
// Interface: SomaTreeClusterReceiver
func (teb *Bucket) receiveCluster(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "cluster":
			teb.Children[r.Cluster.GetID()] = r.Cluster
			r.Cluster.setParent(teb)
			r.Cluster.setAction(teb.Action)
			r.Cluster.setFault(teb.Fault)
		default:
			panic(`Bucket.receiveCluster`)
		}
		return
	}
	panic(`Bucket.receiveCluster`)
}

//
// Interface: SomaTreeNodeReceiver
func (teb *Bucket) receiveNode(r ReceiveRequest) {
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
			panic(`Bucket.receiveNode`)
		}
		return
	}
	panic(`Bucket.receiveNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
