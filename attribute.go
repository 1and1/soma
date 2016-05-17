package proto

type Attribute struct {
	Name        string            `json:"name,omitempty"`
	Cardinality string            `json:"cardinality,omitempty"`
	Details     *AttributeDetails `json:"details,omitempty"`
}

type AttributeDetails struct {
	DetailsCreation
}

func NewAttributeRequest() Request {
	return Request{
		Flags:     &Flags{},
		Attribute: &Attribute{},
	}
}

func NewAttributeResult() Result {
	return Result{
		Errors:     &[]string{},
		Attributes: &[]Attribute{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
