package somaproto

import (
//"github.com/satori/go.uuid"
)

type ProtoRequestProperty struct {
	Custom  ProtoCustomProperty  `json:"custom,omitempty"`
	System  ProtoSystemProperty  `json:"system,omitempty"`
	Service ProtoServiceProperty `json:"service,omitempty"`
}

type ProtoResultProperty struct {
	Code    uint16                 `json:"code,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Text    []string               `json:"text,omitempty"`
	Custom  []ProtoCustomProperty  `json:"custom,omitempty"`
	System  []ProtoSystemProperty  `json:"system,omitempty"`
	Service []ProtoServiceProperty `json:"service,omitempty"`
}

type ProtoCustomProperty struct {
	Repository string `json:"repository,omitempty"`
	Property   string `json:"property,omitempty"`
	Value      string `json:"value,omitempty"`
}

type ProtoSystemProperty struct {
	Property string `json:"property,omitempty"`
	Value    string `json:"value,omitempty"`
}

type ProtoServiceProperty struct {
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
