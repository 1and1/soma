package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// DatacenterList function
func DatacenterList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `datacenter`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterReadHandler"].(*somaDatacenterReadHandler)
	handler.input <- somaDatacenterRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

// DatacenterSync function
func DatacenterSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `datacenter`,
		Action:     `sync`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterReadHandler"].(*somaDatacenterReadHandler)
	handler.input <- somaDatacenterRequest{
		action: `sync`,
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

// DatacenterShow function
func DatacenterShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `datacenter`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterReadHandler"].(*somaDatacenterReadHandler)
	handler.input <- somaDatacenterRequest{
		action: "show",
		Datacenter: proto.Datacenter{
			Locode: params.ByName("datacenter"),
		},
		reply: returnChannel,
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

// DatacenterAdd function
func DatacenterAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `datacenter`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterWriteHandler"].(*somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action: "add",
		reply:  returnChannel,
		Datacenter: proto.Datacenter{
			Locode: cReq.Datacenter.Locode,
		},
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

// DatacenterRemove function
func DatacenterRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `datacenter`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterWriteHandler"].(*somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action: "delete",
		reply:  returnChannel,
		Datacenter: proto.Datacenter{
			Locode: params.ByName("datacenter"),
		},
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

// DatacenterRename function
func DatacenterRename(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `datacenter`,
		Action:     `rename`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["datacenterWriteHandler"].(*somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action: "rename",
		Datacenter: proto.Datacenter{
			Locode: params.ByName("datacenter"),
		},
		rename: cReq.Datacenter.Locode,
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendDatacenterReply(&w, &result)
}

// SendDatacenterReply function
func SendDatacenterReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewDatacenterResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Datacenters {
		*result.Datacenters = append(*result.Datacenters, i.Datacenter)
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
