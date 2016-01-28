package somaproto

type ProtoRequestGroup struct {
	Group  ProtoGroup       `json:"group,omitempty"`
	Filter ProtoGroupFilter `json:"filter,omitempty"`
}

type ProtoResultGroup struct {
	Code   uint16       `json:"code,omitempty"`
	Status string       `json:"status,omitempty"`
	Text   []string     `json:"text,omitempty"`
	Groups []ProtoGroup `json:"group,omitempty"`
	JobId  string       `json:"jobid,omitempty"`
}

type ProtoGroup struct {
	Id             string               `json:"id,omitempty"`
	Name           string               `json:"name,omitempty"`
	Bucket         string               `json:"bucket,omitempty"`
	BucketId       string               `json:"bucketid,omitempty"`
	ObjectState    string               `json:"objectstate,omitempty"`
	MemberGroups   []ProtoGroup         `json:"membergroups,omitempty"`
	MemberClusters []ProtoCluster       `json:"memberclusters,omitempty"`
	MemberNodes    []ProtoNode          `json:"membernodes,omitempty"`
	Details        *ProtoGroupDetails   `json:"details,omitempty"`
	Properties     []ProtoGroupProperty `json:"properties,omitempty"`
}

type ProtoGroupFilter struct {
	Name     string `json:"name,omitempty"`
	Bucket   string `json:"bucket,omitempty"`
	BucketId string `json:"bucketid,omitempty"`
}

type ProtoGroupDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

type ProtoGroupProperty struct {
	Type         string `json:"type,omitempty"`
	View         string `json:"view,omitempty"`
	Property     string `json:"property,omitempty"`
	Value        string `json:"value,omitempty"`
	Inheritance  bool   `json:"inheritance,omitempty"`
	ChildrenOnly bool   `json:"children,omitempty"`
	Source       string `json:"source,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
