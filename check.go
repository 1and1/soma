package somaproto

type Check struct {
	CheckId       string `json:"checkId, omitempty"`
	SourceCheckId string `json:"sourceCheckId, omitempty"`
	CheckConfigId string `json:"checkConfigId, omitempty"`
	SourceType    string `json:"sourceType, omitempty"`
	IsInherited   bool   `json:"isInherited, omitempty"`
	InheritedFrom string `json:"inheritedFrom, omitempty"`
	Inheritance   bool   `json:"inheritance, omitempty"`
	ChildrenOnly  bool   `json:"childrenOnly, omitempty"`
	RepositoryId  string `json:"repositoryId, omitempty"`
	BucketId      string `json:"bucketId, omitempty"`
	CapabilityId  string `json:"capabilityId, omitempty"`
}

func (t *Check) DeepCompare(a *Check) bool {
	if t.CheckId != a.CheckId || t.SourceCheckId != a.SourceCheckId ||
		t.CheckConfigId != a.CheckConfigId || t.SourceType != a.SourceType ||
		t.IsInherited != a.IsInherited || t.InheritedFrom != a.InheritedFrom ||
		t.Inheritance != a.Inheritance || t.ChildrenOnly != a.ChildrenOnly ||
		t.RepositoryId != a.RepositoryId || t.BucketId != a.BucketId ||
		t.CapabilityId != a.CapabilityId {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
