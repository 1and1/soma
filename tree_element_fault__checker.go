package somatree

//
// Interface: Checker
func (tef *SomaTreeElemFault) SetCheck(c Check) {
}

func (tef *SomaTreeElemFault) inheritCheck(c Check) {
}

func (tef *SomaTreeElemFault) inheritCheckDeep(c Check) {
}

func (tef *SomaTreeElemFault) storeCheck(c Check) {
}

func (tef *SomaTreeElemFault) syncCheck(childId string) {
}

func (tef *SomaTreeElemFault) checkCheck(checkId string) bool {
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
