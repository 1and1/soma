package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListPredicate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaPredicateResult)
	handler := handlerMap["predicateReadHandler"].(somaPredicateReadHandler)
	handler.input <- somaPredicateRequest{
		action: "list",
		reply:  returnChannel,
	}
	results := <-returnChannel
	SendPredicateReply(&w, &results)
}

func ShowPredicate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaPredicateResult)
	handler := handlerMap["predicateReadHandler"].(somaPredicateReadHandler)
	handler.input <- somaPredicateRequest{
		action: "show",
		reply:  returnChannel,
		predicate: somaproto.ProtoPredicate{
			Predicate: params.ByName("predicate"),
		},
	}
	results := <-returnChannel
	SendPredicateReply(&w, &results)
}

func AddPredicate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	var cReq somaproto.ProtoRequestPredicate
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan []somaPredicateResult)
	handler := handlerMap["predicateWriteHandler"].(somaPredicateWriteHandler)
	handler.input <- somaPredicateRequest{
		action: "add",
		reply:  returnChannel,
		predicate: somaproto.ProtoPredicate{
			Predicate: cReq.Predicate.Predicate,
		},
	}
	results := <-returnChannel
	SendPredicateReply(&w, &results)
}

func DeletePredicate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaPredicateResult)
	handler := handlerMap["predicateWriteHandler"].(somaPredicateWriteHandler)
	handler.input <- somaPredicateRequest{
		action: "delete",
		reply:  returnChannel,
		predicate: somaproto.ProtoPredicate{
			Predicate: params.ByName("predicate"),
		},
	}
	results := <-returnChannel
	SendPredicateReply(&w, &results)
}

func SendPredicateReply(w *http.ResponseWriter, r *[]somaPredicateResult) {
	var res somaproto.ProtoResultPredicate
	dispatchError := CheckErrorHandler(r, &res)
	if dispatchError {
		goto dispatch
	}
	res.Text = make([]string, 0)
	res.Predicates = make([]somaproto.ProtoPredicate, 0)
	for _, l := range *r {
		res.Predicates = append(res.Predicates, l.predicate)
		if l.lErr != nil {
			res.Text = append(res.Text, l.lErr.Error())
		}
	}

dispatch:
	json, err := json.Marshal(res)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
