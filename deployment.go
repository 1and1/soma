package proto

type Deployment struct {
	Repository       string            `json:"repository"`
	Environment      string            `json:"environment"`
	Bucket           string            `json:"bucket"`
	ObjectType       string            `json:"objectType"`
	View             string            `json:"view"`
	Task             string            `json:"task"`
	Datacenter       string            `json:"datacenter"`
	Capability       *Capability       `json:"capability"`
	Monitoring       *Monitoring       `json:"monitoringSystem"`
	Metric           *Metric           `json:"metric"`
	Unit             *Unit             `json:"unit"`
	Team             *Team             `json:"organizationalTeam"`
	Oncall           *Oncall           `json:"oncallDuty,omitempty"`
	Service          *PropertyService  `json:"service,omitempty"`
	Properties       *[]PropertySystem `json:"properties,omitempty"`
	CustomProperties *[]PropertyCustom `json:"customProperties,omitempty"`
	Group            *Group            `json:"group,omitempty"`
	Cluster          *Cluster          `json:"cluster,omitempty"`
	Node             *Node             `json:"node,omitempty"`
	Server           *Server           `json:"server,omitempty"`
	CheckConfig      *CheckConfig      `json:"checkConfig"`
	Check            *Check            `json:"check"`
	CheckInstance    *CheckInstance    `json:"checkInstance"`
}

func (dd *Deployment) DeepCompare(alternate *Deployment) bool {
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
	// not equal if one of them is nil while the other is not
	if (dd.Oncall == nil && alternate.Oncall != nil) ||
		(dd.Oncall != nil && alternate.Oncall == nil) {
		return false
	}
	// do not compare if any of them are nil
	if !(dd.Oncall == nil || alternate.Oncall == nil) {
		if !dd.Oncall.DeepCompare(alternate.Oncall) {
			return false
		}
	}
	if !dd.Service.DeepCompare(alternate.Service) {
		return false
	}
	if dd.Properties != nil && *dd.Properties != nil {
	proploop:
		for _, pr := range *dd.Properties {
			if pr.DeepCompareSlice(alternate.Properties) {
				continue proploop
			}
			return false
		}
	}
	if alternate.Properties != nil && *alternate.Properties != nil {
	revproploop:
		for _, pr := range *alternate.Properties {
			if pr.DeepCompareSlice(dd.Properties) {
				continue revproploop
			}
			return false
		}
	}
	if dd.CustomProperties != nil && *dd.CustomProperties != nil {
	cproploop:
		for _, pr := range *dd.CustomProperties {
			if pr.DeepCompareSlice(alternate.CustomProperties) {
				continue cproploop
			}
			return false
		}
	}
	if alternate.CustomProperties != nil && *alternate.CustomProperties != nil {
	revcproploop:
		for _, pr := range *alternate.CustomProperties {
			if pr.DeepCompareSlice(dd.CustomProperties) {
				continue revcproploop
			}
			return false
		}
	}
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
	if !dd.CheckConfig.DeepCompare(alternate.CheckConfig) {
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

func NewDeploymentResult() Result {
	return Result{
		Errors:      &[]string{},
		Deployments: &[]Deployment{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
