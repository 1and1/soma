package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListOncall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaOncallResult)
	handler := handlerMap["oncallReadHandler"].(somaOncallReadHandler)
	handler.input <- somaOncallRequest{
		action: "list",
		reply:  returnChannel,
	}
	results := <-returnChannel

	var cReq somaproto.ProtoRequestOncall
	b, _ := ioutil.ReadAll((*r).Body)
	if string(b) != "" {
		err := DecodeJsonBody(r, &cReq)
		if err != nil || cReq.Filter.Name == "" {
			DispatchBadRequest(&w, err)
			return
		}

		filtered := make([]somaOncallResult, 0)
	filterloop:
		for _, onCall := range results {
			if onCall.rErr != nil {
				filtered = append(filtered, onCall)
				break filterloop
			}
			if onCall.oncall.Name == cReq.Filter.Name {
				filtered = append(filtered, onCall)
			}
		}
		results = filtered
	}

	SendOncallReply(&w, &results)
}

func ShowOncall(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaOncallResult)
	handler := handlerMap["oncallReadHandler"].(somaOncallReadHandler)
	handler.input <- somaOncallRequest{
		action: "show",
		reply:  returnChannel,
		oncall: somaproto.ProtoOncall{
			Id: params.ByName("oncall"),
		},
	}
	results := <-returnChannel
	SendOncallReply(&w, &results)
}

/*
 * Write functions
 */
/*
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
*/

/*
 * Utility
 */
func SendOncallReply(w *http.ResponseWriter, r *[]somaOncallResult) {
	var res somaproto.ProtoResultOncall
	dispatchError := CheckErrorHandler(r, &res)
	if dispatchError {
		goto dispatch
	}
	res.Text = make([]string, 0)
	res.Oncalls = make([]somaproto.ProtoOncall, 0)
	for _, l := range *r {
		res.Oncalls = append(res.Oncalls, l.oncall)
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
