package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// CapabilityList function
func CapabilityList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `capability`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["capabilityReadHandler"].(*somaCapabilityReadHandler)
	handler.input <- somaCapabilityRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.Request{}
	cReq.Filter = &proto.Filter{}
	cReq.Filter.Capability = &proto.CapabilityFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Capability.MonitoringId != "" {
		filtered := []somaCapabilityResult{}
		for _, i := range result.Capabilities {
			if i.Capability.MonitoringId == cReq.Filter.Capability.MonitoringId &&
				i.Capability.Metric == cReq.Filter.Capability.Metric &&
				i.Capability.View == cReq.Filter.Capability.View {
				filtered = append(filtered, i)
			}
		}
		result.Capabilities = filtered
	}

	// cleanup reply
	for i := range result.Capabilities {
		result.Capabilities[i].Capability.MonitoringId = ""
		result.Capabilities[i].Capability.Metric = ""
		result.Capabilities[i].Capability.View = ""
	}

skip:
	SendCapabilityReply(&w, &result)
}

// CapabilityShow function
func CapabilityShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:     params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `capability`,
		Action:       `show`,
		CapabilityID: params.ByName(`capability`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["capabilityReadHandler"].(*somaCapabilityReadHandler)
	handler.input <- somaCapabilityRequest{
		action: "show",
		reply:  returnChannel,
		Capability: proto.Capability{
			Id: params.ByName("capability"),
		},
	}
	result := <-returnChannel
	SendCapabilityReply(&w, &result)
}

// CapabilityAdd function
func CapabilityAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if !IsAuthorized(&msg.Authorization{
		AuthUser:     params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `capability`,
		Action:       `add`,
		MonitoringID: cReq.Capability.MonitoringId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["capabilityWriteHandler"].(*somaCapabilityWriteHandler)
	handler.input <- somaCapabilityRequest{
		action: "add",
		reply:  returnChannel,
		Capability: proto.Capability{
			MonitoringId: cReq.Capability.MonitoringId,
			Metric:       cReq.Capability.Metric,
			View:         cReq.Capability.View,
			Thresholds:   cReq.Capability.Thresholds,
		},
	}
	result := <-returnChannel
	SendCapabilityReply(&w, &result)
}

// CapabilityRemove function
func CapabilityRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:     params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `capability`,
		Action:       `remove`,
		CapabilityID: params.ByName(`capability`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["capabilityWriteHandler"].(*somaCapabilityWriteHandler)
	handler.input <- somaCapabilityRequest{
		action: "delete",
		reply:  returnChannel,
		Capability: proto.Capability{
			Id: params.ByName("capability"),
		},
	}
	result := <-returnChannel
	SendCapabilityReply(&w, &result)
}

// SendCapabilityReply function
func SendCapabilityReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewCapabilityResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Capabilities {
		*result.Capabilities = append(*result.Capabilities, i.Capability)
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
		}
	}

dispatch:
	result.Clean()
	json, err := json.Marshal(result)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
