package somaproto

// An Entity is a Type without the golang keyword problem
type Entity struct {
	Name    string         `json:"entity, omitempty"`
	Details *EntityDetails `json:"details, omitempty"`
}

type EntityDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
