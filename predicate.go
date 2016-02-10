package somaproto

type ProtoRequestPredicate struct {
	Predicate *ProtoPredicate `json:"predicate,omitempty"`
}

type ProtoResultPredicate struct {
	Code       uint16           `json:"code,omitempty"`
	Status     string           `json:"status,omitempty"`
	Text       []string         `json:"text,omitempty"`
	Predicates []ProtoPredicate `json:"predicates,omitempty"`
	JobId      string           `json:"jobid,omitempty"`
}

type ProtoPredicate struct {
	Predicate string                 `json:"predicate,omitempty"`
	Details   *ProtoPredicateDetails `json:"details,omitempty"`
}

type ProtoPredicateDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultPredicate) ErrorMark(err error, imp bool, found bool,
	length int, jobid string) bool {
	if p.markError(err) {
		return true
	}
	if p.markImplemented(imp) {
		return true
	}
	if p.markFound(found, length) {
		return true
	}
	if p.hasJobId(jobid) {
		return p.markAccepted()
	}
	return p.markOk()
}

func (p *ProtoResultPredicate) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultPredicate) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultPredicate) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultPredicate) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultPredicate) hasJobId(s string) bool {
	if s != "" {
		return true
	}
	return false
}

func (p *ProtoResultPredicate) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
