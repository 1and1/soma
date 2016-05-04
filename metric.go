package somaproto

type ProtoRequestMetric struct {
	Metric *ProtoMetric       `json:"metric,omitempty"`
	Filter *ProtoMetricFilter `json:"filter,omitempty"`
}

type ProtoResultMetric struct {
	Code    uint16        `json:"code,omitempty"`
	Status  string        `json:"status,omitempty"`
	Text    []string      `json:"text,omitempty"`
	Metrics []ProtoMetric `json:"metrics,omitempty"`
	JobId   string        `json:"jobid,omitempty"`
}

type ProtoMetric struct {
	Metric      string                        `json:"metric,omitempty"`
	Unit        string                        `json:"unit,omitempty"`
	Description string                        `json:"description,omitempty"`
	Packages    *[]ProtoMetricProviderPackage `json:"packages,omitempty"`
	Details     *ProtoMetricDetails           `json:"details,omitempty"`
}

type ProtoMetricFilter struct {
	Unit     string `json:"unit,omitempty"`
	Provider string `json:"provider,omitempty"`
	Package  string `json:"package,omitempty"`
}

type ProtoMetricDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

type ProtoMetricProviderPackage struct {
	Provider string `json:"provider,omitempty"`
	Package  string `json:"package,omitempty"`
}

//
func (p *ProtoMetric) DeepCompare(a *ProtoMetric) bool {
	if p.Metric != a.Metric || p.Unit != a.Unit || p.Description != a.Description {
		return false
	}
packageloop:
	for _, pkg := range *p.Packages {
		if pkg.DeepCompareSlice(a.Packages) {
			continue packageloop
		}
		return false
	}
	return true
}

func (p *ProtoMetricProviderPackage) DeepCompare(a *ProtoMetricProviderPackage) bool {
	if p.Provider != a.Provider || p.Package != a.Package {
		return false
	}
	return true
}

func (p *ProtoMetricProviderPackage) DeepCompareSlice(a *[]ProtoMetricProviderPackage) bool {
	for _, pkg := range *a {
		if p.DeepCompare(&pkg) {
			return true
		}
	}
	return false
}

//
func (p *ProtoResultMetric) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultMetric) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultMetric) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultMetric) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultMetric) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultMetric) hasJobId(s string) bool {
	if s != "" {
		p.JobId = s
		return true
	}
	return false
}

func (p *ProtoResultMetric) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
