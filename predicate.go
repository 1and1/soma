package somaproto

type Predicate struct {
	Symbol  string            `json:"symbol, omitempty"`
	Details *PredicateDetails `json:"details, omitempty"`
}

type PredicateDetails struct {
	DetailsCreation
}

func NewPredicateRequest() Request {
	return Request{
		Predicate: &Predicate{},
	}
}

func NewPredicateResult() Result {
	return Result{
		Errors:     &[]string{},
		Predicates: &[]Predicate{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
