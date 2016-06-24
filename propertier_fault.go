package somatree

func (tef *SomaTreeElemFault) SetProperty(
	p Property) {
}

func (tef *SomaTreeElemFault) inheritProperty(
	p Property) {
}

func (tef *SomaTreeElemFault) inheritPropertyDeep(
	p Property) {
}

func (tef *SomaTreeElemFault) setCustomProperty(
	p Property) {
}

func (tef *SomaTreeElemFault) setServiceProperty(
	p Property) {
}

func (tef *SomaTreeElemFault) setSystemProperty(
	p Property) {
}

func (tef *SomaTreeElemFault) setOncallProperty(
	p Property) {
}

func (tef *SomaTreeElemFault) syncProperty(
	childId string) {
}

func (tef *SomaTreeElemFault) checkProperty(
	propType string, propId string) bool {
	return false
}

func (tef *SomaTreeElemFault) checkDuplicate(p Property) (
	bool, bool, Property) {
	return true, false, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
