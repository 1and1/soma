package somaproto

type ProtoRequestTeam struct {
	Team   *ProtoTeam       `json:"team,omitempty"`
	Filter *ProtoTeamFilter `json:"filter,omitempty"`
}

type ProtoResultTeam struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Teams  []ProtoTeam `json:"teams,omitempty"`
}

type ProtoTeam struct {
	Id      string            `json:"id,omitempty"`
	Name    string            `json:"name,omitempty"`
	Ldap    string            `json:"ldap,omitempty"`
	System  bool              `json:"system,omitempty"`
	Details *ProtoTeamDetails `json:"details,omitempty"`
}

type ProtoTeamDetails struct {
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	Members   []string `json:"members,omitempty"`
}

type ProtoTeamFilter struct {
	Name   string `json:"name,omitempty"`
	Ldap   string `json:"ldap,omitempty"`
	System bool   `json:"system,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
