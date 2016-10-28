/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import "github.com/1and1/soma/lib/proto"

//
// Interface: Receiver
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
		teb.Children[child].(Receiver).Receive(r)
	}
}

//
// Interface: GroupReceiver
func (teb *Bucket) receiveGroup(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "group":
			teb.Children[r.Group.GetID()] = r.Group
			r.Group.setParent(teb)
			r.Group.setAction(teb.Action)
			r.Group.setFault(teb.Fault)
			r.Group.setLoggerDeep(teb.log)
			teb.ordChildrenGrp[teb.ordNumChildGrp] = r.Group.GetID()
			teb.ordNumChildGrp++
		default:
			panic(`Bucket.receiveGroup`)
		}
		return
	}
	panic(`Bucket.receiveGroup`)
}

//
// Interface: ClusterReceiver
func (teb *Bucket) receiveCluster(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "cluster":
			teb.Children[r.Cluster.GetID()] = r.Cluster
			r.Cluster.setParent(teb)
			r.Cluster.setAction(teb.Action)
			r.Cluster.setFault(teb.Fault)
			r.Cluster.setLoggerDeep(teb.log)
			teb.ordChildrenClr[teb.ordNumChildClr] = r.Cluster.GetID()
			teb.ordNumChildClr++
		default:
			panic(`Bucket.receiveCluster`)
		}
		return
	}
	panic(`Bucket.receiveCluster`)
}

//
// Interface: NodeReceiver
func (teb *Bucket) receiveNode(r ReceiveRequest) {
	if receiveRequestCheck(r, teb) {
		switch r.ChildType {
		case "node":
			teb.Children[r.Node.GetID()] = r.Node
			r.Node.setParent(teb)
			r.Node.setAction(teb.Action)
			r.Node.setFault(teb.Fault)
			r.Node.setLoggerDeep(teb.log)
			teb.ordChildrenNod[teb.ordNumChildNod] = r.Node.GetID()
			teb.ordNumChildNod++

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
