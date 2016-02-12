package somaproto

type ProtoRequestRepository struct {
	Repository *ProtoRepository       `json:"repository,omitempty"`
	Filter     *ProtoRepositoryFilter `json:"filter,omitempty"`
	Restore    bool                   `json:"restore,omitempty"`
	Purge      bool                   `json:"purge,omitempty"`
	Clear      bool                   `json:"clear,omitempty"`
	Activate   bool                   `json:"activate,omitempty"`
}

type ProtoResultRepository struct {
	Code         uint16            `json:"code,omitempty"`
	Status       string            `json:"status,omitempty"`
	Text         []string          `json:"text,omitempty"`
	Repositories []ProtoRepository `json:"repositories,omitempty"`
	JobId        string            `json:"jobid,omitempty"`
}

type ProtoRepository struct {
	Id         string                    `json:"id,omitempty"`
	Name       string                    `json:"name,omitempty"`
	Team       string                    `json:"team,omitempty"`
	IsDeleted  bool                      `json:"deleted,omitempty"`
	IsActive   bool                      `json:"active,omitempty"`
	Details    *ProtoRepositoryDetails   `json:"details,omitempty"`
	Properties []ProtoRepositoryProperty `json:"properties,omitempty"`
}

type ProtoRepositoryFilter struct {
	Name      string `json:"name,omitempty"`
	Team      string `json:"team,omitempty"`
	IsDeleted bool   `json:"deleted,omitempty"`
	IsActive  bool   `json:"active,omitempty"`
}

type ProtoRepositoryDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

type ProtoRepositoryProperty struct {
	Type         string `json:"type,omitempty"`
	View         string `json:"view,omitempty"`
	Property     string `json:"property,omitempty"`
	Value        string `json:"value,omitempty"`
	Inheritance  bool   `json:"inheritance,omitempty"`
	ChildrenOnly bool   `json:"children,omitempty"`
	Source       string `json:"source,omitempty"`
}

//
func (p *ProtoResultRepository) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultRepository) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultRepository) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultRepository) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultRepository) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultRepository) hasJobId(s string) bool {
	if s != "" {
		p.JobId = s
		return true
	}
	return false
}

func (p *ProtoResultRepository) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
