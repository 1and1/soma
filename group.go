package proto

type Group struct {
	Id             string        `json:"id,omitempty"`
	Name           string        `json:"name,omitempty"`
	BucketId       string        `json:"bucketId,omitempty"`
	ObjectState    string        `json:"objectState,omitempty"`
	TeamId         string        `json:"teamId,omitempty"`
	MemberGroups   []Group       `json:"memberGroups,omitempty"`
	MemberClusters []Cluster     `json:"memberClusters,omitempty"`
	MemberNodes    []Node        `json:"memberNodes,omitempty"`
	Details        *GroupDetails `json:"details,omitempty"`
	Properties     *[]Property   `json:"properties,omitempty"`
}

type GroupFilter struct {
	Name     string `json:"name,omitempty"`
	BucketId string `json:"bucketId,omitempty"`
}

type GroupDetails struct {
	DetailsCreation
}

//
func (p *Group) DeepCompare(a *Group) bool {
	if a == nil {
		return false
	}
	if p.Id != a.Id || p.Name != a.Name || p.BucketId != a.BucketId ||
		p.ObjectState != a.ObjectState || p.TeamId != a.TeamId {
		return false
	}
	return true
}

func NewGroupRequest() Request {
	return Request{
		Group: &Group{},
	}
}

func NewGroupFilter() Request {
	return Request{
		Filter: &Filter{
			Group: &GroupFilter{},
		},
	}
}

func NewGroupResult() Result {
	return Result{
		Errors: &[]string{},
		Groups: &[]Group{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
