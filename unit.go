package somaproto

type Unit struct {
	Unit    string       `json:"unit, omitempty"`
	Name    string       `json:"name, omitempty"`
	Details *UnitDetails `json:"details, omitempty"`
}

type UnitFilter struct {
	Unit string `json:"unit, omitempty"`
	Name string `json:"name, omitempty"`
}

type UnitDetails struct {
	DetailsCreation
}

//
func (p *Unit) DeepCompare(a *Unit) bool {
	if p.Unit != a.Unit || p.Name != a.Name {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
