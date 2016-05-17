package proto

type Mode struct {
	Mode    string       `json:"mode,omitempty"`
	Details *ModeDetails `json:"details,omitempty"`
}

type ModeDetails struct {
	DetailsCreation
}

func NewModeRequest() Request {
	return Request{
		Flags: &Flags{},
		Mode:  &Mode{},
	}
}

func NewModeResult() Result {
	return Result{
		Errors: &[]string{},
		Modes:  &[]Mode{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
