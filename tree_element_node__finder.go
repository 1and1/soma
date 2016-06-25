package somatree

//
// Interface: SomaTreeFinder
func (ten *Node) Find(f FindRequest, b bool) Attacher {
	if findRequestCheck(f, ten) {
		return ten
	} else if b {
		return ten.Fault
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
