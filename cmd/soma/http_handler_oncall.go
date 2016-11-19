package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// OncallList function
func OncallList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `oncall`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["oncallReadHandler"].(*somaOncallReadHandler)
	handler.input <- somaOncallRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.NewOncallFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Oncall.Name != "" {
		filtered := []somaOncallResult{}
		for _, i := range result.Oncall {
			if i.Oncall.Name == cReq.Filter.Oncall.Name {
				filtered = append(filtered, i)
			}
		}
		result.Oncall = filtered
	}

skip:
	SendOncallReply(&w, &result)
}

// OncallShow function
func OncallShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `oncall`,
		Action:     `show`,
		OncallID:   params.ByName(`oncall`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["oncallReadHandler"].(*somaOncallReadHandler)
	handler.input <- somaOncallRequest{
		action: "show",
		reply:  returnChannel,
		Oncall: proto.Oncall{
			Id: params.ByName("oncall"),
		},
	}
	result := <-returnChannel
	SendOncallReply(&w, &result)
}

// OncallAdd function
func OncallAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `oncall`,
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
	handler := handlerMap["oncallWriteHandler"].(*somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "add",
		reply:  returnChannel,
		Oncall: proto.Oncall{
			Name:   cReq.Oncall.Name,
			Number: cReq.Oncall.Number,
		},
	}
	result := <-returnChannel
	SendOncallReply(&w, &result)
}

// OncallUpdate function
func OncallUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `oncall`,
		Action:     `update`,
		OncallID:   params.ByName(`oncall`),
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
	handler := handlerMap["oncallWriteHandler"].(*somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "update",
		reply:  returnChannel,
		Oncall: proto.Oncall{
			Id:     params.ByName("oncall"),
			Name:   cReq.Oncall.Name,
			Number: cReq.Oncall.Number,
		},
	}
	result := <-returnChannel
	SendOncallReply(&w, &result)
}

// OncallRemove function
func OncallRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `oncall`,
		Action:     `remove`,
		OncallID:   params.ByName(`oncall`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["oncallWriteHandler"].(*somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "delete",
		reply:  returnChannel,
		Oncall: proto.Oncall{
			Id: params.ByName("oncall"),
		},
	}
	result := <-returnChannel
	SendOncallReply(&w, &result)
}

// SendOncallReply function
func SendOncallReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewOncallResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Oncall {
		*result.Oncalls = append(*result.Oncalls, i.Oncall)
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
