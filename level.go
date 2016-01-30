package somaproto

type ProtoRequestLevel struct {
	Level  *ProtoLevel       `json:"level,omitempty"`
	Filter *ProtoLevelFilter `json:"filter,omitempty"`
}

type ProtoResultLevel struct {
	Code   uint16       `json:"code,omitempty"`
	Status string       `json:"status,omitempty"`
	Text   []string     `json:"text,omitempty"`
	Levels []ProtoLevel `json:"levels,omitempty"`
	JobId  string       `json:"jobid,omitempty"`
}

type ProtoLevel struct {
	Name      string             `json:"name,omitempty"`
	ShortName string             `json:"shortname,omitempty"`
	Numeric   uint16             `json:"numeric,omitempty"`
	Details   *ProtoLevelDetails `json:"details,omitempty"`
}

type ProtoLevelFilter struct {
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortname,omitempty"`
	Numeric   uint16 `json:"numeric,omitempty"`
}

type ProtoLevelDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultLevel) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultLevel) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultLevel) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultLevel) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultLevel) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
