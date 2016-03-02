package somatree


type Action struct {
	Action     string `json:",omitempty"`
	Type       string `json:",omitempty"`
	Repository somaproto.ProtoRepository
	Bucket     somaproto.ProtoBucket
	Group      somaproto.ProtoGroup
	Cluster    somaproto.ProtoCluster
	Node       somaproto.ProtoNode

	ChildType    string `json:",omitempty"`
	ChildGroup   somaproto.ProtoGroup
	ChildCluster somaproto.ProtoCluster
	ChildNode    somaproto.ProtoNode

	Property somaproto.TreeProperty
	Check    somaproto.TreeCheck

	Id              string `json:",omitempty"`
	SourceId        string `json:",omitempty"`
	Name            string `json:",omitempty"`
	Team            string `json:",omitempty"`
	ChildId         string `json:",omitempty"`
	PropertyType    string `json:",omitempty"`
	PropertyId      string `json:",omitempty"`
	PropertySource  string `json:",omitempty"`
	CheckId         string `json:",omitempty"`
	CheckSource     string `json:",omitempty"`
	CheckCapability string `json:",omitempty"`
	InstanceId      string `json:",omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
