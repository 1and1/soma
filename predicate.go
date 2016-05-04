package somaproto

type Predicate struct {
	Symbol  string            `json:"symbol, omitempty"`
	Details *PredicateDetails `json:"details, omitempty"`
}

type PredicateDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
