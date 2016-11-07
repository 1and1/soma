/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// WorkflowSummary returns information about the current workflow
// status distribution
func WorkflowSummary(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`workflow_summary`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`workflow_r`].(*workflowRead)
	handler.input <- msg.Request{
		Type:       `workflow`,
		Action:     `summary`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    false,
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func WorkflowList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`workflow_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewWorkflowFilter()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if cReq.Filter.Workflow.Status == `` {
		DispatchBadRequest(&w, fmt.Errorf(
			`No workflow status specified`))
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`workflow_r`].(*workflowRead)
	handler.input <- msg.Request{
		Type:       `workflow`,
		Action:     `list`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    false,
		Workflow: proto.Workflow{
			Status: cReq.Filter.Workflow.Status,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
