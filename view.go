package somaproto

type View struct {
	Name    string       `json:"name, omitempty"`
	Details *ViewDetails `json:"details, omitempty"`
}

type ViewDetails struct {
	DetailsCreation
}

func NewViewRequest() Request {
	return Request{
		View: &View{},
	}
}

func NewViewResult() Result {
	return Result{
		Errors: &[]string{},
		Views:  &[]View{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
