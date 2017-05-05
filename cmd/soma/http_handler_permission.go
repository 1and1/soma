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
	"strings"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// PermissionList function
func PermissionList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `permission`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `permission`,
		Action:     `list`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Permission: proto.Permission{
			Category: params.ByName(`category`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// PermissionShow function
func PermissionShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `permission`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `permission`,
		Action:     `show`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Permission: proto.Permission{
			Id:       params.ByName(`permission`),
			Category: params.ByName(`category`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// PermissionSearch function
func PermissionSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `permission`,
		Action:     `search`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewPermissionFilter()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	mr := msg.Request{
		Section:    `permission`,
		Action:     `search/name`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Permission: proto.Permission{
			Name:     cReq.Filter.Permission.Name,
			Category: cReq.Filter.Permission.Category,
		},
	}

	handler.input <- mr
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// PermissionAdd function
func PermissionAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `permission`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewPermissionRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if cReq.Permission.Category != params.ByName(`category`) {
		DispatchBadRequest(&w, fmt.Errorf(`Category mismatch`))
		return
	}
	if strings.Contains(params.ByName(`category`), `:grant`) {
		DispatchBadRequest(&w, fmt.Errorf(
			`Permissions in :grant categories are auto-managed.`))
		return
	}
	if params.ByName(`category`) == `system` ||
		params.ByName(`category`) == `omnipotence` {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `permission`,
		Action:     `add`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Permission: proto.Permission{
			Name:     cReq.Permission.Name,
			Category: cReq.Permission.Category,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// PermissionRemove function
func PermissionRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `permission`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}
	if strings.Contains(params.ByName(`category`), `:grant`) {
		DispatchBadRequest(&w, fmt.Errorf(
			`Permissions in :grant categories are auto-managed.`))
		return
	}
	if params.ByName(`category`) == `system` ||
		params.ByName(`category`) == `omnipotence` {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `permission`,
		Action:     `remove`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Permission: proto.Permission{
			Id:       params.ByName(`permission`),
			Category: params.ByName(`category`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// PermissionEdit function
func PermissionEdit(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `permission`,
		Action:     `edit`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewPermissionRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if cReq.Permission.Category != params.ByName(`category`) {
		DispatchBadRequest(&w, fmt.Errorf(`Category mismatch`))
		return
	}
	if cReq.Permission.Id != params.ByName(`permission`) {
		DispatchBadRequest(&w, fmt.Errorf(`PermissionId mismatch`))
		return
	}
	if strings.Contains(params.ByName(`category`), `:grant`) {
		DispatchBadRequest(&w, fmt.Errorf(
			`Permissions in :grant categories can not be mapped`))
		return
	}
	if params.ByName(`category`) == `system` ||
		params.ByName(`category`) == `omnipotence` {
		DispatchForbidden(&w, nil)
		return
	}
	// invalid: map+unmap at the same time
	if cReq.Flags.Add && cReq.Flags.Remove {
		DispatchBadRequest(&w, fmt.Errorf(`Ambiguous instruction`))
		return
	}
	// invalid: batched mapping
	if cReq.Permission.Actions != nil && cReq.Permission.Sections != nil {
		DispatchBadRequest(&w, fmt.Errorf(`Invalid batch mapping`))
		return
	}
	if cReq.Permission.Actions != nil {
		if len(*cReq.Permission.Actions) != 1 ||
			params.ByName(`category`) != (*cReq.Permission.Actions)[0].Category {
			DispatchBadRequest(&w, fmt.Errorf(`Invalid action specification`))
			return
		}
	}
	if cReq.Permission.Sections != nil {
		if len(*cReq.Permission.Sections) != 1 ||
			params.ByName(`category`) != (*cReq.Permission.Sections)[0].Category {
			DispatchBadRequest(&w, fmt.Errorf(`Invalid section specification`))
			return
		}
	}

	var task string
	if cReq.Flags.Add {
		task = `map`
	}
	if cReq.Flags.Remove {
		task = `unmap`
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `permission`,
		Action:     task,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Permission: proto.Permission{
			Id:       cReq.Permission.Id,
			Name:     cReq.Permission.Name,
			Category: cReq.Permission.Category,
			Sections: cReq.Permission.Sections,
			Actions:  cReq.Permission.Actions,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
