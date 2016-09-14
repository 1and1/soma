/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

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
	getCurrentProperty(p Property) Property

	DeleteProperty(p Property)
	deletePropertyInherited(p Property)
	deletePropertyOnChildren(p Property)
	deletePropertyAllInherited()
	deletePropertyAllLocal()
	rmProperty(p Property) bool

	verifySourceInstance(id, prop string) bool
	findIdForSource(source, prop string) string
	syncProperty(childId string)
	checkProperty(propType, propId string) bool
	checkDuplicate(p Property) (bool, bool, Property)
	resyncProperty(srcId, pType, childId string)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
