package somaproto

import (
	"github.com/satori/go.uuid"
)

type ProtoRequestProperty struct {
	Custom  ProtoPropertyCustom  `json:"custom,omitempty"`
	System  ProtoPropertySystem  `json:"system,omitempty"`
	Service ProtoPropertyService `json:"service,omitempty"`
	Filter  ProtoPropertyFilter  `json:"filter,omitempty"`
}

type ProtoResultProperty struct {
	Code    uint16                 `json:"code,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Text    []string               `json:"text,omitempty"`
	Custom  []ProtoPropertyCustom  `json:"custom,omitempty"`
	System  []ProtoPropertySystem  `json:"system,omitempty"`
	Service []ProtoPropertyService `json:"service,omitempty"`
}

type ProtoPropertyCustom struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Repository string    `json:"repository,omitempty"`
	Property   string    `json:"property,omitempty"`
	Value      string    `json:"value,omitempty"`
}

type ProtoPropertySystem struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Property string    `json:"property,omitempty"`
	Value    string    `json:"value,omitempty"`
}

type ProtoPropertyService struct {
	Id         uuid.UUID              `json:"id,omitempty"`
	Property   string                 `json:"property,omitempty"`
	Team       string                 `json:"team,omitempty"`
	Attributes ProtoServiceAttributes `json:"attributes,omitempty"`
}

type ProtoServiceAttributes struct {
	Transport   []string `json:"proto_transport,omitempty"`
	Application []string `json:"proto_application,omitempty"`
	Port        []string `json:"port,omitempty"`
	Process     []string `json:"process,omitempty"`
	File        []string `json:"file,omitempty"`
	Directory   []string `json:"directory,omitempty"`
	Socket      []string `json:"socket,omitempty"`
	Uid         []string `json:"uid,omitempty"`
	Tls         string   `json:"tls,omitempty"`
	Provider    string   `json:"provider,omitempty"`
}

type ProtoPropertyFilter struct {
	Name string `json:"name,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
