package somaproto

type Node struct {
	Id         string       `json:"id,omitempty"`
	AssetId    uint64       `json:"assetId,omitempty"`
	Name       string       `json:"name,omitempty"`
	TeamId     string       `json:"teamId,omitempty"`
	ServerId   string       `json:"serverId,omitempty"`
	State      string       `json:"state,omitempty"`
	IsOnline   bool         `json:"isOnline,omitempty"`
	IsDeleted  bool         `json:"isDeleted,omitempty"`
	Details    *NodeDetails `json:"details,omitempty"`
	Config     *NodeConfig  `json:"config,omitempty"`
	Properties *[]Property  `json:"properties,omitempty"`
}

type NodeDetails struct {
	DetailsCreation
	Server Server `json:"server,omitempty"`
}

type NodeFilter struct {
	Name       string `json:"name,omitempty"`
	TeamId     string `json:"teamId,omitempty"`
	ServerId   string `json:"serverId,omitempty"`
	IsOnline   bool   `json:"isOnline,omitempty"`
	NotOnline  bool   `json:"notOnline,omitempty"`
	Deleted    bool   `json:"isDeleted,omitempty"`
	NotDeleted bool   `json:"notDeleted,omitempty"`
	/*
		PropertyType  string `json:"propertytype,omitempty"`
		Property      string `json:"property,omitempty"`
		LocalProperty bool   `json:"localproperty,omitempty"`
	*/
}

type NodeConfig struct {
	RepositoryId string `json:"repositoryId,omitempty"`
	BucketId     string `json:"bucketId,omitempty"`
}

//
func (p *Node) DeepCompare(a *Node) bool {
	if a == nil {
		return false
	}

	if p.Id != a.Id || p.AssetId != a.AssetId || p.Name != a.Name ||
		p.TeamId != a.TeamId || p.ServerId != a.ServerId || p.State != a.State ||
		p.IsOnline != a.IsOnline || p.IsDeleted != a.IsDeleted {
		return false
	}
	return true
}

func NewNodeRequest() Request {
	return Request{
		Node: &Node{},
	}
}

func NewNodeFilter() Request {
	return Request{
		Filter: &Filter{
			Node: &NodeFilter{},
		},
	}
}

func NewNodeResult() Result {
	return Result{
		Errors: &[]string{},
		Nodes:  &[]Node{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
