package somaproto

type ProtoRequestPredicate struct {
	Predicate ProtoPredicate `json:"predicate,omitempty"`
}

type ProtoResultPredicate struct {
	Code       uint16           `json:"code,omitempty"`
	Status     string           `json:"status,omitempty"`
	Text       []string         `json:"text,omitempty"`
	Predicates []ProtoPredicate `json:"predicates,omitempty"`
	JobId      string           `json:"jobid,omitempty"`
}

type ProtoPredicate struct {
	Predicate string                `json:"predicate,omitempty"`
	Details   ProtoPredicateDetails `json:"details,omitempty"`
}

type ProtoPredicateDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
