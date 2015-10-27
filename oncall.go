package somaproto

import (
	"github.com/satori/go.uuid"
)

type ProtoRequestOncall struct {
	OnCall  ProtoOncall         `json:"oncall,omitempty"`
	Filter  ProtoOncallFilter   `json:"filter,omitempty"`
	Members []ProtoOncallMember `json:"members,omitempty"`
}

type ProtoResultOncall struct {
	Code    uint16        `json:"code,omitempty"`
	Status  string        `json:"status,omitempty"`
	Text    []string      `json:"text,omitempty"`
	Oncalls []ProtoOncall `json:"oncalls,omitempty"`
}

type ProtoOncall struct {
	Id      uuid.UUID          `json:"id,omitempty"`
	Name    string             `json:"name,omitempty"`
	Number  string             `json:"number,omitempty"`
	Details ProtoOncallDetails `json:"details,omitempty"`
}

type ProtoOncallDetails struct {
	CreatedAt string   `json:"createdat,omitempty"`
	CreatedBy string   `json:"createdby,omitempty"`
	Members   []string `json:"members,omitempty"`
}

type ProtoOncallMember struct {
	UserName string    `json:"username,omitempty"`
	UserId   uuid.UUID `json"userid,omitempty"`
}

type ProtoOncallFilter struct {
	Name   string `json:"name,omitempty"`
	Number string `json:"number,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
