package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*Read functions
 */
func ListPredicate(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["predicateReadHandler"].(somaPredicateReadHandler)
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

	returnChannel := make(chan somaResult)
	handler := handlerMap["predicateReadHandler"].(somaPredicateReadHandler)
	handler.input <- somaPredicateRequest{
		action: "show",
		reply:  returnChannel,
		Predicate: somaproto.ProtoPredicate{
			Predicate: params.ByName("predicate"),
		},
	}
	result := <-returnChannel
	SendPredicateReply(&w, &result)
}

/* Write functions
 */
func AddPredicate(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestPredicate{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["predicateWriteHandler"].(somaPredicateWriteHandler)
	handler.input <- somaPredicateRequest{
		action: "add",
		reply:  returnChannel,
		Predicate: somaproto.ProtoPredicate{
			Predicate: cReq.Predicate.Predicate,
		},
	}
	result := <-returnChannel
	SendPredicateReply(&w, &result)
}

func DeletePredicate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["predicateWriteHandler"].(somaPredicateWriteHandler)
	handler.input <- somaPredicateRequest{
		action: "delete",
		reply:  returnChannel,
		Predicate: somaproto.ProtoPredicate{
			Predicate: params.ByName("predicate"),
		},
	}
	result := <-returnChannel
	SendPredicateReply(&w, &result)
}

/* Utility
 */
func SendPredicateReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultPredicate{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Predicates = make([]somaproto.ProtoPredicate, 0)
	for _, i := range (*r).Predicates {
		result.Predicates = append(result.Predicates, i.Predicate)
		if i.ResultError != nil {
			result.Text = append(result.Text, i.ResultError.Error())
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
