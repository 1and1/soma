package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListUnit(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["unitReadHandler"].(somaUnitReadHandler)
	handler.input <- somaUnitRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendUnitReply(&w, &result)
}

func ShowUnit(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["unitReadHandler"].(somaUnitReadHandler)
	handler.input <- somaUnitRequest{
		action: "show",
		reply:  returnChannel,
		Unit: somaproto.ProtoUnit{
			Unit: params.ByName("unit"),
		},
	}
	result := <-returnChannel
	SendUnitReply(&w, &result)
}

/* Write functions
 */
func AddUnit(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestUnit{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["unitWriteHandler"].(somaUnitWriteHandler)
	handler.input <- somaUnitRequest{
		action: "add",
		reply:  returnChannel,
		Unit: somaproto.ProtoUnit{
			Unit: cReq.Unit.Unit,
			Name: cReq.Unit.Name,
		},
	}
	result := <-returnChannel
	SendUnitReply(&w, &result)
}

func DeleteUnit(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["unitWriteHandler"].(somaUnitWriteHandler)
	handler.input <- somaUnitRequest{
		action: "delete",
		reply:  returnChannel,
		Unit: somaproto.ProtoUnit{
			Unit: params.ByName("unit"),
		},
	}
	result := <-returnChannel
	SendUnitReply(&w, &result)
}

/* Utility
 */
func SendUnitReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultUnit{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Units = make([]somaproto.ProtoUnit, 0)
	for _, i := range (*r).Units {
		result.Units = append(result.Units, i.Unit)
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
