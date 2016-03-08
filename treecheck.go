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

func (t *TreeCheck) DeepCompare(a *TreeCheck) bool {
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

func (t *TreeCheckInstance) DeepCompare(a *TreeCheckInstance) bool {
	if t.InstanceId != a.InstanceId || t.CheckId != a.CheckId || t.ConfigId != a.ConfigId ||
		t.ConstraintHash != a.ConstraintHash || t.ConstraintValHash != a.ConstraintValHash ||
		t.InstanceSvcCfgHash != a.InstanceSvcCfgHash || t.InstanceService != a.InstanceService {
		// - InstanceConfigId is a randomly generated uuid on every instance calculation
		// - Version is incremented on every instance calculation
		// - InstanceServiceConfig is compared as deploymentdetails.Service
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
