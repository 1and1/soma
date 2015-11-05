package somaproto

import (
	"github.com/satori/go.uuid"
)

type ProtoRequestRepository struct {
	Repository ProtoRepository       `json:"repository,omitempty"`
	Filter     ProtoRepositoryFilter `json:"filter,omitempty"`
	Restore    bool                  `json:"restore,omitempty"`
	Purge      bool                  `json:"purge,omitempty"`
	Clear      bool                  `json:"clear,omitempty"`
	Activate   bool                  `json:"activate,omitempty"`
}

type ProtoResultRepository struct {
	Code         uint16            `json:"code,omitempty"`
	Status       string            `json:"status,omitempty"`
	Text         []string          `json:"text,omitempty"`
	Repositories []ProtoRepository `json:"repositories,omitempty"`
	JobId        uuid.UUID         `json:"jobid,omitempty"`
}

type ProtoRepository struct {
	Id         uuid.UUID                 `json:"id,omitempty"`
	Name       string                    `json:"name,omitempty"`
	Team       string                    `json:"team,omitempty"`
	IsDeleted  bool                      `json:"deleted,omitempty"`
	IsActive   bool                      `json:"active,omitempty"`
	Details    ProtoRepositoryDetails    `json:"details,omitempty"`
	Properties []ProtoRepositoryProperty `json:"properties,omitempty"`
}

type ProtoRepositoryFilter struct {
	Name      string `json:"name,omitempty"`
	Team      string `json:"team,omitempty"`
	IsDeleted bool   `json:"deleted,omitempty"`
	IsActive  bool   `json:"active,omitempty"`
}

type ProtoRepositoryDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

type ProtoRepositoryProperty struct {
	Type         string `json:"type,omitempty"`
	View         string `json:"view,omitempty"`
	Property     string `json:"property,omitempty"`
	Value        string `json:"value,omitempty"`
	Inheritance  bool   `json:"inheritance,omitempty"`
	ChildrenOnly bool   `json:"children,omitempty"`
	Source       string `json:"source,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
