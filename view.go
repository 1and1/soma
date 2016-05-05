package somaproto

type View struct {
	Name    string       `json:"name, omitempty"`
	Details *ViewDetails `json:"details, omitempty"`
}

type ViewDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
