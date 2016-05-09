package somaproto

type Validity struct {
	SystemProperty string           `json:"systemProperty,omitempty"`
	ObjectType     string           `json:"objectType,omitempty"`
	Direct         bool             `json:"direct, string"`
	Inherited      bool             `json:"inherited, string"`
	Details        *ValidityDetails `json:"details,omitempty"`
}

type ValidityDetails struct {
	DetailsCreation
}

func NewValidityRequest() Request {
	return Request{
		Validity: &Validity{},
	}
}

func NewValidityResult() Result {
	return Result{
		Errors:     &[]string{},
		Validities: &[]Validity{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
