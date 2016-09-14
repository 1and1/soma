package proto

type Bucket struct {
	Id             string      `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	RepositoryId   string      `json:"repositoryId,omitempty"`
	TeamId         string      `json:"teamId,omitempty"`
	Environment    string      `json:"environment,omitempty"`
	IsDeleted      bool        `json:"isDeleted,omitempty"`
	IsFrozen       bool        `json:"isFrozen,omitempty"`
	MemberGroups   *[]Group    `json:"memberGroups,omitempty"`
	MemberClusters *[]Cluster  `json:"memberClusters,omitempty"`
	MemberNodes    *[]Node     `json:"memberNodes,omitempty"`
	Details        *Details    `json:"details,omitempty"`
	Properties     *[]Property `json:"properties,omitempty"`
}

type BucketFilter struct {
	Name         string `json:"name,omitempty"`
	Id           string `json:"id,omitempty"`
	RepositoryId string `json:"repositoryId,omitempty"`
	IsDeleted    bool   `json:"isDeleted,omitempty"`
	IsFrozen     bool   `json:"isFrozen,omitempty"`
}

func NewBucketRequest() Request {
	return Request{
		Flags:  &Flags{},
		Bucket: &Bucket{},
	}
}

func NewBucketFilter() Request {
	return Request{
		Filter: &Filter{
			Bucket: &BucketFilter{},
		},
	}
}

func NewBucketResult() Result {
	return Result{
		Errors:  &[]string{},
		Buckets: &[]Bucket{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
