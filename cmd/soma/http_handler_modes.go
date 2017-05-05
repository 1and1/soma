package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// ModeList function
func ModeList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `mode`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["modeReadHandler"].(*somaModeReadHandler)
	handler.input <- somaModeRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendModeReply(&w, &result)
}

// ModeShow function
func ModeShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `mode`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["modeReadHandler"].(*somaModeReadHandler)
	handler.input <- somaModeRequest{
		action: "show",
		reply:  returnChannel,
		Mode: proto.Mode{
			Mode: params.ByName("mode"),
		},
	}
	result := <-returnChannel
	SendModeReply(&w, &result)
}

// ModeAdd function
func ModeAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `mode`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewModeRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["modeWriteHandler"].(*somaModeWriteHandler)
	handler.input <- somaModeRequest{
		action: "add",
		reply:  returnChannel,
		Mode: proto.Mode{
			Mode: cReq.Mode.Mode,
		},
	}
	result := <-returnChannel
	SendModeReply(&w, &result)
}

// ModeRemove function
func ModeRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `mode`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["modeWriteHandler"].(*somaModeWriteHandler)
	handler.input <- somaModeRequest{
		action: "delete",
		reply:  returnChannel,
		Mode: proto.Mode{
			Mode: params.ByName("mode"),
		},
	}
	result := <-returnChannel
	SendModeReply(&w, &result)
}

// SendModeReply function
func SendModeReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewModeResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Modes {
		*result.Modes = append(*result.Modes, i.Mode)
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
