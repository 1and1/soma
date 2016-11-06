/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package proto

type Workflow struct {
	InstanceId       string           `json:"instanceId,omitempty"`
	InstanceConfigId string           `json:"instanceConfigId,omitempty"`
	Status           string           `json:"status,omitempty"`
	NextStatus       string           `json:"nextStatus,omitempty"`
	Summary          *WorkflowSummary `json:"summary,omitempty"`
	Instances        *[]Instance      `json:"instances,omitempty"`
}

type WorkflowSummary struct {
	AwaitingComputation   uint64 `json:"awaitingComputation"`
	Computed              uint64 `json:"computed"`
	AwaitingRollout       uint64 `json:"awaitingRollout"`
	RolloutInProgress     uint64 `json:"rolloutInProgress"`
	RolloutFailed         uint64 `json:"rolloutFailed"`
	Active                uint64 `json:"active"`
	AwaitingDeprovision   uint64 `json:"awaitingDeprovision"`
	DeprovisionInProgress uint64 `json:"deprovisionInProgress"`
	DeprovisionFailed     uint64 `json:"deprovisionFailed"`
	Deprovisioned         uint64 `json:"deprovisioned"`
	AwaitingDeletion      uint64 `json:"awaitingDeletion"`
	Blocked               uint64 `json:"blocked"`
}

type WorkflowFilter struct {
	Status string `json:"status"`
}

func NewWorkflowRequest() Request {
	return Request{
		Flags:    &Flags{},
		Workflow: &Workflow{},
	}
}

func NewWorkflowFilter() Request {
	return Request{
		Filter: &Filter{
			Workflow: &WorkflowFilter{},
		},
	}
}

func NewWorkflowResult() Result {
	return Result{
		Errors:    &[]string{},
		Workflows: &[]Workflow{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
