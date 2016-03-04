package somaproto

type DeploymentDetails struct {
	Repository         string                `json:"repository"`
	Environment        string                `json:"environment"`
	Bucket             string                `json:"bucket"`
	ObjectType         string                `json:"object_type"`
	View               string                `json:"view"`
	Task               string                `json:"task"`
	Datacenter         string                `json:"datacenter"`
	Capability         *ProtoCapability      `json:"capability"`
	Monitoring         *ProtoMonitoring      `json:"monitoring_system"`
	Metric             *ProtoMetric          `json:"metric"`
	Team               *ProtoTeam            `json:"organizational_team"`
	Oncall             *ProtoOncall          `json:"oncall,omitempty"`
	Service            *TreePropertyService  `json:"service,omitempty"`
	Properties         *[]TreePropertySystem `json:"properties,omitempty"`
	CustomProperties   *[]TreePropertyCustom `json:"custom_properties,omitempty"`
	Group              *ProtoGroup           `json:"group,omitempty"`
	Cluster            *ProtoCluster         `json:"cluster,omitempty"`
	Node               *ProtoNode            `json:"node,omitempty"`
	Server             *ProtoServer          `json:"server,omitempty"`
	CheckConfiguration *CheckConfiguration   `json:"check_configuration"`
	Check              *TreeCheck            `json:"check"`
	CheckInstance      *TreeCheckInstance    `json:"check_instance"`
}

func (dd *DeploymentDetails) DeepCompare(alternate *DeploymentDetails) bool {
	if dd.Repository != alternate.Repository {
		return false
	}
	if dd.Environment != alternate.Environment {
		return false
	}
	if dd.Bucket != alternate.Bucket {
		return false
	}
	if dd.ObjectType != alternate.ObjectType {
		return false
	}
	if dd.View != alternate.View {
		return false
	}
	if dd.Task != alternate.Task {
		return false
	}
	if dd.Datacenter != alternate.Datacenter {
		return false
	}
	//if !dd.Capability.DeepCompare(alternate.Capability) {
	//	return false
	//}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
