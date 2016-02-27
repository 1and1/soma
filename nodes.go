package somaproto

type ProtoRequestNode struct {
	Node    *ProtoNode       `json:"node,omitempty"`
	Filter  *ProtoNodeFilter `json:"filter,omitempty"`
	Restore bool             `json:"restore,omitempty"`
	Purge   bool             `json:"purge,omitempty"`
}

type ProtoResultNode struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Nodes  []ProtoNode `json:"nodes,omitempty"`
	JobId  string      `json:"jobid,omitempty"`
}

type ProtoNode struct {
	Id         string            `json:"id,omitempty"`
	AssetId    uint64            `json:"assetid,omitempty"`
	Name       string            `json:"name,omitempty"`
	Team       string            `json:"team,omitempty"`
	Server     string            `json:"server,omitempty"`
	State      string            `json:"state,omitempty"`
	IsOnline   bool              `json:"online,omitempty"`
	IsDeleted  bool              `json:"deleted,omitempty"`
	Details    *ProtoNodeDetails `json:"details,omitempty"`
	Config     *ProtoNodeConfig  `json:"config,omitempty"`
	Properties *[]TreeProperty   `json:"properties,omitempty"`
}

type ProtoNodeDetails struct {
	CreatedAt string      `json:"createdat,omitempty"`
	CreatedBy string      `json:"createdby,omitempty"`
	Server    ProtoServer `json:"server,omitempty"`
}

type ProtoNodeFilter struct {
	Name          string `json:"name,omitempty"`
	Team          string `json:"team,omitempty"`
	Server        string `json:"server,omitempty"`
	Online        bool   `json:"online,omitempty"`
	NotOnline     bool   `json:"notonline,omitempty"`
	Deleted       bool   `json:"deleted,omitempty"`
	NotDeleted    bool   `json:"notdeleted,omitempty"`
	PropertyType  string `json:"propertytype,omitempty"`
	Property      string `json:"property,omitempty"`
	LocalProperty bool   `json:"localproperty,omitempty"`
}

type ProtoNodeConfig struct {
	RepositoryId string `json:"repository_id,omitempty"`
	BucketId     string `json:"bucket_id,omitempty"`
}

//
func (p *ProtoResultNode) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultNode) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultNode) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultNode) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultNode) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultNode) hasJobId(s string) bool {
	if s != "" {
		p.JobId = s
		return true
	}
	return false
}

func (p *ProtoResultNode) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
