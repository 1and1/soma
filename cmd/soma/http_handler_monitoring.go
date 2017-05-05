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

// MonitoringList function
func MonitoringList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	// check for operations runtime privileges
	admin := IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `runtime`,
		Action:     `monitoringsystem_list_all`,
	})

	// skip the regular permission check, if the user has
	// the operations permission
	if admin {
		goto authorized
	}

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `monitoringsystem`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

authorized:
	returnChannel := make(chan msg.Result)
	handler := handlerMap[`monitoring_r`].(*monitoringRead)
	handler.input <- msg.Request{
		Section:    `monitoringsystem`,
		Action:     `list`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Flag: msg.Flags{
			Unscoped: admin,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// MonitoringSearch function
func MonitoringSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	// check for operations runtime privileges
	admin := IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `runtime`,
		Action:     `monitoringsystem_list_all`,
	})

	// skip the regular permission check, if the user has
	// the operations permission
	if admin {
		goto authorized
	}

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `monitoringsystem`,
		Action:     `search`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

authorized:
	cReq := proto.NewMonitoringFilter()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if cReq.Filter.Monitoring.Name == `` {
		DispatchBadRequest(&w, fmt.Errorf(
			`Empty search request: name missing`))
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`monitoring_r`].(*monitoringRead)
	handler.input <- msg.Request{
		Section:    `monitoringsystem`,
		Action:     `search`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Flag: msg.Flags{
			Unscoped: admin,
		},
		Monitoring: proto.Monitoring{
			Name: cReq.Filter.Monitoring.Name,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// MonitoringShow function
func MonitoringShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:     params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `monitoringsystem`,
		Action:       `show`,
		MonitoringID: params.ByName(`monitoring`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`monitoring_r`].(*monitoringRead)
	handler.input <- msg.Request{
		Section:    `monitoringsystem`,
		Action:     `show`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Monitoring: proto.Monitoring{
			Id: params.ByName(`monitoring`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// MonitoringAdd function
func MonitoringAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `monitoringsystem`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewMonitoringRequest()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.Monitoring.Name, `.`) {
		DispatchBadRequest(&w, fmt.Errorf(
			`Invalid monitoring system`+
				` name containing . character`))
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`monitoring_w`].(*monitoringWrite)
	handler.input <- msg.Request{
		Section:    `monitoringsystem`,
		Action:     `add`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Monitoring: proto.Monitoring{
			Name:     cReq.Monitoring.Name,
			Mode:     cReq.Monitoring.Mode,
			Contact:  cReq.Monitoring.Contact,
			TeamId:   cReq.Monitoring.TeamId,
			Callback: cReq.Monitoring.Callback,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// MonitoringRemove function
func MonitoringRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `monitoringsystem`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`monitoring_w`].(*monitoringWrite)
	handler.input <- msg.Request{
		Section:    `monitoringsystem`,
		Action:     `remove`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Monitoring: proto.Monitoring{
			Id: params.ByName(`monitoring`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
