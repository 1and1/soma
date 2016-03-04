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

type TreeCheckInstance struct {
	InstanceId            string `json:"instance_id,omitempty"`
	CheckId               string `json:"check_id,omitempty"`
	ConfigId              string `json:"config_id,omitempty"`
	InstanceConfigId      string `json:"instance_config_id,omitempty"`
	Version               uint64 `json:"version,omitempty"`
	ConstraintHash        string `json:"constraint_hash,omitempty"`
	ConstraintValHash     string `json:"constraint_val_hash,omitempty"`
	InstanceSvcCfgHash    string `json:"instance_svc_cfghash,omitempty"`
	InstanceService       string `json:"instance_service,omitempty"`
	InstanceServiceConfig string `json:"instance_service_cfg,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
