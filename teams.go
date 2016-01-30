package somaproto

type ProtoRequestTeam struct {
	Team   *ProtoTeam       `json:"team,omitempty"`
	Filter *ProtoTeamFilter `json:"filter,omitempty"`
}

type ProtoResultTeam struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Teams  []ProtoTeam `json:"teams,omitempty"`
	JobId  string      `json:"jobid,omitempty"`
}

type ProtoTeam struct {
	Id      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Ldap    string            `json:"ldap,omitempty"`
	System  bool              `json:"system,omitempty"`
	Details *ProtoTeamDetails `json:"details,omitempty"`
}

type ProtoTeamDetails struct {
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	Members   []string `json:"members,omitempty"`
}

type ProtoTeamFilter struct {
	Name   string `json:"name,omitempty"`
	Ldap   string `json:"ldap,omitempty"`
	System bool   `json:"system,omitempty"`
}

//
func (p *ProtoResultTeam) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultTeam) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultTeam) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultTeam) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultTeam) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
