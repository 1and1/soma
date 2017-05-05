package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// PredicateList function
func PredicateList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `predicate`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["predicateReadHandler"].(*somaPredicateReadHandler)
	handler.input <- somaPredicateRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendPredicateReply(&w, &result)
}

// PredicateShow function
func PredicateShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `predicate`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["predicateReadHandler"].(*somaPredicateReadHandler)
	handler.input <- somaPredicateRequest{
		action: "show",
		reply:  returnChannel,
		Predicate: proto.Predicate{
			Symbol: params.ByName("predicate"),
		},
	}
	result := <-returnChannel
	SendPredicateReply(&w, &result)
}

// PredicateAdd function
func PredicateAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `predicate`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewPredicateRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["predicateWriteHandler"].(*somaPredicateWriteHandler)
	handler.input <- somaPredicateRequest{
		action: "add",
		reply:  returnChannel,
		Predicate: proto.Predicate{
			Symbol: cReq.Predicate.Symbol,
		},
	}
	result := <-returnChannel
	SendPredicateReply(&w, &result)
}

// PredicateRemove function
func PredicateRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `predicate`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["predicateWriteHandler"].(*somaPredicateWriteHandler)
	handler.input <- somaPredicateRequest{
		action: "delete",
		reply:  returnChannel,
		Predicate: proto.Predicate{
			Symbol: params.ByName("predicate"),
		},
	}
	result := <-returnChannel
	SendPredicateReply(&w, &result)
}

// SendPredicateReply function
func SendPredicateReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewPredicateResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Predicates {
		*result.Predicates = append(*result.Predicates, i.Predicate)
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
