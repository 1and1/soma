package tree

//
// Interface: Receiver
func (tec *Cluster) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, tec) {
		switch r.ChildType {
		case "node":
			tec.receiveNode(r)
		default:
			panic(`Cluster.Receive`)
		}
	}
	// no passing along since only nodes are a SomeTreeClusterAttacher
	// and nodes can have no children
	return
}

//
// Interface: NodeReceiver
func (tec *Cluster) receiveNode(r ReceiveRequest) {
	if receiveRequestCheck(r, tec) {
		switch r.ChildType {
		case "node":
			tec.Children[r.Node.GetID()] = r.Node
			r.Node.setParent(tec)
			r.Node.setAction(tec.Action)
			r.Node.setFault(tec.Fault)

			tec.actionMemberNew(Action{
				ChildType: "node",
				ChildNode: r.Node.export(),
			})
		default:
			panic(`Cluster.receiveNode`)
		}
		return
	}
	panic(`Cluster.receiveNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
