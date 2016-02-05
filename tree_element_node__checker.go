package somatree

//
// Interface: SomaTreeChecker
func (ten *SomaTreeElemNode) SetCheck(c SomaTreeCheck) {
	c.InheritedFrom = ten.Id
	c.Inherited = false
	ten.storeCheck(c)
}

func (ten *SomaTreeElemNode) inheritCheck(c SomaTreeCheck) {
	ten.storeCheck(c)
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) inheritCheckDeep(c SomaTreeCheck) {
}

func (ten *SomaTreeElemNode) storeCheck(c SomaTreeCheck) {
	ten.Checks[c.Id.String()] = c
}

// noop, satisfy interface
func (ten *SomaTreeElemNode) syncCheck(childId string) {
}

func (ten *SomaTreeElemNode) checkCheck(checkId string) bool {
	if _, ok := ten.Checks[checkId]; ok {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
