package somaproto

type Environment struct {
	Name    string              `json:"name, omitempty"`
	Details *EnvironmentDetails `json:"details, omitempty"`
}

type EnvironmentDetails struct {
	DetailsCreation
}

func NewEnvironmentRequest() Request {
	return Request{
		Environment: &Environment{},
	}
}

func NewEnvironmentResult() Result {
	return Result{
		Errors:       &[]string{},
		Environments: &[]Environment{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
