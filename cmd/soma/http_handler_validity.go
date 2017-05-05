package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// ValidityList function
func ValidityList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `validity`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["validityReadHandler"].(*somaValidityReadHandler)
	handler.input <- somaValidityRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendValidityReply(&w, &result)
}

// ValidityShow function
func ValidityShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `validity`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["validityReadHandler"].(*somaValidityReadHandler)
	handler.input <- somaValidityRequest{
		action: "show",
		reply:  returnChannel,
		Validity: proto.Validity{
			SystemProperty: params.ByName("property"),
		},
	}
	result := <-returnChannel
	SendValidityReply(&w, &result)
}

// ValidityAdd function
func ValidityAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `validity`,
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
	handler := handlerMap["validityWriteHandler"].(*somaValidityWriteHandler)
	handler.input <- somaValidityRequest{
		action:   "add",
		reply:    returnChannel,
		Validity: *cReq.Validity,
	}
	result := <-returnChannel
	SendValidityReply(&w, &result)
}

// ValidityRemove function
func ValidityRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `validity`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["validityWriteHandler"].(*somaValidityWriteHandler)
	handler.input <- somaValidityRequest{
		action: "delete",
		reply:  returnChannel,
		Validity: proto.Validity{
			SystemProperty: params.ByName("property"),
		},
	}
	result := <-returnChannel
	SendValidityReply(&w, &result)
}

// SendValidityReply function
func SendValidityReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewValidityResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Validity {
		*result.Validities = append(*result.Validities, i.Validity)
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
