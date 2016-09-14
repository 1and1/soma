package proto

type Status struct {
	Name    string         `json:"name,omitempty"`
	Details *StatusDetails `json:"details,omitempty"`
}

type StatusDetails struct {
	DetailsCreation
}

func NewStatusRequest() Request {
	return Request{
		Flags:  &Flags{},
		Status: &Status{},
	}
}

func NewStatusResult() Result {
	return Result{
		Errors: &[]string{},
		Status: &[]Status{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
