package somaproto

import (
	"github.com/satori/go.uuid"
)

type ProtoRequestTeam struct {
	Team   ProtoTeam       `json:"team,omitempty"`
	Filter ProtoTeamFilter `json:"filter,omitempty"`
}

type ProtoResultTeam struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Teams  []ProtoTeam `json:"teams,omitempty"`
}

type ProtoTeam struct {
	TeamId   uuid.UUID        `json:"teamid,omitempty"`
	TeamName string           `json:"teamname,omitempty"`
	LdapId   string           `json:"ldapid,omitempty"`
	System   bool             `json:"system,omitempty"`
	Details  ProtoTeamDetails `json:"details,omitempty"`
}

type ProtoTeamDetails struct {
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	Members   []string `json:"members,omitempty"`
}

type ProtoTeamFilter struct {
	TeamName string `json:"teamname,omitempty"`
	LdapId   string `json:"ldapid,omitempty"`
	System   bool   `json:"system,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
