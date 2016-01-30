package somaproto

type ProtoRequestUnit struct {
	Unit   *ProtoUnit       `json:"unit,omitempty"`
	Filter *ProtoUnitFilter `json:"filter,omitempty"`
}

type ProtoResultUnit struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Units  []ProtoUnit `json:"units,omitempty"`
	JobId  string      `json:"jobid,omitempty"`
}

type ProtoUnit struct {
	Unit    string            `json:"unit,omitempty"`
	Name    string            `json:"name,omitempty"`
	Details *ProtoUnitDetails `json:"details,omitempty"`
}

type ProtoUnitFilter struct {
	Unit string `json:"unit,omitempty"`
	Name string `json:"name,omitempty"`
}

type ProtoUnitDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultUnit) ErrorMark(err error, imp bool, found bool,
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
	return false
}

func (p *ProtoResultUnit) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultUnit) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultUnit) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
