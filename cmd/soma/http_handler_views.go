package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// ViewList function
func ViewList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `view`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewReadHandler"].(*somaViewReadHandler)
	handler.input <- somaViewRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

// ViewShow function
func ViewShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `view`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewReadHandler"].(*somaViewReadHandler)
	handler.input <- somaViewRequest{
		action: "show",
		reply:  returnChannel,
		View: proto.View{
			Name: params.ByName("view"),
		},
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

// ViewAdd function
func ViewAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `view`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewViewRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.View.Name, `.`) {
		DispatchBadRequest(&w, fmt.Errorf(`Invalid view name containing . character`))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewWriteHandler"].(*somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "add",
		reply:  returnChannel,
		View: proto.View{
			Name: cReq.View.Name,
		},
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

// ViewRemove function
func ViewRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `view`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewWriteHandler"].(*somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "delete",
		reply:  returnChannel,
		View: proto.View{
			Name: params.ByName("view"),
		},
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

// ViewRename function
func ViewRename(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `view`,
		Action:     `rename`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewViewRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewWriteHandler"].(*somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "rename",
		reply:  returnChannel,
		name:   params.ByName("view"),
		View: proto.View{
			Name: cReq.View.Name,
		},
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

// SendViewReply function
func SendViewReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewViewResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Views {
		*result.Views = append(*result.Views, i.View)
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
