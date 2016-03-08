package somaproto

type CheckConfigurationRequest struct {
	CheckConfiguration *CheckConfiguration       `json:"check_configuration,omitempty"`
	Filter             *CheckConfigurationFilter `json:"filter,omitempty"`
}

type CheckConfigurationResult struct {
	Code                uint16               `json:"code,omitempty"`
	Status              string               `json:"status,omitempty"`
	Text                []string             `json:"text,omitempty"`
	CheckConfigurations []CheckConfiguration `json:"check_configurations,omitempty"`
	JobId               string               `json:"jobid,omitempty"`
}

type CheckConfiguration struct {
	Id           string                         `json:"id,omitempty"`
	Name         string                         `json:"name,omitempty"`
	Interval     uint64                         `json:"interval,omitempty"`
	RepositoryId string                         `json:"repository_id,omitempty"`
	BucketId     string                         `json:"bucket_id,omitempty"`
	CapabilityId string                         `json:"capability_id,omitempty"`
	ObjectId     string                         `json:"object_id,omitempty"`
	ObjectType   string                         `json:"object_type,omitempty"`
	IsActive     bool                           `json:"is_active,omitempty"`
	IsEnabled    bool                           `json:"is_enabled,omitempty"`
	Inheritance  bool                           `json:"inheritance,omitempty"`
	ChildrenOnly bool                           `json:"children_only,omitempty"`
	ExternalId   string                         `json:"external_id,omitempty"`
	Constraints  []CheckConfigurationConstraint `json:"constraints,omitempty"`
	Thresholds   []CheckConfigurationThreshold  `json:"thresholds,omitempty"`
	Details      *CheckConfigurationDetails     `json:"details,omitempty"`
}

type CheckConfigurationConstraint struct {
	ConstraintType string                `json:"constraint_type,omitempty"`
	Native         *TreePropertyNative   `json:"native,omitempty"`
	Oncall         *TreePropertyOncall   `json:"oncall,omitempty"`
	Custom         *TreePropertyCustom   `json:"custom,omitempty"`
	System         *TreePropertySystem   `json:"system,omitempty"`
	Service        *TreePropertyService  `json:"service,omitempty"`
	Attribute      *TreeServiceAttribute `json:"attribute,omitempty"`
}

type CheckConfigurationThreshold struct {
	Predicate ProtoPredicate
	Level     ProtoLevel
	Value     int64
}

func (c *CheckConfigurationThreshold) DeepCompareSlice(a []CheckConfigurationThreshold) bool {
	for _, thr := range a {
		if c.DeepCompare(&thr) {
			return true
		}
	}
	return false
}

func (c *CheckConfigurationThreshold) DeepCompare(a *CheckConfigurationThreshold) bool {
	if c.Value != a.Value || c.Level.Name != a.Level.Name ||
		c.Predicate.Predicate != a.Predicate.Predicate {
		return false
	}
	return true
}

type CheckConfigurationDetails struct {
	CreatedAt string `json:"created_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
}

type CheckConfigurationFilter struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	CapabilityId string `json:"capability_id,omitempty"`
}

//
func (c *CheckConfiguration) DeepCompare(a *CheckConfiguration) bool {
	if a == nil {
		return false
	}
	if c.Id != a.Id || c.Name != a.Name || c.Interval != a.Interval ||
		c.RepositoryId != a.RepositoryId || c.BucketId != a.BucketId ||
		c.CapabilityId != a.CapabilityId || c.ObjectId != a.ObjectId ||
		c.ObjectType != a.ObjectType || c.IsActive != a.IsActive ||
		c.IsEnabled != a.IsEnabled || c.Inheritance != a.Inheritance ||
		c.ChildrenOnly != a.ChildrenOnly || c.ExternalId != a.ExternalId {
		return false
	}
threshloop:
	for _, thr := range c.Thresholds {
		if thr.DeepCompareSlice(a.Thresholds) {
			continue threshloop
		}
		return false
	}
	// TODO: constraints
	return true
}

//
func (p *CheckConfigurationResult) ErrorMark(err error, imp bool, found bool,
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

func (p *CheckConfigurationResult) markError(err error) bool {
	if err != nil {
		p.Code = 500
		p.Status = "ERROR"
		p.Text = []string{err.Error()}
		return true
	}
	return false
}

func (p *CheckConfigurationResult) markImplemented(f bool) bool {
	if f {
		p.Code = 501
		p.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (p *CheckConfigurationResult) markFound(f bool, i int) bool {
	if f || i == 0 {
		p.Code = 404
		p.Status = "NOT FOUND"
		return true
	}
	return false
}

func (p *CheckConfigurationResult) markOk() bool {
	p.Code = 200
	p.Status = "OK"
	return false
}

func (p *CheckConfigurationResult) hasJobId(s string) bool {
	if s != "" {
		p.JobId = s
		return true
	}
	return false
}

func (p *CheckConfigurationResult) markAccepted() bool {
	p.Code = 202
	p.Status = "ACCEPTED"
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
