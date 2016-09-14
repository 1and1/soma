package proto

// An Entity is a Type without the golang keyword problem
type Entity struct {
	Name    string         `json:"entity,omitempty"`
	Details *EntityDetails `json:"details,omitempty"`
}

type EntityDetails struct {
	DetailsCreation
}

func NewEntityRequest() Request {
	return Request{
		Flags:  &Flags{},
		Entity: &Entity{},
	}
}

func NewEntityResult() Result {
	return Result{
		Errors:   &[]string{},
		Entities: &[]Entity{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
