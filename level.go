package somaproto

type ProtoRequestLevel struct {
	Level  ProtoLevel       `json:"level,omitempty"`
	Filter ProtoLevelFilter `json:"filter,omitempty"`
}

type ProtoResultLevel struct {
	Code   uint16       `json:"code,omitempty"`
	Status string       `json:"status,omitempty"`
	Text   []string     `json:"text,omitempty"`
	Levels []ProtoLevel `json:"levels,omitempty"`
	JobId  string       `json:"jobid,omitempty"`
}

type ProtoLevel struct {
	Name      string            `json:"name,omitempty"`
	ShortName string            `json:"shortname,omitempty"`
	Numeric   uint16            `json:"numeric,omitempty"`
	Details   ProtoLevelDetails `json:"details,omitempty"`
}

type ProtoLevelFilter struct {
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortname,omitempty"`
	Numeric   uint16 `json:"numeric,omitempty"`
}

type ProtoLevelDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
