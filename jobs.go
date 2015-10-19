package somaproto

import (
	"github.com/satori/go.uuid"
)

type ProtoRequestJob struct {
	JobId      uuid.UUID          `json:"jobid,omitempty"`
	JobType    string             `json:"jobtype"`
	Server     ProtoJobServer     `json:"server,omitempty"`
	Team       ProtoJobTeam       `json:"team,omitempty"`
	Bucket     ProtoJobBucket     `json:"bucket,omitempty"`
	Repository ProtoJobRepository `json:"repository,omitempty"`
}

type ProtoResultJob struct {
	Code   uint16    `json:"code,omitempty"`
	Status string    `json:"status,omitempty"`
	Text   []string  `json:"text,omitempty"`
	JobId  uuid.UUID `json:"jobid"`
}

type ProtoJobServer struct {
	Action string      `json:"action"`
	Server ProtoServer `json:"server"`
}

type ProtoJobTeam struct {
}

type ProtoJobBucket struct {
}

type ProtoJobRepository struct {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
