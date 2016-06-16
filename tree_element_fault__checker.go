package somatree

//
// Interface: Checker
func (tef *SomaTreeElemFault) SetCheck(c Check) {
}

func (tef *SomaTreeElemFault) setCheckInherited(c Check) {
}

func (tef *SomaTreeElemFault) setCheckOnChildren(c Check) {
}

func (tef *SomaTreeElemFault) addCheck(c Check) {
}

func (tef *SomaTreeElemFault) DeleteCheck(c Check) {
}

func (tef *SomaTreeElemFault) deleteCheckInherited(c Check) {
}

func (tef *SomaTreeElemFault) deleteCheckOnChildren(c Check) {
}

func (tef *SomaTreeElemFault) rmCheck(c Check) {
}

func (tef *SomaTreeElemFault) syncCheck(childId string) {
}

func (tef *SomaTreeElemFault) checkCheck(checkId string) bool {
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
