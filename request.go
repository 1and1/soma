package somaproto

type Request struct {
	Filter *Filter `json:"filter, omitempty"`
	Flags  *Flags  `json:"flags, omitempty"`

	Attribute       *Attribute       `json:"attribute, omitempty"`
	Bucket          *Bucket          `json:"bucket, omitempty"`
	Capability      *Capability      `json:"capability, omitempty"`
	Category        *Category        `json:"category, omitempty"`
	CheckConfig     *CheckConfig     `json:"checkConfig, omitempty"`
	Cluster         *Cluster         `json:"cluster, omitempty"`
	Datacenter      *Datacenter      `json:"datacenter, omitempty"`
	DatacenterGroup *DatacenterGroup `json:"datacenterGroup, omitempty"`
	Entity          *Entity          `json:"entity, omitempty"`
	Environment     *Environment     `json:"environment, omitempty"`
	Grant           *Grant           `json:"grant, omitempty"`
	Group           *Group           `json:"group, omitempty"`
	HostDeployment  *HostDeployment  `json:"hostDeployment, omitempty"`
	Level           *Level           `json:"level, omitempty"`
	Metric          *Metric          `json:"metric, omitempty"`
	Mode            *Mode            `json:"mode, omitempty"`
	Monitoring      *Monitoring      `json:"monitoring, omitempty"`
	Node            *Node            `json:"node, omitempty"`
	Oncall          *Oncall          `json:"oncall, omitempty"`
	Permission      *Permission      `json:"permission, omitempty"`
	Predicate       *Predicate       `json:"predicate, omitempty"`
	Property        *Property        `json:"property, omitempty"`
	Provider        *Provider        `json:"provider, omitempty"`
	Repository      *Repository      `json:"repository, omitempty"`
	Server          *Server          `json:"server, omitempty"`
	State           *State           `json:"state, omitempty"`
	Status          *Status          `json:"status, omitempty"`
	Team            *Team            `json:"team, omitempty"`
	Unit            *Unit            `json:"unit, omitempty"`
	User            *User            `json:"user, omitempty"`
	Validity        *Validity        `json:"validity, omitempty"`
	View            *View            `json:"view, omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
