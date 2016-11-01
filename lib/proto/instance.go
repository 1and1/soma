/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Instance struct {
	Id               string               `json:"id,omitempty"`
	Version          uint64               `json:"version"`
	CheckId          string               `json:"checkId,omitempty"`
	ConfigId         string               `json:"configId,omitempty"`
	InstanceConfigId string               `json:"instanceConfigId,omitempty"`
	RepositoryId     string               `json:"repositoryId,omitempty"`
	BucketId         string               `json:"bucketId,omitempty"`
	ObjectId         string               `json:"objectId,omitempty"`
	ObjectType       string               `json:"objectType,omitempty"`
	CurrentStatus    string               `json:"currentStatus,omitempty"`
	NextStatus       string               `json:"nextStatus,omitempty"`
	IsInherited      bool                 `json:"isInherited"`
	Info             *InstanceVersionInfo `json:"instanceVersionInfo,omitempty"`
}

type InstanceVersionInfo struct {
	CreatedAt           string `json:"createdAt"`
	ActivatedAt         string `json:"activatedAt"`
	DeprovisionedAt     string `json:"deprovisionedAt"`
	StatusLastUpdatedAt string `json:"statusLastUpdatedAt"`
	NotifiedAt          string `json:"notifiedAt"`
}

func NewInstanceResult() Result {
	return Result{
		Errors:    &[]string{},
		Instances: &[]Instance{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
