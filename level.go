package somaproto

type Level struct {
	Name      string        `json:"name,omitempty"`
	ShortName string        `json:"shortName,omitempty"`
	Numeric   uint16        `json:"numeric,omitempty"`
	Details   *LevelDetails `json:"details,omitempty"`
}

type LevelFilter struct {
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortName,omitempty"`
	Numeric   uint16 `json:"numeric,omitempty"`
}

type LevelDetails struct {
	DetailsCreation
}

func NewLevelRequest() Request {
	return Request{
		Level: &Level{},
	}
}

func NewLevelFilter() Request {
	return Request{
		Filter: &Filter{
			Level: &LevelFilter{},
		},
	}
}

func NewLevelResult() Result {
	return Result{
		Errors: &[]string{},
		Levels: &[]Level{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
