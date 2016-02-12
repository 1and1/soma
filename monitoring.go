package somaproto

type ProtoRequestMonitoring struct {
	Monitoring *ProtoMonitoring       `json:"monitoring,omitempty"`
	Filter     *ProtoMonitoringFilter `json:"filter,omitempty"`
}

type ProtoResultMonitoring struct {
	Code    uint16            `json:"code,omitempty"`
	Status  string            `json:"status,omitempty"`
	Text    []string          `json:"text,omitempty"`
	Systems []ProtoMonitoring `json:"systems,omitempty"`
	JobId   string            `json:"jobid,omitempty"`
}

type ProtoMonitoring struct {
	Id       string                  `json:"id,omitempty"`
	Name     string                  `json:"name,omitempty"`
	Mode     string                  `json:"mode,omitempty"`
	Contact  string                  `json:"contact,omitempty"`
	Team     string                  `json:"team,omitempty"`
	Callback string                  `json:"callback,omitempty"`
	Details  *ProtoMonitoringDetails `json:"details,omitempty"`
}

type ProtoMonitoringFilter struct {
	Name    string `json:"name,omitempty"`
	Mode    string `json:"mode,omitempty"`
	Contact string `json:"contact,omitempty"`
	Team    string `json:"team,omitempty"`
}

type ProtoMonitoringDetails struct {
	Checks    uint64 `json:"checks,omitempty"`
	Instances uint64 `json:"instances,omitempty"`
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultMonitoring) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultMonitoring) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultMonitoring) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultMonitoring) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultMonitoring) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultMonitoring) hasJobId(s string) bool {
	if s != "" {
		p.JobId = s
		return true
	}
	return false
}

func (p *ProtoResultMonitoring) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
