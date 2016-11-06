/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Filter struct {
	Bucket      *BucketFilter      `json:"bucket,omitempty"`
	Capability  *CapabilityFilter  `json:"capability,omitempty"`
	CheckConfig *CheckConfigFilter `json:"checkConfig,omitempty"`
	Cluster     *ClusterFilter     `json:"cluster,omitempty"`
	Grant       *GrantFilter       `json:"grant,omitempty"`
	Group       *GroupFilter       `json:"group,omitempty"`
	Job         *JobFilter         `json:"job,omitempty"`
	Level       *LevelFilter       `json:"level,omitempty"`
	Metric      *MetricFilter      `json:"metric,omitempty"`
	Monitoring  *MonitoringFilter  `json:"monitoring,omitempty"`
	Node        *NodeFilter        `json:"node,omitempty"`
	Oncall      *OncallFilter      `json:"oncall,omitempty"`
	Permission  *PermissionFilter  `json:"permission,omitempty"`
	Property    *PropertyFilter    `json:"property,omitempty"`
	Provider    *ProviderFilter    `json:"provider,omitempty"`
	Repository  *RepositoryFilter  `json:"repository,omitempty"`
	Server      *ServerFilter      `json:"server,omitempty"`
	Team        *TeamFilter        `json:"team,omitempty"`
	Unit        *UnitFilter        `json:"unit,omitempty"`
	User        *UserFilter        `json:"user,omitempty"`
	Workflow    *WorkflowFilter    `json:"workflow,omitempty"`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
