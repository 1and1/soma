package somaproto

type Level struct {
	Name      string        `json:"name, omitempty"`
	ShortName string        `json:"shortName, omitempty"`
	Numeric   uint16        `json:"numeric, omitempty"`
	Details   *LevelDetails `json:"details, omitempty"`
}

type LevelFilter struct {
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortName,omitempty"`
	Numeric   uint16 `json:"numeric,omitempty"`
}

type LevelDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
