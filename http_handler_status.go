package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaStatusResult)
	handler := handlerMap["statusReadHandler"].(somaStatusReadHandler)
	handler.input <- somaStatusRequest{
		action: "list",
		reply:  returnChannel,
	}
	results := <-returnChannel
	SendStatusReply(&w, &results)
}

func ShowStatus(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaStatusResult)
	handler := handlerMap["statusReadHandler"].(somaStatusReadHandler)
	handler.input <- somaStatusRequest{
		action: "show",
		reply:  returnChannel,
		status: somaproto.ProtoStatus{
			Status: params.ByName("status"),
		},
	}
	results := <-returnChannel
	SendStatusReply(&w, &results)
}

func AddStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	var cReq somaproto.ProtoRequestStatus
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan []somaStatusResult)
	handler := handlerMap["statusWriteHandler"].(somaStatusWriteHandler)
	handler.input <- somaStatusRequest{
		action: "add",
		reply:  returnChannel,
		status: somaproto.ProtoStatus{
			Status: cReq.Status.Status,
		},
	}
	results := <-returnChannel
	SendStatusReply(&w, &results)
}

func DeleteStatus(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaStatusResult)
	handler := handlerMap["statusWriteHandler"].(somaStatusWriteHandler)
	handler.input <- somaStatusRequest{
		action: "delete",
		reply:  returnChannel,
		status: somaproto.ProtoStatus{
			Status: params.ByName("status"),
		},
	}
	results := <-returnChannel
	SendStatusReply(&w, &results)
}

func SendStatusReply(w *http.ResponseWriter, r *[]somaStatusResult) {
	var res somaproto.ProtoResultStatus
	dispatchError := CheckErrorHandler(r, &res)
	if dispatchError {
		goto dispatch
	}
	res.Text = make([]string, 0)
	res.StatusList = make([]somaproto.ProtoStatus, 0)
	for _, l := range *r {
		res.StatusList = append(res.StatusList, l.status)
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
