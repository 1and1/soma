package proto

type CheckConfig struct {
	Id           string                  `json:"id,omitempty"`
	Name         string                  `json:"name,omitempty"`
	Interval     uint64                  `json:"interval,omitempty"`
	RepositoryId string                  `json:"repositoryId,omitempty"`
	BucketId     string                  `json:"bucketId,omitempty"`
	CapabilityId string                  `json:"capabilityId,omitempty"`
	ObjectId     string                  `json:"objectId,omitempty"`
	ObjectType   string                  `json:"objectType,omitempty"`
	IsActive     bool                    `json:"isActive,omitempty"`
	IsEnabled    bool                    `json:"isEnabled,omitempty"`
	Inheritance  bool                    `json:"inheritance,omitempty"`
	ChildrenOnly bool                    `json:"childrenOnly,omitempty"`
	ExternalId   string                  `json:"externalId,omitempty"`
	Constraints  []CheckConfigConstraint `json:"constraints,omitempty"`
	Thresholds   []CheckConfigThreshold  `json:"thresholds,omitempty"`
	Details      *CheckConfigDetails     `json:"details,omitempty"`
}

type CheckConfigConstraint struct {
	ConstraintType string            `json:"constraintType,omitempty"`
	Native         *PropertyNative   `json:"native,omitempty"`
	Oncall         *PropertyOncall   `json:"oncall,omitempty"`
	Custom         *PropertyCustom   `json:"custom,omitempty"`
	System         *PropertySystem   `json:"system,omitempty"`
	Service        *PropertyService  `json:"service,omitempty"`
	Attribute      *ServiceAttribute `json:"attribute,omitempty"`
}

func (c *CheckConfigConstraint) DeepCompare(a *CheckConfigConstraint) bool {
	if c.ConstraintType != a.ConstraintType {
		return false
	}
	switch c.ConstraintType {
	case "native":
		if c.Native.DeepCompare(a.Native) {
			return true
		}
	case "oncall":
		if c.Oncall.DeepCompare(a.Oncall) {
			return true
		}
	case "custom":
		if c.Custom.DeepCompare(a.Custom) {
			return true
		}
	case "system":
		if c.System.DeepCompare(a.System) {
			return true
		}
	case "service":
		if c.Service.DeepCompare(a.Service) {
			return true
		}
	case "attribute":
		if c.Attribute.DeepCompare(a.Attribute) {
			return true
		}
	}
	return false
}

func (c *CheckConfigConstraint) DeepCompareSlice(a []CheckConfigConstraint) bool {
	if a == nil {
		return false
	}
	for _, constr := range a {
		if c.DeepCompare(&constr) {
			return true
		}
	}
	return false
}

type CheckConfigThreshold struct {
	Predicate Predicate
	Level     Level
	Value     int64
}

func (c *CheckConfigThreshold) DeepCompareSlice(a []CheckConfigThreshold) bool {
	if a == nil {
		return false
	}
	for _, thr := range a {
		if c.DeepCompare(&thr) {
			return true
		}
	}
	return false
}

func (c *CheckConfigThreshold) DeepCompare(a *CheckConfigThreshold) bool {
	if c.Value != a.Value || c.Level.Name != a.Level.Name ||
		c.Predicate.Symbol != a.Predicate.Symbol {
		return false
	}
	return true
}

type CheckConfigDetails struct {
	CreatedAt string              `json:"createdAt,omitempty"`
	CreatedBy string              `json:"createdBy,omitempty"`
	Instances []CheckInstanceInfo `json:"instances,omitempty"`
}

type CheckConfigFilter struct {
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	CapabilityId string `json:"capabilityId,omitempty"`
}

type CheckInstanceInfo struct {
	Id            string `json:"id,omitempty"`
	ObjectId      string `json:"objectId,omitempty"`
	ObjectType    string `json:"objectType,omitempty"`
	CurrentStatus string `json:"currentStatus,omitempty"`
	NextStatus    string `json:"nextStatus,omitempty"`
}

func (c *CheckConfig) DeepCompare(a *CheckConfig) bool {
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
revthreshloop:
	for _, thr := range a.Thresholds {
		if thr.DeepCompareSlice(c.Thresholds) {
			continue revthreshloop
		}
		return false
	}
constrloop:
	for _, constr := range c.Constraints {
		if constr.DeepCompareSlice(a.Constraints) {
			continue constrloop
		}
		return false
	}
revconstrloop:
	for _, constr := range a.Constraints {
		if constr.DeepCompareSlice(c.Constraints) {
			continue revconstrloop
		}
		return false
	}
	return true
}

func NewCheckConfigRequest() Request {
	return Request{
		Flags:       &Flags{},
		CheckConfig: &CheckConfig{},
	}
}

func NewCheckConfigFilter() Request {
	return Request{
		Filter: &Filter{
			CheckConfig: &CheckConfigFilter{},
		},
	}
}

func NewCheckConfigResult() Result {
	return Result{
		Errors:       &[]string{},
		CheckConfigs: &[]CheckConfig{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
