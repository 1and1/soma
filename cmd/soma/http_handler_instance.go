/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
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

// InstanceShow returns information about a check instance
func InstanceShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	/* repository scope is currently not checked
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`instance_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	*/

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`instance_r`].(*instance)
	handler.input <- msg.Request{
		Type:       `instance`,
		Action:     `show`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    false,
		Instance: proto.Instance{
			Id: params.ByName(`instance`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// InstanceVersions returns information about a check instance's
// version history
func InstanceVersions(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	/* repository scope is currently not checked
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`instance_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	*/

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`instance_r`].(*instance)
	handler.input <- msg.Request{
		Type:       `instance`,
		Action:     `versions`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    false,
		Instance: proto.Instance{
			Id: params.ByName(`instance`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// InstanceList returns the list of instances in the subtree
// below the queried object.
// Currently only supports repositories and buckets as target.
func InstanceList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	/* repository scope is currently not checked
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`instance_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	*/

	listT := ``
	switch {
	case params.ByName(`repository`) != ``:
		listT = `repository`
	case params.ByName(`bucket`) != ``:
		listT = `bucket`
	case params.ByName(`group`) != ``:
		fallthrough
	case params.ByName(`cluster`) != ``:
		fallthrough
	case params.ByName(`node`) != ``:
		DispatchNotImplemented(&w, nil)
		return
	default:
		DispatchBadRequest(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`instance_r`].(*instance)
	handler.input <- msg.Request{
		Type:       `instance`,
		Action:     `list`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    false,
		Instance: proto.Instance{
			Id:         params.ByName(`instance`),
			ObjectId:   params.ByName(listT),
			ObjectType: listT,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// InstanceListAll is an administrative action that lists all
// check instances on the system
func InstanceListAll(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if ok, isAdmin := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`instance_list_all`, ``, ``, ``); !(ok && isAdmin) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`instance_r`].(*instance)
	handler.input <- msg.Request{
		Type:       `instance`,
		Action:     `list_all`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		IsAdmin:    false,
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
