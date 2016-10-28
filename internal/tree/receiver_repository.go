/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

// Interface: Receiver
func (ter *Repository) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "bucket":
			ter.receiveBucket(r)
		case "fault":
			ter.receiveFault(r)
		default:
			panic(`Repository.Receive`)
		}
		return
	}
	for child, _ := range ter.Children {
		ter.Children[child].(Receiver).Receive(r)
	}
}

// Interface: BucketReceiver
func (ter *Repository) receiveBucket(r ReceiveRequest) {
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "bucket":
			ter.Children[r.Bucket.GetID()] = r.Bucket
			r.Bucket.setParent(ter)
			r.Bucket.setAction(ter.Action)
			r.Bucket.setFault(ter.Fault)
			r.Bucket.setLoggerDeep(ter.log)
			ter.ordChildrenBck[ter.ordNumChildBck] = r.Bucket.GetID()
			ter.ordNumChildBck++
		default:
			panic(`Repository.receiveBucket`)
		}
	}
}

// Interface: FaultReceiver
func (ter *Repository) receiveFault(r ReceiveRequest) {
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "fault":
			ter.setFault(r.Fault)
			ter.Fault.setParent(ter)
			ter.updateFaultRecursive(ter.Fault)
		default:
			panic(`Repository.receiveFault`)
		}
		return
	}
	panic(`Repository.receiveFault`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
