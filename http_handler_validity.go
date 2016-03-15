package main

import (
	"encoding/json"
	"net/http"


	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListValidity(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["validityReadHandler"].(somaValidityReadHandler)
	handler.input <- somaValidityRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendValidityReply(&w, &result)
}

func ShowValidity(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["validityReadHandler"].(somaValidityReadHandler)
	handler.input <- somaValidityRequest{
		action: "show",
		reply:  returnChannel,
		Validity: somaproto.Validity{
			SystemProperty: params.ByName("property"),
		},
	}
	result := <-returnChannel
	SendValidityReply(&w, &result)
}

/* Write functions
 */
func AddValidity(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ValidityRequest{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["validityWriteHandler"].(somaValidityWriteHandler)
	handler.input <- somaValidityRequest{
		action:   "add",
		reply:    returnChannel,
		Validity: *cReq.Validity,
	}
	result := <-returnChannel
	SendValidityReply(&w, &result)
}

func DeleteValidity(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["validityWriteHandler"].(somaValidityWriteHandler)
	handler.input <- somaValidityRequest{
		action: "delete",
		reply:  returnChannel,
		Validity: somaproto.Validity{
			SystemProperty: params.ByName("property"),
		},
	}
	result := <-returnChannel
	SendValidityReply(&w, &result)
}

/* Utility
 */
func SendValidityReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ValidityResult{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Validity = make([]somaproto.Validity, 0)
	for _, i := range (*r).Validity {
		result.Validity = append(result.Validity, i.Validity)
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
