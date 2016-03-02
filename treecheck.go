package somaproto

type TreeCheck struct {
	CheckId       string `json:"check_id,omitempty"`
	SourceCheckId string `json:"source_check_id,omitempty"`
	CheckConfigId string `json:"check_config_id,omitempty"`
	SourceType    string `json:"source_type,omitempty"`
	IsInherited   bool   `json:"is_inherited,omitempty"`
	InheritedFrom string `json:"inherited_from,omitempty"`
	Inheritance   bool   `json:"inheritance,omitempty"`
	ChildrenOnly  bool   `json:"children_only,omitempty"`
	RepositoryId  string `json:"repository_id,omitempty"`
	BucketId      string `json:"bucket_id,omitempty"`
	CapabilityId  string `json:"capability_id,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
