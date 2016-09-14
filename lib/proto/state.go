package proto

type State struct {
	Name    string        `json:"Name,omitempty"`
	Details *StateDetails `json:"details,omitempty"`
}

type StateDetails struct {
	DetailsCreation
}

func NewStateRequest() Request {
	return Request{
		Flags: &Flags{},
		State: &State{},
	}
}

func NewStateResult() Result {
	return Result{
		Errors: &[]string{},
		States: &[]State{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
