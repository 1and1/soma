package somaproto

type Category struct {
	Name    string           `json:"name, omitempty"`
	Details *CategoryDetails `json:"details, omitempty"`
}

type CategoryDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
