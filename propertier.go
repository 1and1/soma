package tree

type Propertier interface {
	SetProperty(p Property)
	setPropertyInherited(p Property)
	setPropertyOnChildren(p Property)
	addProperty(p Property)

	UpdateProperty(p Property)
	updatePropertyInherited(p Property)
	updatePropertyOnChildren(p Property)
	switchProperty(p Property) bool

	DeleteProperty(p Property)
	deletePropertyInherited(p Property)
	deletePropertyOnChildren(p Property)
	rmProperty(p Property)

	verifySourceInstance(id, prop string) bool
	findIdForSource(source, prop string) string
	syncProperty(childId string)
	checkProperty(propType, propId string) bool
	checkDuplicate(p Property) (bool, bool, Property)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
