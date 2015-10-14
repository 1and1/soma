package somaproto

type ProtoRequestObjectType struct {
	Type string `json:"type,omitempty"`
}

type ProtoResultObjectType struct {
	Code   uint16   `json:"code,omitempty"`
	Status string   `json:"status,omitempty"`
	Text   []string `json:"text,omitempty"`
}

type ProtoResultObjectTypeList struct {
	Code   uint16   `json:"code,omitempty"`
	Status string   `json:"status,omitempty"`
	Text   []string `json:"text,omitempty"`
	Types  []string `json:"types,omitempty"`
}

type ProtoResultObjectTypeDetail struct {
	Code    uint16                 `json:"code,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Text    []string               `json:"text,omitempty"`
	Details ProtoObjectTypeDetails `json:"details,omitempty"`
}

type ProtoObjectTypeDetails struct {
	Type      string   `json:"type,omitempty"`
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	UsedBy    []string `json:"usedby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
