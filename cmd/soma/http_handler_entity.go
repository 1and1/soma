/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// EntityList function
func EntityList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `entity`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`entity_r`].(*entityRead)
	handler.input <- msg.Request{
		Section:    `entity`,
		Action:     `list`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// EntityShow function
func EntityShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `entity`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`entity_r`].(*entityRead)
	handler.input <- msg.Request{
		Section:    `entity`,
		Action:     `show`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Entity: proto.Entity{
			Name: params.ByName(`entity`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// EntityAdd function
func EntityAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `entity`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewEntityRequest()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`entity_w`].(*entityWrite)
	handler.input <- msg.Request{
		Section:    `entity`,
		Action:     `add`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Entity: proto.Entity{
			Name: cReq.Entity.Name,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// EntityRemove function
func EntityRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `entity`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`entity_w`].(*entityWrite)
	handler.input <- msg.Request{
		Section:    `entity`,
		Action:     `remove`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Entity: proto.Entity{
			Name: params.ByName(`entity`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// EntityRename function
func EntityRename(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `entity`,
		Action:     `rename`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewEntityRequest()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`entity_w`].(*entityWrite)
	handler.input <- msg.Request{
		Section:    `entity`,
		Action:     `remove`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Update: msg.UpdateData{
			Entity: proto.Entity{
				Name: cReq.Entity.Name,
			},
		},
		Entity: proto.Entity{
			Name: params.ByName(`entity`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
