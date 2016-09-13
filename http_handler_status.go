package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListStatus(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`status_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusReadHandler"].(*somaStatusReadHandler)
	handler.input <- somaStatusRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

func ShowStatus(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`status_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusReadHandler"].(*somaStatusReadHandler)
	handler.input <- somaStatusRequest{
		action: "show",
		reply:  returnChannel,
		Status: proto.Status{
			Name: params.ByName("status"),
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

/* Write functions
 */
func AddStatus(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`status_create`, ``, ``, ``); !ok {
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
	handler := handlerMap["statusWriteHandler"].(*somaStatusWriteHandler)
	handler.input <- somaStatusRequest{
		action: "add",
		reply:  returnChannel,
		Status: proto.Status{
			Name: cReq.Status.Name,
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

func DeleteStatus(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`status_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusWriteHandler"].(*somaStatusWriteHandler)
	handler.input <- somaStatusRequest{
		action: "delete",
		reply:  returnChannel,
		Status: proto.Status{
			Name: params.ByName("status"),
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

/* Utility
 */
func SendStatusReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewStatusResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Status {
		*result.Status = append(*result.Status, i.Status)
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
