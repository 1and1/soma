package proto

type Cluster struct {
	Id          string      `json:"id,omitempty"`
	Name        string      `json:"name,omitempty"`
	BucketId    string      `json:"bucketId,omitempty"`
	ObjectState string      `json:"objectState,omitempty"`
	TeamId      string      `json:"teamId,omitempty"`
	Members     *[]Node     `json:"members,omitempty"`
	Details     *Details    `json:"details,omitempty"`
	Properties  *[]Property `json:"properties,omitempty"`
}

type ClusterFilter struct {
	Name     string `json:"name,omitempty"`
	BucketId string `json:"bucketid,omitempty"`
	TeamId   string `json:"teamId,omitempty"`
}

//
func (c *Cluster) DeepCompare(a *Cluster) bool {
	if a == nil {
		return false
	}

	if c.Id != a.Id ||
		c.Name != a.Name ||
		c.BucketId != a.BucketId ||
		c.ObjectState != a.ObjectState ||
		c.TeamId != a.TeamId {
		return false
	}

	if c.Members != nil && a.Members != nil {
	member:
		for i, _ := range *c.Members {
			for j, _ := range *a.Members {
				if (*c.Members)[i].Id == (*a.Members)[j].Id {
					continue member
				}
			}
			return false
		}
		return true
	} else if c.Members != nil && a.Members == nil {
		return false
	} else if c.Members == nil && a.Members != nil {
		return false
	}
	return true
}

func NewClusterRequest() Request {
	return Request{
		Flags:   &Flags{},
		Cluster: &Cluster{},
	}
}

func NewClusterFilter() Request {
	return Request{
		Filter: &Filter{
			Cluster: &ClusterFilter{},
		},
	}
}

func NewClusterResult() Result {
	return Result{
		Errors:   &[]string{},
		Clusters: &[]Cluster{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
