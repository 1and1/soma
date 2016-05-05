package somaproto

type DeploymentDetailsResult struct {
	Code        uint16              `json:"code,omitempty"`
	Status      string              `json:"status,omitempty"`
	Deployments []DeploymentDetails `json:"deployments,omitempty"`
	List        []string            `json:"list,omitempty"`
	JobId       string              `json:"jobid,omitempty"`
}

func (dd *DeploymentDetailsResult) ErrorMark(err error, imp bool,
	found bool, length int, jobid string) bool {
	if dd.markError(err) {
		return true
	}
	if dd.markImplemented(imp) {
		return true
	}
	if dd.markFound(found, length) {
		return true
	}
	if dd.hasJobId(jobid) {
		return dd.markAccepted()
	}
	return dd.markOk()
}

func (dd *DeploymentDetailsResult) markError(err error) bool {
	if err != nil {
		dd.Code = 500
		dd.Status = "ERROR"
		return true
	}
	return false
}

func (dd *DeploymentDetailsResult) markImplemented(f bool) bool {
	if f {
		dd.Code = 501
		dd.Status = "NOT IMPLEMENTED"
		return true
	}
	return false
}

func (dd *DeploymentDetailsResult) markFound(f bool, i int) bool {
	if f || i == 0 {
		dd.Code = 404
		dd.Status = "NOT FOUND"
		return true
	}
	return false
}

func (dd *DeploymentDetailsResult) markOk() bool {
	dd.Code = 200
	dd.Status = "OK"
	return false
}

func (dd *DeploymentDetailsResult) hasJobId(s string) bool {
	if s != "" {
		dd.JobId = s
		return true
	}
	return false
}

func (dd *DeploymentDetailsResult) markAccepted() bool {
	dd.Code = 202
	dd.Status = "ACCEPTED"
	return false
}

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
	Unit               *ProtoUnit            `json:"unit"`
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
	//
	if !dd.Capability.DeepCompare(alternate.Capability) {
		return false
	}
	if !dd.Monitoring.DeepCompare(alternate.Monitoring) {
		return false
	}
	if !dd.Metric.DeepCompare(alternate.Metric) {
		return false
	}
	if !dd.Unit.DeepCompare(alternate.Unit) {
		return false
	}
	if !dd.Team.DeepCompare(alternate.Team) {
		return false
	}
	if !dd.Oncall.DeepCompare(alternate.Oncall) {
		return false
	}
	if !dd.Service.DeepCompare(alternate.Service) {
		return false
	}
	// TODO: Properties
	// TODO: CustomProperties
	if dd.Group != nil && !dd.Group.DeepCompare(alternate.Group) {
		return false
	} else if dd.Group == nil && alternate.Group != nil {
		return false
	}
	if dd.Cluster != nil && !dd.Cluster.DeepCompare(alternate.Cluster) {
		return false
	} else if dd.Cluster == nil && alternate.Cluster != nil {
		return false
	}
	if dd.Node != nil && !dd.Node.DeepCompare(alternate.Node) {
		return false
	} else if dd.Node == nil && alternate.Node != nil {
		return false
	}
	if dd.Server != nil && !dd.Server.DeepCompare(alternate.Server) {
		return false
	} else if dd.Server == nil && alternate.Server != nil {
		return false
	}
	if !dd.CheckConfiguration.DeepCompare(alternate.CheckConfiguration) {
		return false
	}
	if !dd.Check.DeepCompare(alternate.Check) {
		return false
	}
	if !dd.CheckInstance.DeepCompare(alternate.CheckInstance) {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
