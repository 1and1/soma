package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// StatusList function
func StatusList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `status`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusReadHandler"].(*somaStatusReadHandler)
	handler.input <- somaStatusRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

// StatusShow function
func StatusShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `status`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusReadHandler"].(*somaStatusReadHandler)
	handler.input <- somaStatusRequest{
		action: "show",
		reply:  returnChannel,
		Status: proto.Status{
			Name: params.ByName("status"),
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

// StatusAdd function
func StatusAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `status`,
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

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusWriteHandler"].(*somaStatusWriteHandler)
	handler.input <- somaStatusRequest{
		action: "add",
		reply:  returnChannel,
		Status: proto.Status{
			Name: cReq.Status.Name,
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

// StatusRemove function
func StatusRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `status`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusWriteHandler"].(*somaStatusWriteHandler)
	handler.input <- somaStatusRequest{
		action: "delete",
		reply:  returnChannel,
		Status: proto.Status{
			Name: params.ByName("status"),
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

// SendStatusReply function
func SendStatusReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewStatusResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Status {
		*result.Status = append(*result.Status, i.Status)
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
		}
	}

dispatch:
	json, err := json.Marshal(result)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
