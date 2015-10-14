package somaproto

type ProtoRequestDatacenterGroup struct {
	DatacenterGroup string `json:"datacentergroup,omitempty"`
}

type ProtoResultDatacenterGroup struct {
	Code   uint16   `json:"code,omitempty"`
	Status string   `json:"status,omitempty"`
	Text   []string `json:"text,omitempty"`
}

type ProtoResultDatacenterGroupList struct {
	Code             uint16   `json:"code,omitempty"`
	Status           string   `json:"status,omitempty"`
	Text             []string `json:"text,omitempty"`
	DatacenterGroups []string `json:"datacentergroups,omitempty"`
}

type ProtoResultDatacenterGroupDetail struct {
	Code    uint16                      `json:"code,omitempty"`
	Status  string                      `json:"status,omitempty"`
	Text    []string                    `json:"text,omitempty"`
	Details ProtoDatacenterGroupDetails `json:"details,omitempty"`
}

type ProtoDatacenterGroupDetails struct {
	DatacenterGroup string   `json:"datacentergroup"`
	CreatedAt       string   `json:"createdat,omitempty"`
	CreatedBy       string   `json:"createdby,omitempty"`
	UsedBy          []string `json:"usedby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
