package somaproto

type ProtoRequestServer struct {
	Server ProtoServer       `json:"server,omitempty"`
	Filter ProtoServerFilter `json:"filter,omitempty"`
	Purge  bool              `json:"purge,omitempty"`
}

type ProtoResultServer struct {
	Code    uint16        `json:"code,omitempty"`
	Status  string        `json:"status,omitempty"`
	Text    []string      `json:"text,omitempty"`
	Servers []ProtoServer `json:"servers,omitempty"`
}

type ProtoServer struct {
	AssetId    uint64             `json:"assetid,omitempty"`
	Datacenter string             `json:"datacenter,omitempty"`
	Location   string             `json:"location,omitempty"`
	Name       string             `json:"name,omitempty"`
	Online     bool               `json:"online,omitempty"`
	Details    ProtoServerDetails `json:"details,omitempty"`
}

type ProtoServerDetails struct {
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	Nodes     []string `json:"nodes,omitempty"`
}

type ProtoServerFilter struct {
	Online     bool   `json:"online,omitempty"`
	Deleted    bool   `json:"deleted,omitempty"`
	Datacenter string `json:"datacenter,omitempty"`
	Name       string `json:"name,omitempty"`
}
