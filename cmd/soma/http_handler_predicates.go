package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

/*Read functions
 */
func ListPredicate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`predicates_list`, ``, ``, ``); !ok {
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

func ShowPredicate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`predicates_show`, ``, ``, ``); !ok {
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

/* Write functions
 */
func AddPredicate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`predicates_create`, ``, ``, ``); !ok {
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

func DeletePredicate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`predicates_delete`, ``, ``, ``); !ok {
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

/* Utility
 */
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
