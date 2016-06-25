package tree

func (tef *Fault) SetProperty(p Property) {
}

func (tef *Fault) setPropertyInherited(p Property) {
}

func (tef *Fault) setPropertyOnChildren(p Property) {
}

func (tef *Fault) addProperty(p Property) {
}

func (tef *Fault) UpdateProperty(p Property) {
}

func (tef *Fault) updatePropertyInherited(p Property) {
}

func (tef *Fault) updatePropertyOnChildren(p Property) {
}

func (tef *Fault) switchProperty(p Property) {
}

func (tef *Fault) DeleteProperty(p Property) {
}

func (tef *Fault) deletePropertyInherited(p Property) {
}

func (tef *Fault) deletePropertyOnChildren(p Property) {
}

func (tef *Fault) rmProperty(p Property) {
}

func (tef *Fault) verifySourceInstance(id, prop string) bool {
	return false
}

func (tef *Fault) findIdForSource(source, prop string) string {
	return ``
}

func (tef *Fault) syncProperty(childId string) {
}

func (tef *Fault) checkProperty(propType, propId string) bool {
	return false
}

func (tef *Fault) checkDuplicate(p Property) (bool, bool, Property) {
	return true, false, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
