package somaproto

type Property struct {
	PropertyType     string           `json:"propertyType"`
	RepositoryId     string           `json:"repositoryId, omitempty"`
	BucketId         string           `json:"bucketId, omitempty"`
	InstanceId       string           `json:"instanceId, omitempty"`
	View             string           `json:"view, omitempty"`
	Inheritance      bool             `json:"inheritance"`
	ChildrenOnly     bool             `json:"childrenOnly"`
	IsInherited      bool             `json:"isInherited, omitempty"`
	SourceInstanceId string           `json:"sourceInstanceId, omitempty"`
	SourceType       string           `json:"sourceType, omitempty"`
	InheritedFrom    string           `json:"inheritedFrom, omitempty"`
	Custom           *PropertyCustom  `json:"custom, omitempty"`
	System           *PropertySystem  `json:"system, omitempty"`
	Service          *PropertyService `json:"service, omitempty"`
	Native           *PropertyNative  `json:"native, omitempty"`
	Oncall           *PropertyOncall  `json:"oncall, omitempty"`
	Details          *PropertyDetails `json:"details, omitempty"`
}

type PropertyFilter struct {
	Name         string `json:"name, omitempty"`
	Type         string `json:"type, omitempty"`
	RepositoryId string `json:"repositoryId, omitempty"`
}

type PropertyDetails struct {
	DetailsCreation
}

type PropertyCustom struct {
	CustomId     string `json:"customId, omitempty"`
	RepositoryId string `json:"repositoryId, omitempty"`
	Name         string `json:"name, omitempty"`
	Value        string `json:"value, omitempty"`
}

type PropertySystem struct {
	Name  string `json:"name, omitempty"`
	Value string `json:"value, omitempty"`
}

type PropertyService struct {
	Name       string             `json:"name, omitempty"`
	TeamId     string             `json:"teamId, omitempty"`
	Attributes []ServiceAttribute `json:"attributes"`
}

type PropertyNative struct {
	Name  string `json:"name, omitempty"`
	Value string `json:"value, omitempty"`
}

type PropertyOncall struct {
	OncallId string `json:"oncallId, omitempty"`
	Name     string `json:"name, omitempty"`
	Number   string `json:"number, omitempty"`
}

type ServiceAttribute struct {
	Name  string `json:"name, omitempty"`
	Value string `json:"value, omitempty"`
}

func (t *PropertyService) DeepCompare(a *PropertyService) bool {
	if t.Name != a.Name || t.TeamId != a.TeamId {
		return false
	}
attrloop:
	for _, attr := range t.Attributes {
		if attr.DeepCompareSlice(&a.Attributes) {
			continue attrloop
		}
		return false
	}
	return true
}

func (t *ServiceAttribute) DeepCompare(a *ServiceAttribute) bool {
	if t.Name != a.Name || t.Value != a.Value {
		return false
	}
	return true
}

func (t *ServiceAttribute) DeepCompareSlice(a *[]ServiceAttribute) bool {
	for _, attr := range *a {
		if t.DeepCompare(&attr) {
			return true
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
