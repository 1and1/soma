package somaproto

type ProtoRequestMode struct {
	Mode *ProtoMode `json:"mode,omitempty"`
}

type ProtoResultMode struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Modes  []ProtoMode `json:"modes,omitempty"`
	JobId  string      `json:"jobid,omitempty"`
}

type ProtoMode struct {
	Mode    string            `json:"mode,omitempty"`
	Details *ProtoModeDetails `json:"details,omitempty"`
}

type ProtoModeDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultMode) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultMode) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultMode) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultMode) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultMode) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultMode) hasJobId(s string) bool {
	if s != "" {
		p.JobId = s
		return true
	}
	return false
}

func (p *ProtoResultMode) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
