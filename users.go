package somaproto

import (
	"github.com/satori/go.uuid"
)

type ProtoRequestUser struct {
	User        ProtoUser            `json:"user,omitempty"`
	Credentials ProtoUserCredentials `json:"credentials,omitempty"`
	Filter      ProtoUserFilter      `json:"filter,omitempty"`
	Restore     bool                 `json:"restore,omitempty"`
	Purge       bool                 `json:"purge,omitempty"`
}

type ProtoResultUser struct {
	Code   uint16      `json:"code,omitempty"`
	Status string      `json:"status,omitempty"`
	Text   []string    `json:"text,omitempty"`
	Users  []ProtoUser `json:"users,omitempty"`
}

type ProtoUser struct {
	Id             uuid.UUID         `json:"id,omitempty"`
	UserName       string            `json:"username,omitempty"`
	FirstName      string            `json:"firstname,omitempty"`
	LastName       string            `json:"lastname,omitempty"`
	EmployeeNumber string            `json:"employeenumber,omitempty"`
	MailAddress    string            `json:"mailaddress,omitempty"`
	IsActive       bool              `json:"active,omitempty"`
	IsSystem       bool              `json:"system,omitempty"`
	IsDeleted      bool              `json:"deleted,omitempty"`
	Team           string            `json:"team,omitempty"`
	Details        *ProtoUserDetails `json:"details,omitempty"`
}

type ProtoUserCredentials struct {
	Reset    bool   `json:"reset,omitempty"`
	Force    bool   `json:"force,omitempty"`
	Password string `json:"password,omitempty"`
}

type ProtoUserDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

type ProtoUserFilter struct {
	UserName  string `json:"username,omitempty"`
	IsActive  bool   `json:"active,omitempty"`
	IsSystem  bool   `json:"system,omitempty"`
	IsDeleted bool   `json:"deleted,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
