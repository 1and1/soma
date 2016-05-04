package somaproto

type Attribute struct {
	Name        string            `json:"name, omitempty"`
	Cardinality string            `json:"cardinality, omitempty"`
	Details     *AttributeDetails `json:"details, omitempty"`
}

type AttributeDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
