package proto

type Unit struct {
	Unit    string       `json:"unit,omitempty"`
	Name    string       `json:"name,omitempty"`
	Details *UnitDetails `json:"details,omitempty"`
}

type UnitFilter struct {
	Unit string `json:"unit,omitempty"`
	Name string `json:"name,omitempty"`
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

func NewUnitRequest() Request {
	return Request{
		Unit: &Unit{},
	}
}

func NewUnitFilter() Request {
	return Request{
		Filter: &Filter{
			Unit: &UnitFilter{},
		},
	}
}

func NewUnitResult() Result {
	return Result{
		Errors: &[]string{},
		Units:  &[]Unit{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
