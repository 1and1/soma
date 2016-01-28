package main

import (
	"encoding/json"
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

	cReq := somaproto.ProtoRequestOncall{}
	cFil := somaproto.ProtoOncallFilter{}
	cReq.Filter = &cFil

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
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
func AddOncall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	var cReq somaproto.ProtoRequestOncall
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan []somaOncallResult)
	handler := handlerMap["oncallWriteHandler"].(somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "add",
		reply:  returnChannel,
		oncall: somaproto.ProtoOncall{
			Name:   cReq.OnCall.Name,
			Number: cReq.OnCall.Number,
		},
	}
	results := <-returnChannel
	SendOncallReply(&w, &results)
}

func UpdateOncall(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	var cReq somaproto.ProtoRequestOncall
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan []somaOncallResult)
	handler := handlerMap["oncallWriteHandler"].(somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "update",
		reply:  returnChannel,
		oncall: somaproto.ProtoOncall{
			Id:     params.ByName("oncall"),
			Name:   cReq.OnCall.Name,
			Number: cReq.OnCall.Number,
		},
	}
	results := <-returnChannel
	SendOncallReply(&w, &results)
}

func DeleteOncall(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaOncallResult)
	handler := handlerMap["oncallWriteHandler"].(somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "delete",
		reply:  returnChannel,
		oncall: somaproto.ProtoOncall{
			Id: params.ByName("oncall"),
		},
	}
	results := <-returnChannel
	SendOncallReply(&w, &results)
}

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
