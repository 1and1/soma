package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListLevel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaLevelResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "list",
		reply:  returnChannel,
	}
	results := <-returnChannel
	SendLevelReply(&w, &results)
}

func ShowLevel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaLevelResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "show",
		reply:  returnChannel,
		level: somaproto.ProtoLevel{
			Name: params.ByName("level"),
		},
	}
	results := <-returnChannel
	SendLevelReply(&w, &results)
}

func AddLevel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	var cReq somaproto.ProtoRequestLevel
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan []somaLevelResult)
	handler := handlerMap["levelWriteHandler"].(somaLevelWriteHandler)
	handler.input <- somaLevelRequest{
		action: "add",
		reply:  returnChannel,
		level: somaproto.ProtoLevel{
			Name:      cReq.Level.Name,
			ShortName: cReq.Level.ShortName,
			Numeric:   cReq.Level.Numeric,
		},
	}
	results := <-returnChannel
	SendLevelReply(&w, &results)
}

func DeleteLevel(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaLevelResult)
	handler := handlerMap["levelWriteHandler"].(somaLevelWriteHandler)
	handler.input <- somaLevelRequest{
		action: "delete",
		reply:  returnChannel,
		level: somaproto.ProtoLevel{
			Name: params.ByName("level"),
		},
	}
	results := <-returnChannel
	SendLevelReply(&w, &results)
}

func SendLevelReply(w *http.ResponseWriter, r *[]somaLevelResult) {
	var res somaproto.ProtoResultLevel
	dispatchError := CheckErrorHandler(r, &res)
	if dispatchError {
		goto dispatch
	}
	res.Text = make([]string, 0)
	res.Levels = make([]somaproto.ProtoLevel, 0)
	for _, l := range *r {
		res.Levels = append(res.Levels, l.level)
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
