package somaproto

type ProtoRequestStatus struct {
	Status ProtoStatus `json:"status,omitempty"`
}

type ProtoResultStatus struct {
	Code       uint16        `json:"code,omitempty"`
	Status     string        `json:"status,omitempty"`
	Text       []string      `json:"text,omitempty"`
	StatusList []ProtoStatus `json:"statuslist,omitempty"`
	JobId      string        `json:"jobid,omitempty"`
}

type ProtoStatus struct {
	Status  string              `json:"status,omitempty"`
	Details *ProtoStatusDetails `json:"details,omitempty"`
}

type ProtoStatusDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
