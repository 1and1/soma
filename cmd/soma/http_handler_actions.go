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

// ActionList accepts requests to list actions in a specific section
func ActionList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `action`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `action`,
		Action:     `list`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		ActionObj: proto.Action{
			SectionId: params.ByName(`section`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// ActionShow accepts requests to show details about a specific
// action
func ActionShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `action`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `action`,
		Action:     `show`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		ActionObj: proto.Action{
			Id:        params.ByName(`action`),
			SectionId: params.ByName(`section`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func ActionSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `action`,
		Action:     `search`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if cReq.Action.SectionId == `` || cReq.Action.Name == `` {
		DispatchBadRequest(&w,
			fmt.Errorf(`Invalid action search specification`))
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `action`,
		Action:     `search`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		ActionObj: proto.Action{
			Name:      cReq.Action.Name,
			SectionId: cReq.Action.SectionId,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// ActionAdd accepts requests to add a new action to a section
func ActionAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `action`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if cReq.Action.SectionId != params.ByName(`section`) {
		DispatchBadRequest(&w, fmt.Errorf("SectionId mismatch: %s, %s",
			cReq.Action.SectionId, params.ByName(`section`)))
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `action`,
		Action:     `add`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		ActionObj: proto.Action{
			Name:      cReq.Action.Name,
			SectionId: cReq.Action.SectionId,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// ActionRemove accepts requests to remove an action form a section
func ActionRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `action`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `action`,
		Action:     `remove`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		ActionObj: proto.Action{
			Id:        params.ByName(`action`),
			SectionId: params.ByName(`section`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
