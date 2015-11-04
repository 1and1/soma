package somaproto

import (
	"github.com/satori/go.uuid"
)

type ProtoRequestRepository struct {
	Filter ProtoRepositoryFilter `json:"filter,omitempty"`
}

type ProtoResultRepository struct {
	Code         uint16            `json:"code,omitempty"`
	Status       string            `json:"status,omitempty"`
	Text         []string          `json:"text,omitempty"`
	Repositories []ProtoRepository `json:"repositories,omitempty"`
}

type ProtoRepository struct {
	Id   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}

type ProtoRepositoryFilter struct {
	Name string `json:"name,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
