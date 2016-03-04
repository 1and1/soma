package somatree


type Action struct {
	Action        string                      `json:"action,omitempty"`
	Type          string                      `json:"type,omitempty"`
	Repository    somaproto.ProtoRepository   `json:"repository,omitempty"`
	Bucket        somaproto.ProtoBucket       `json:"bucket,omitempty"`
	Group         somaproto.ProtoGroup        `json:"group,omitempty"`
	Cluster       somaproto.ProtoCluster      `json:"cluster,omitempty"`
	Node          somaproto.ProtoNode         `json:"node,omitempty"`
	Property      somaproto.TreeProperty      `json:"property,omitempty"`
	Check         somaproto.TreeCheck         `json:"check,omitempty"`
	CheckInstance somaproto.TreeCheckInstance `json:"check_instance,omitempty"`
	ChildType     string                      `json:"child_type,omitempty"`
	ChildGroup    somaproto.ProtoGroup        `json:"child_group,omitempty"`
	ChildCluster  somaproto.ProtoCluster      `json:"child_cluster,omitempty"`
	ChildNode     somaproto.ProtoNode         `json:"child_node,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
