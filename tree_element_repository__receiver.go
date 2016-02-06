package somatree

// Interface: SomaTreeReceiver
func (ter *SomaTreeElemRepository) Receive(r ReceiveRequest) {
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "bucket":
			ter.receiveBucket(r)
		case "fault":
			ter.receiveFault(r)
		default:
			panic(`SomaTreeElemRepository.Receive`)
		}
		return
	}
	for child, _ := range ter.Children {
		ter.Children[child].(SomaTreeReceiver).Receive(r)
	}
}

// Interface: SomaTreeBucketReceiver
func (ter *SomaTreeElemRepository) receiveBucket(r ReceiveRequest) {
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "bucket":
			ter.Children[r.Bucket.GetID()] = r.Bucket
			r.Bucket.setParent(ter)
			r.Bucket.setAction(ter.Action)
		default:
			panic(`SomaTreeElemRepository.receiveBucket`)
		}
	}
}

// Interface: SomaTreeFaultReceiver
func (ter *SomaTreeElemRepository) receiveFault(r ReceiveRequest) {
	if receiveRequestCheck(r, ter) {
		switch r.ChildType {
		case "fault":
			ter.setFault(r.Fault)
			ter.Fault.setParent(ter)
			ter.updateFaultRecursive(ter.Fault)
		default:
			panic(`SomaTreeElemRepository.receiveFault`)
		}
		return
	}
	panic(`SomaTreeElemRepository.receiveFault`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
