package somatree

func (tef *SomaTreeElemFault) SetProperty(p Property) {
}

func (tef *SomaTreeElemFault) setPropertyInherited(p Property) {
}

func (tef *SomaTreeElemFault) setPropertyOnChildren(p Property) {
}

func (tef *SomaTreeElemFault) addProperty(p Property) {
}

func (tef *SomaTreeElemFault) UpdateProperty(p Property) {
}

func (tef *SomaTreeElemFault) updatePropertyInherited(p Property) {
}

func (tef *SomaTreeElemFault) updatePropertyOnChildren(p Property) {
}

func (tef *SomaTreeElemFault) switchProperty(p Property) {
}

func (tef *SomaTreeElemFault) DeleteProperty(p Property) {
}

func (tef *SomaTreeElemFault) deletePropertyInherited(p Property) {
}

func (tef *SomaTreeElemFault) deletePropertyOnChildren(p Property) {
}

func (tef *SomaTreeElemFault) rmProperty(p Property) {
}

func (tef *SomaTreeElemFault) verifySourceInstance(id, prop string) bool {
	return false
}

func (tef *SomaTreeElemFault) findIdForSource(source, prop string) string {
	return ``
}

func (tef *SomaTreeElemFault) syncProperty(childId string) {
}

func (tef *SomaTreeElemFault) checkProperty(propType, propId string) bool {
	return false
}

func (tef *SomaTreeElemFault) checkDuplicate(p Property) (bool, bool, Property) {
	return true, false, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
