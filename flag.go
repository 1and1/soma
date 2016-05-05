package somaproto

type Flags struct {
	Restore  bool `json:"restore"`
	Purge    bool `json:"purge"`
	Freeze   bool `json:"freeze"`
	Thaw     bool `json:"thaw"`
	Clear    bool `json:"clear"`    // repository
	Activate bool `json:"activate"` // repository
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
