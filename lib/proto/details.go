package proto

type Details struct {
	CreatedAt    string         `json:"createdAt,omitempty"`
	CreatedBy    string         `json:"createdBy,omitempty"`
	Server       Server         `json:"server,omitempty"`
	CheckConfigs *[]CheckConfig `json:"checkConfigs,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
