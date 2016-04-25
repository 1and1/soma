package somaproto

type PropertyRequest struct {
	PropertyType string               `json:"propertytype,omitempty"`
	Custom       *TreePropertyCustom  `json:"custom,omitempty"`
	System       *TreePropertySystem  `json:"system,omitempty"`
	Service      *TreePropertyService `json:"service,omitempty"`
	Native       *TreePropertyNative  `json:"native,omitempty"`
	Filter       *PropertyFilter      `json:"filter,omitempty"`
}

type PropertyResult struct {
	Code    uint16                `json:"code,omitempty"`
	Status  string                `json:"status,omitempty"`
	Text    []string              `json:"text,omitempty"`
	Custom  []TreePropertyCustom  `json:"custom,omitempty"`
	System  []TreePropertySystem  `json:"system,omitempty"`
	Service []TreePropertyService `json:"service,omitempty"`
	Native  []TreePropertyNative  `json:"native,omitempty"`
	JobId   string                `json:"jobid,omitempty"`
}

type PropertyFilter struct {
	Property   string `json:"property,omitempty"`
	Type       string `json:"type,omitempty"`
	Repository string `json:"repository,omitempty"`
}

//
func (p *PropertyResult) ErrorMark(err error, imp bool, found bool,
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

func (p *PropertyResult) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *PropertyResult) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *PropertyResult) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *PropertyResult) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *PropertyResult) hasJobId(s string) bool {
	if s != "" {
		p.JobId = s
		return true
	}
	return false
}

func (p *PropertyResult) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
