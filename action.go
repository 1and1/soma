package somatree


type Action struct {
	Action        string              `json:"action,omitempty"`
	Type          string              `json:"type,omitempty"`
	Bucket        proto.Bucket        `json:"bucket,omitempty"`
	Check         proto.Check         `json:"check,omitempty"`
	CheckInstance proto.CheckInstance `json:"check_instance,omitempty"`
	ChildCluster  proto.Cluster       `json:"child_cluster,omitempty"`
	ChildGroup    proto.Group         `json:"child_group,omitempty"`
	ChildNode     proto.Node          `json:"child_node,omitempty"`
	ChildType     string              `json:"child_type,omitempty"`
	Cluster       proto.Cluster       `json:"cluster,omitempty"`
	Group         proto.Group         `json:"group,omitempty"`
	Node          proto.Node          `json:"node,omitempty"`
	Property      proto.Property      `json:"property,omitempty"`
	Repository    proto.Repository    `json:"repository,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
