package somaproto

type ProtoRequestBucket struct {
	Bucket  *ProtoBucket       `json:"bucket,omitempty"`
	Filter  *ProtoBucketFilter `json:"filter,omitempty"`
	Restore bool               `json:"restore,omitempty"`
	Purge   bool               `json:"purge,omitempty"`
	Freeze  bool               `json:"freeze,omitempty"`
	Thaw    bool               `json:"thaw,omitempty"`
}

type ProtoResultBucket struct {
	Code    uint16        `json:"code,omitempty"`
	Status  string        `json:"status,omitempty"`
	Text    []string      `json:"text,omitempty"`
	Buckets []ProtoBucket `json:"buckets,omitempty"`
	JobId   string        `json:"jobid,omitempty"`
}

type ProtoBucket struct {
	Id          string              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Repository  string              `json:"repositoryid,omitempty"`
	Team        string              `json:"team,omitempty"`
	Environment string              `json:"environment,omitempty"`
	IsDeleted   bool                `json:"deleted,omitempty"`
	IsFrozen    bool                `json:"frozen,omitempty"`
	Details     *ProtoBucketDetails `json:"details,omitempty"`
	//	Properties []ProtoBucketProperty `json:"properties,omitempty"`
}

type ProtoBucketFilter struct {
	Name         string `json:"name,omitempty"`
	Id           string `json:"id,omitempty"`
	RepositoryId string `json:"repositoryid,omitempty"`
	IsDeleted    bool   `json:"deleted,omitempty"`
	IsFrozen     bool   `json:"frozen,omitempty"`
}

type ProtoBucketDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

//
func (p *ProtoResultBucket) ErrorMark(err error, imp bool, found bool,
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

func (p *ProtoResultBucket) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *ProtoResultBucket) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *ProtoResultBucket) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *ProtoResultBucket) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *ProtoResultBucket) hasJobId(s string) bool {
	if s != "" {
		return true
	}
	return false
}

func (p *ProtoResultBucket) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
