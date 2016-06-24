package somatree

func (tef *SomaTreeElemFault) SetProperty(
	p SomaTreeProperty) {
}

func (tef *SomaTreeElemFault) inheritProperty(
	p SomaTreeProperty) {
}

func (tef *SomaTreeElemFault) inheritPropertyDeep(
	p SomaTreeProperty) {
}

func (tef *SomaTreeElemFault) setCustomProperty(
	p SomaTreeProperty) {
}

func (tef *SomaTreeElemFault) setServiceProperty(
	p SomaTreeProperty) {
}

func (tef *SomaTreeElemFault) setSystemProperty(
	p SomaTreeProperty) {
}

func (tef *SomaTreeElemFault) setOncallProperty(
	p SomaTreeProperty) {
}

func (tef *SomaTreeElemFault) syncProperty(
	childId string) {
}

func (tef *SomaTreeElemFault) checkProperty(
	propType string, propId string) bool {
	return false
}

func (tef *SomaTreeElemFault) checkDuplicate(p SomaTreeProperty) (
	bool, bool, SomaTreeProperty) {
	return true, false, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
