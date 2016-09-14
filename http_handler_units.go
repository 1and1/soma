package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListUnit(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`units_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["unitReadHandler"].(*somaUnitReadHandler)
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
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`units_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["unitReadHandler"].(*somaUnitReadHandler)
	handler.input <- somaUnitRequest{
		action: "show",
		reply:  returnChannel,
		Unit: proto.Unit{
			Unit: params.ByName("unit"),
		},
	}
	result := <-returnChannel
	SendUnitReply(&w, &result)
}

/* Write functions
 */
func AddUnit(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`units_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewUnitRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["unitWriteHandler"].(*somaUnitWriteHandler)
	handler.input <- somaUnitRequest{
		action: "add",
		reply:  returnChannel,
		Unit: proto.Unit{
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
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`units_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["unitWriteHandler"].(*somaUnitWriteHandler)
	handler.input <- somaUnitRequest{
		action: "delete",
		reply:  returnChannel,
		Unit: proto.Unit{
			Unit: params.ByName("unit"),
		},
	}
	result := <-returnChannel
	SendUnitReply(&w, &result)
}

/* Utility
 */
func SendUnitReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewUnitResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Units {
		*result.Units = append(*result.Units, i.Unit)
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
