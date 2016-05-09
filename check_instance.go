package proto

type CheckInstance struct {
	InstanceId            string `json:"instanceId,omitempty"`
	CheckId               string `json:"checkId,omitempty"`
	ConfigId              string `json:"configId,omitempty"`
	InstanceConfigId      string `json:"instanceConfigId,omitempty"`
	Version               uint64 `json:"version,omitempty"`
	ConstraintHash        string `json:"constraintHash,omitempty"`
	ConstraintValHash     string `json:"constraintValHash,omitempty"`
	InstanceSvcCfgHash    string `json:"instanceSvcCfghash,omitempty"`
	InstanceService       string `json:"instanceService,omitempty"`
	InstanceServiceConfig string `json:"instanceServiceCfg,omitempty"`
}

func (t *CheckInstance) DeepCompare(a *CheckInstance) bool {
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
