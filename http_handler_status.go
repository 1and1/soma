package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListStatus(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusReadHandler"].(somaStatusReadHandler)
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

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusReadHandler"].(somaStatusReadHandler)
	handler.input <- somaStatusRequest{
		action: "show",
		reply:  returnChannel,
		Status: somaproto.ProtoStatus{
			Status: params.ByName("status"),
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

/* Write functions
 */
func AddStatus(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestStatus{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusWriteHandler"].(somaStatusWriteHandler)
	handler.input <- somaStatusRequest{
		action: "add",
		reply:  returnChannel,
		Status: somaproto.ProtoStatus{
			Status: cReq.Status.Status,
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

func DeleteStatus(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["statusWriteHandler"].(somaStatusWriteHandler)
	handler.input <- somaStatusRequest{
		action: "delete",
		reply:  returnChannel,
		Status: somaproto.ProtoStatus{
			Status: params.ByName("status"),
		},
	}
	result := <-returnChannel
	SendStatusReply(&w, &result)
}

/* Utility
 */
func SendStatusReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultStatus{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.StatusList = make([]somaproto.ProtoStatus, 0)
	for _, i := range (*r).Status {
		result.StatusList = append(result.StatusList, i.Status)
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
