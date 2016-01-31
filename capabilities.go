package somaproto

type ProtoRequestCapability struct {
	Capability *ProtoCapability       `json:"metric,omitempty"`
	Filter     *ProtoCapabilityFilter `json:"filter,omitempty"`
}

type ProtoResultCapability struct {
	Code         uint16            `json:"code,omitempty"`
	Status       string            `json:"status,omitempty"`
	Text         []string          `json:"text,omitempty"`
	Capabilities []ProtoCapability `json:"metrics,omitempty"`
	JobId        string            `json:"jobid,omitempty"`
}

type ProtoCapability struct {
	Id         string                  `json:"id,omitempty"`
	Monitoring string                  `json:"monitoring,omitempty"`
	Metric     string                  `json:"metric,omitempty"`
	View       string                  `json:"view,omitempty"`
	Thresholds uint64                  `json:"thresholds,omitempty"`
	Name       string                  `json:"name,omitempty"`
	Details    *ProtoCapabilityDetails `json:"details,omitempty"`
}

type ProtoCapabilityFilter struct {
	Monitoring string `json:"monitoring,omitempty"`
	Metric     string `json:"metric,omitempty"`
	View       string `json:"view,omitempty"`
}

type ProtoCapabilityDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultCapability) ErrorMark(err error, imp bool,
	found bool, length int) bool {
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

func (p *ProtoResultCapability) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultCapability) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultCapability) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultCapability) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
