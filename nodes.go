package somaproto

type ProtoRequestNode struct {
	Node    *ProtoNode       `json:"node,omitempty"`
	Filter  *ProtoNodeFilter `json:"filter,omitempty"`
	Restore bool             `json:"restore,omitempty"`
	Purge   bool             `json:"purge,omitempty"`
}

type ProtoResultNode struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Nodes  []ProtoNode `json:"nodes,omitempty"`
	JobId  string      `json:"jobid,omitempty"`
}

type ProtoNode struct {
	Id         string              `json:"id,omitempty"`
	AssetId    uint64              `json:"assetid,omitempty"`
	Name       string              `json:"name,omitempty"`
	Team       string              `json:"team,omitempty"`
	Server     string              `json:"server,omitempty"`
	State      string              `json:"state,omitempty"`
	IsOnline   bool                `json:"online,omitempty"`
	IsDeleted  bool                `json:"deleted,omitempty"`
	Details    *ProtoNodeDetails   `json:"details,omitempty"`
	Config     *ProtoNodeConfig    `json:"config,omitempty"`
	Properties []ProtoNodeProperty `json:"properties,omitempty"`
}

type ProtoNodeDetails struct {
	CreatedAt string      `json:"createdat,omitempty"`
	CreatedBy string      `json:"createdby,omitempty"`
	Server    ProtoServer `json:"server,omitempty"`
}

type ProtoNodeFilter struct {
	Name          string `json:"name,omitempty"`
	Team          string `json:"team,omitempty"`
	Server        string `json:"server,omitempty"`
	IsOnline      bool   `json:"online,omitempty"`
	IsDeleted     bool   `json:"deleted,omitempty"`
	PropertyType  string `json:"propertytype,omitempty"`
	Property      string `json:"property,omitempty"`
	LocalProperty bool   `json:"localproperty,omitempty"`
}

type ProtoNodeConfig struct {
	RepositoryId   string `json:"repositoryid,omitempty"`
	RepositoryName string `json:"repositoryname,omitempty"`
	BucketId       string `json:"bucketid,omitempty"`
	BucketName     string `json:"bucketname,omitempty"`
}

type ProtoNodeProperty struct {
	Type         string `json:"type,omitempty"`
	View         string `json:"view,omitempty"`
	Property     string `json:"property,omitempty"`
	Value        string `json:"value,omitempty"`
	Inheritance  bool   `json:"inheritance,omitempty"`
	ChildrenOnly bool   `json:"children,omitempty"`
	Source       string `json:"source,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
