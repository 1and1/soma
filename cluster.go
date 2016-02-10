package somaproto

type ProtoRequestCluster struct {
	Cluster *ProtoCluster       `json:"cluster,omitempty"`
	Filter  *ProtoClusterFilter `json:"filter,omitempty"`
}

type ProtoResultCluster struct {
	Code     uint16         `json:"code,omitempty"`
	Status   string         `json:"status,omitempty"`
	Text     []string       `json:"text,omitempty"`
	Clusters []ProtoCluster `json:"clusters,omitempty"`
	JobId    string         `json:"jobid,omitempty"`
}

type ProtoCluster struct {
	Id          string                 `json:"id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	BucketId    string                 `json:"bucketid,omitempty"`
	ObjectState string                 `json:"objectstate,omitempty"`
	TeamId      string                 `json:"teamid,omitempty"`
	Members     []ProtoNode            `json:"members,omitempty"`
	Details     *ProtoClusterDetails   `json:"details,omitempty"`
	Properties  []ProtoClusterProperty `json:"properties,omitempty"`
}

type ProtoClusterFilter struct {
	Name     string `json:"name,omitempty"`
	BucketId string `json:"bucketid,omitempty"`
}

type ProtoClusterDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

type ProtoClusterProperty struct {
	Type         string `json:"type,omitempty"`
	View         string `json:"view,omitempty"`
	Property     string `json:"property,omitempty"`
	Value        string `json:"value,omitempty"`
	Inheritance  bool   `json:"inheritance,omitempty"`
	ChildrenOnly bool   `json:"children,omitempty"`
	Source       string `json:"source,omitempty"`
}

//
func (p *ProtoResultCluster) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultCluster) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultCluster) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultCluster) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultCluster) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultCluster) hasJobId(s string) bool {
	if s != "" {
		return true
	}
	return false
}

func (p *ProtoResultCluster) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
