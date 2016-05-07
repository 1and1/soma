package somaproto

type Filter struct {
	Bucket      *BucketFilter      `json:"bucket, omitempty"`
	Capability  *CapabilityFilter  `json:"capability, omitempty"`
	CheckConfig *CheckConfigFilter `json:"checkConfig, omitempty"`
	Cluster     *ClusterFilter     `json:"cluster, omitempty"`
	Group       *GroupFilter       `json:"group, omitempty"`
	Level       *LevelFilter       `json:"level, omitempty"`
	Metric      *MetricFilter      `json:"metric, omitempty"`
	Monitoring  *MonitoringFilter  `json:"monitoring, omitempty"`
	Node        *NodeFilter        `json:"node, omitempty"`
	Oncall      *OncallFilter      `json:"oncall, omitempty"`
	Property    *PropertyFilter    `json:"property, omitempty"`
	Provider    *ProviderFilter    `json:"provider, omitempty"`
	Repository  *RepositoryFilter  `json:"repository, omitempty"`
	Server      *ServerFilter      `json:"server, omitempty"`
	Team        *TeamFilter        `json:"team, omitempty"`
	Unit        *UnitFilter        `json:"unit, omitempty"`
	User        *UserFilter        `json:"user, omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
