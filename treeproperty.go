package somaproto

type TreeProperty struct {
	InstanceId       string               `json:"instance_id,omitempty"`
	SourceInstanceId string               `json:"source_instance_id,omitempty"`
	SourceType       string               `json:"source_type,omitempty"`
	IsInherited      bool                 `json:"is_inherited,omitempty"`
	InheritedFrom    string               `json:"inherited_from,omitempty"`
	Inheritance      bool                 `json:"inheritance"`
	ChildrenOnly     bool                 `json:"children_only"`
	View             string               `json:"view,omitempty"`
	PropertyType     string               `json:"property_type"`
	RepositoryId     string               `json:"repository_id,omitempty"`
	BucketId         string               `json:"bucket_id,omitempty"`
	Custom           *TreePropertyCustom  `json:"custom,omitempty"`
	System           *TreePropertySystem  `json:"system,omitempty"`
	Service          *TreePropertyService `json:"service,omitempty"`
	Native           *TreePropertyNative  `json:"native,omitempty"`
	Oncall           *TreePropertyOncall  `json:"oncall,omitempty"`
}

type TreePropertyCustom struct {
	CustomId     string `json:"custom_id,omitempty"`
	RepositoryId string `json:"repository_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Value        string `json:"value,omitempty"`
}

type TreePropertySystem struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type TreePropertyService struct {
	Name       string                 `json:"name,omitempty"`
	TeamId     string                 `json:"team_id,omitempty"`
	Attributes []TreeServiceAttribute `json:"attributes"`
}

type TreePropertyNative struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type TreePropertyOncall struct {
	OncallId string `json:"oncall_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Number   string `json:"number,omitempty"`
}

type TreeServiceAttribute struct {
	Attribute string `json:"attribute,omitempty"`
	Value     string `json:"value,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
