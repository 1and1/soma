package somaproto

type ProtoRequestObjectState struct {
	State string `json:"state,omitempty"`
}

type ProtoResultObjectState struct {
	Code   uint16   `json:"code,omitempty"`
	Status string   `json:"status,omitempty"`
	Text   []string `json:"text,omitempty"`
}

type ProtoResultObjectStateList struct {
	Code   uint16   `json:"code,omitempty"`
	Status string   `json:"status,omitempty"`
	Text   []string `json:"text,omitempty"`
	States []string `json:"states,omitempty"`
}

type ProtoResultObjectStateDetail struct {
	Code    uint16                  `json:"code,omitempty"`
	Status  string                  `json:"status,omitempty"`
	Text    []string                `json:"text,omitempty"`
	Details ProtoObjectStateDetails `json:"details,omitempty"`
}

type ProtoObjectStateDetails struct {
	State     string   `json:"state,omitempty"`
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	UsedBy    []string `json:"usedby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
