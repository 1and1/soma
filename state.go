package somaproto

type State struct {
	Name    string        `json:"Name, omitempty"`
	Details *StateDetails `json:"details, omitempty"`
}

type StateDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
