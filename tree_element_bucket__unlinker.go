package somatree


//
// Interface: SomaTreeUnlinker
func (teb *SomaTreeElemBucket) Unlink(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "group":
			teb.unlinkGroup(u)
		case "cluster":
			teb.unlinkCluster(u)
		case "node":
			teb.unlinkNode(u)
		default:
			panic(`SomaTreeElemBucket.Unlink`)
		}
		return
	}
loop:
	for child, _ := range teb.Children {
		if teb.Children[child].(SomaTreeBuilder).GetType() == "node" {
			continue loop
		}
		teb.Children[child].(SomaTreeUnlinker).Unlink(u)
	}
}

//
// Interface: SomaTreeGroupUnlinker
func (teb *SomaTreeElemBucket) unlinkGroup(u UnlinkRequest) {
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
			panic(`SomaTreeElemBucket.unlinkGroup`)
		}
		return
	}
	panic(`SomaTreeElemBucket.unlinkGroup`)
}

//
// Interface: SomaTreeClusterUnlinker
func (teb *SomaTreeElemBucket) unlinkCluster(u UnlinkRequest) {
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
			panic(`SomaTreeElemBucket.unlinkCluster`)
		}
		return
	}
	panic(`SomaTreeElemBucket.unlinkCluster`)
}

//
// Interface: SomaTreeNodeUnlinker
func (teb *SomaTreeElemBucket) unlinkNode(u UnlinkRequest) {
	if unlinkRequestCheck(u, teb) {
		switch u.ChildType {
		case "node":
			if _, ok := teb.Children[u.ChildId]; ok {
				if u.ChildName == teb.Children[u.ChildId].GetName() {
					teb.Children[u.ChildId].clearParent()
					delete(teb.Children, u.ChildId)

					teb.Action <- &Action{
						Action:    "node_removal",
						Type:      "bucket",
						ChildType: "node",
						Bucket: somaproto.ProtoBucket{
							Id:          teb.Id.String(),
							Name:        teb.Name,
							Repository:  teb.Repository.String(),
							Team:        teb.Team.String(),
							Environment: teb.Environment,
							IsDeleted:   teb.Deleted,
							IsFrozen:    teb.Frozen,
						},
						Node: somaproto.ProtoNode{
							Id: u.ChildId,
						},
					}
				}
			}
		default:
			panic(`SomaTreeElemBucket.unlinkNode`)
		}
		return
	}
	panic(`SomaTreeElemBucket.unlinkNode`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
