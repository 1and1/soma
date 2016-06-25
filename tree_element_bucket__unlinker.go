package somatree

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
