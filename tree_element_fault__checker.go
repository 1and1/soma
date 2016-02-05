package somatree

//
// Interface: SomaTreeChecker
func (tef *SomaTreeElemFault) SetCheck(c SomaTreeCheck) {
}

func (tef *SomaTreeElemFault) inheritCheck(c SomaTreeCheck) {
}

func (tef *SomaTreeElemFault) inheritCheckDeep(c SomaTreeCheck) {
}

func (tef *SomaTreeElemFault) storeCheck(c SomaTreeCheck) {
}

func (tef *SomaTreeElemFault) syncCheck(childId string) {
}

func (tef *SomaTreeElemFault) checkCheck(checkId string) bool {
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
