package somaproto

type AttributeRequest struct {
	Attribute *Attribute `json:"attribute,omitempty"`
}

type AttributeResult struct {
	Code       uint16      `json:"code,omitempty"`
	Status     string      `json:"status,omitempty"`
	Text       []string    `json:"text,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
	JobId      string      `json:"jobid,omitempty"`
}

type Attribute struct {
	Attribute   string            `json:"attribute,omitempty"`
	Cardinality string            `json:"cardinality,omitempty"`
	Details     *AttributeDetails `json:"details,omitempty"`
}

type AttributeDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *AttributeResult) ErrorMark(err error, imp bool, found bool,
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

func (p *AttributeResult) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *AttributeResult) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *AttributeResult) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *AttributeResult) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *AttributeResult) hasJobId(s string) bool {
	if s != "" {
		p.JobId = s
		return true
	}
	return false
}

func (p *AttributeResult) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
