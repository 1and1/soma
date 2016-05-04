package somaproto

type Mode struct {
	Mode    string       `json:"mode, omitempty"`
	Details *ModeDetails `json:"details, omitempty"`
}

type ModeDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
