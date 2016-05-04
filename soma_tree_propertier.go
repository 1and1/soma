package somatree

type SomaTreePropertier interface {
	SetProperty(p SomaTreeProperty)
	//DelProperty(p SomaTreeProperty)

	inheritProperty(p SomaTreeProperty)
	inheritPropertyDeep(p SomaTreeProperty)
	setCustomProperty(p SomaTreeProperty)
	setServiceProperty(p SomaTreeProperty)
	setSystemProperty(p SomaTreeProperty)
	setOncallProperty(p SomaTreeProperty)
	//deleteProperty(p SomaTreeProperty)
	//deletePropertyDeep(p SomaTreeProperty)
	//deleteCustomProperty(p SomaTreeProperty)
	//deleteServiceProperty(p SomaTreeProperty)
	//deleteSystemProperty(p SomaTreeProperty)
	//deleteOncallProperty(p SomaTreeProperty)
	syncProperty(childId string)
	checkProperty(propType string, propId string) bool
	checkDuplicate(p SomaTreeProperty) (bool, bool, SomaTreeProperty)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
