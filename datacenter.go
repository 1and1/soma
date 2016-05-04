package somaproto

type Datacenter struct {
	Locode  string             `json:"locode, omitempty"`
	Details *DatacenterDetails `json:"details, omitempty"`
}

type DatacenterDetails struct {
	DetailsCreation
	UsedBy []string `json:"usedBy, omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
