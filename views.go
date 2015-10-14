package somaproto

type ProtoRequestView struct {
	View string `json:"view,omitempty"`
}

type ProtoResultView struct {
	Code   uint16   `json:"code,omitempty"`
	Status string   `json:"status,omitempty"`
	Text   []string `json:"text,omitempty"`
}

type ProtoResultViewList struct {
	Code   uint16   `json:"code,omitempty"`
	Status string   `json:"status,omitempty"`
	Text   []string `json:"text,omitempty"`
	Views  []string `json:"views,omitempty"`
}

type ProtoResultViewDetail struct {
	Code    uint16           `json:"code,omitempty"`
	Status  string           `json:"status,omitempty"`
	Text    []string         `json:"text,omitempty"`
	Details ProtoViewDetails `json:"details,omitempty"`
}

type ProtoViewDetails struct {
	View      string   `json:"view,omitempty"`
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	UsedBy    []string `json:"usedby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
