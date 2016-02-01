package somaproto

type ProtoRequestAttribute struct {
	Attribute *ProtoAttribute `json:"attribute,omitempty"`
}

type ProtoResultAttribute struct {
	Code       uint16           `json:"code,omitempty"`
	Status     string           `json:"status,omitempty"`
	Text       []string         `json:"text,omitempty"`
	Attributes []ProtoAttribute `json:"attributes,omitempty"`
	JobId      string           `json:"jobid,omitempty"`
}

type ProtoAttribute struct {
	Attribute string                 `json:"attribute,omitempty"`
	Details   *ProtoAttributeDetails `json:"details,omitempty"`
}

type ProtoAttributeDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultAttribute) ErrorMark(err error, imp bool, found bool,
	length int) bool {
	if p.markError(err) {
		return true
	}
	if p.markImplemented(imp) {
		return true
	}
	if p.markFound(found, length) {
		return true
	}
	return p.markOk()
}

func (p *ProtoResultAttribute) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultAttribute) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultAttribute) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultAttribute) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
