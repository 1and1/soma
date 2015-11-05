package somaproto

import (
	"github.com/satori/go.uuid"
)

type ProtoRequestBucket struct {
	Bucket  ProtoBucket       `json:"bucket,omitempty"`
	Filter  ProtoBucketFilter `json:"filter,omitempty"`
	Restore bool              `json:"restore,omitempty"`
	Purge   bool              `json:"purge,omitempty"`
	Freeze  bool              `json:"freeze,omitempty"`
	Thaw    bool              `json:"thaw,omitempty"`
}

type ProtoResultBucket struct {
	Code    uint16        `json:"code,omitempty"`
	Status  string        `json:"status,omitempty"`
	Text    []string      `json:"text,omitempty"`
	Buckets []ProtoBucket `json:"buckets,omitempty"`
	JobId   uuid.UUID     `json:"jobid,omitempty"`
}

type ProtoBucket struct {
	Id           uuid.UUID          `json:"id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Repository   string             `json:"repository,omitempty"`
	RepositoryId string             `json:"repositoryid,omitempty"`
	Team         string             `json:"team,omitempty"`
	TeamId       string             `json:"team,omitempty"`
	Environment  string             `json:"environment,omitempty"`
	IsDeleted    bool               `json:"deleted,omitempty"`
	IsFrozen     bool               `json:"frozen,omitempty"`
	Details      ProtoBucketDetails `json:"details,omitempty"`
	//	Properties []ProtoBucketProperty `json:"properties,omitempty"`
}

type ProtoBucketFilter struct {
	Name         string `json:"name,omitempty"`
	RepositoryId string `json:"repositoryid,omitempty"`
	IsDeleted    bool   `json:"deleted,omitempty"`
	IsFrozen     bool   `json:"frozen,omitempty"`
}

type ProtoBucketDetails struct {
	CreatedAt string `json:"createdat,omitempty"`
	CreatedBy string `json:"createdby,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
