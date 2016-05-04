package somaproto

type ProtoRequestDatacenter struct {
	Datacenter string `json:"datacenter,omitempty"`
}

type ProtoResultDatacenter struct {
	Code   uint16   `json:"code,omitempty"`
	Status string   `json:"status,omitempty"`
	Text   []string `json:"text,omitempty"`
}

type ProtoResultDatacenterList struct {
	Code        uint16   `json:"code,omitempty"`
	Status      string   `json:"status,omitempty"`
	Text        []string `json:"text,omitempty"`
	Datacenters []string `json:"datacenters,omitempty"`
}

type ProtoResultDatacenterDetail struct {
	Code    uint16                 `json:"code,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Text    []string               `json:"text,omitempty"`
	Details ProtoDatacenterDetails `json:"details,omitempty"`
}

type ProtoDatacenterDetails struct {
	Datacenter string   `json:"datacenter"`
	CreatedAt  string   `json:"createdat,omitempty"`
	CreatedBy  string   `json:"createdby,omitempty"`
	UsedBy     []string `json:"usedby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
