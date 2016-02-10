package somaproto

type ProtoRequestUser struct {
	User    *ProtoUser       `json:"user,omitempty"`
	Filter  *ProtoUserFilter `json:"filter,omitempty"`
	Restore bool             `json:"restore,omitempty"`
	Purge   bool             `json:"purge,omitempty"`
	//	Credentials ProtoUserCredentials `json:"credentials,omitempty"`
}

type ProtoResultUser struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Users  []ProtoUser `json:"users,omitempty"`
	JobId  string      `json:"jobid,omitempty"`
}

type ProtoUser struct {
	Id             string            `json:"id,omitempty"`
	UserName       string            `json:"username,omitempty"`
	FirstName      string            `json:"firstname,omitempty"`
	LastName       string            `json:"lastname,omitempty"`
	EmployeeNumber string            `json:"employeenumber,omitempty"`
	MailAddress    string            `json:"mailaddress,omitempty"`
	IsActive       bool              `json:"active,omitempty"`
	IsSystem       bool              `json:"system,omitempty"`
	IsDeleted      bool              `json:"deleted,omitempty"`
	Team           string            `json:"team,omitempty"`
	Details        *ProtoUserDetails `json:"details,omitempty"`
}

type ProtoUserCredentials struct {
	Reset    bool   `json:"reset,omitempty"`
	Force    bool   `json:"force,omitempty"`
	Password string `json:"password,omitempty"`
}

type ProtoUserDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

type ProtoUserFilter struct {
	UserName  string `json:"username,omitempty"`
	IsActive  bool   `json:"active,omitempty"`
	IsSystem  bool   `json:"system,omitempty"`
	IsDeleted bool   `json:"deleted,omitempty"`
}

//
func (p *ProtoResultUser) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultUser) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultUser) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultUser) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultUser) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultUser) hasJobId(s string) bool {
	if s != "" {
		return true
	}
	return false
}

func (p *ProtoResultUser) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
