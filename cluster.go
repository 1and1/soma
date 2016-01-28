package somaproto

type ProtoRequestCluster struct {
	Cluster ProtoCluster       `json:"cluster,omitempty"`
	Filter  ProtoClusterFilter `json:"filter,omitempty"`
}

type ProtoResultCluster struct {
	Code     uint16         `json:"code,omitempty"`
	Status   string         `json:"status,omitempty"`
	Text     []string       `json:"text,omitempty"`
	Clusters []ProtoCluster `json:"clusters,omitempty"`
	JobId    string         `json:"jobid,omitempty"`
}

type ProtoCluster struct {
	Id          string                 `json:"id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Bucket      string                 `json:"bucket,omitempty"`
	BucketId    string                 `json:"bucketid,omitempty"`
	ObjectState string                 `json:"objectstate,omitempty"`
	Members     []ProtoNode            `json:"members,omitempty"`
	Details     *ProtoClusterDetails   `json:"details,omitempty"`
	Properties  []ProtoClusterProperty `json:"properties,omitempty"`
}

type ProtoClusterFilter struct {
	Name     string `json:"name,omitempty"`
	Bucket   string `json:"bucket,omitempty"`
	BucketId string `json:"bucketid,omitempty"`
}

type ProtoClusterDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

type ProtoClusterProperty struct {
	Type         string `json:"type,omitempty"`
	View         string `json:"view,omitempty"`
	Property     string `json:"property,omitempty"`
	Value        string `json:"value,omitempty"`
	Inheritance  bool   `json:"inheritance,omitempty"`
	ChildrenOnly bool   `json:"children,omitempty"`
	Source       string `json:"source,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
