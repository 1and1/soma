package somatree

type SomaTreePropertier interface {
	SetProperty(p SomaTreeProperty)

	inheritProperty(p SomaTreeProperty)
	inheritPropertyDeep(p SomaTreeProperty)
	setCustomProperty(p SomaTreeProperty)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
