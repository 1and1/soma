package somaproto

type ProtoRequestOncall struct {
	OnCall  *ProtoOncall         `json:"oncall,omitempty"`
	Filter  *ProtoOncallFilter   `json:"filter,omitempty"`
	Members *[]ProtoOncallMember `json:"members,omitempty"`
}

type ProtoResultOncall struct {
	Code    uint16        `json:"code,omitempty"`
	Status  string        `json:"status,omitempty"`
	Text    []string      `json:"text,omitempty"`
	Oncalls []ProtoOncall `json:"oncalls,omitempty"`
	JobId   string        `json:"jobid,omitempty"`
}

type ProtoOncall struct {
	Id      string              `json:"id,omitempty"`
	Name    string              `json:"name,omitempty"`
	Number  string              `json:"number,omitempty"`
	Details *ProtoOncallDetails `json:"details,omitempty"`
}

type ProtoOncallDetails struct {
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	Members   []string `json:"members,omitempty"`
}

type ProtoOncallMember struct {
	UserName string `json:"username,omitempty"`
	UserId   string `json"userid,omitempty"`
}

type ProtoOncallFilter struct {
	Name   string `json:"name,omitempty"`
	Number string `json:"number,omitempty"`
}

//
func (p *ProtoResultOncall) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultOncall) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultOncall) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultOncall) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultOncall) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
