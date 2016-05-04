package somaproto

type Oncall struct {
	Id      string         `json:"id, omitempty"`
	Name    string         `json:"name, omitempty"`
	Number  string         `json:"number, omitempty"`
	Details *OncallDetails `json:"details, omitempty"`
}

type OncallDetails struct {
	DetailsCreation
	Members *[]OncallMember `json:"members, omitempty"`
}

type OncallMember struct {
	UserName string `json:"userName, omitempty"`
	UserId   string `json"userId, omitempty"`
}

type OncallFilter struct {
	Name   string `json:"name, omitempty"`
	Number string `json:"number, omitempty"`
}

//
func (p *Oncall) DeepCompare(a *Oncall) bool {
	if p.Id != a.Id || p.Name != a.Name || p.Number != a.Number {
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
