package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListLevel(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declase here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestLevel{}
	cReq.Filter = &somaproto.ProtoLevelFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Name != "") || (cReq.Filter.ShortName != "") {
		filtered := make([]somaLevelResult, 0)
		for _, i := range result.Levels {
			if ((cReq.Filter.Name != "") && (cReq.Filter.Name == i.Level.Name)) ||
				((cReq.Filter.ShortName != "") && (cReq.Filter.ShortName == i.Level.ShortName)) {
				filtered = append(filtered, i)
			}
		}
		result.Levels = filtered
	}

skip:
	SendLevelReply(&w, &result)
}

func ShowLevel(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "show",
		reply:  returnChannel,
		Level: somaproto.ProtoLevel{
			Name: params.ByName("level"),
		},
	}
	result := <-returnChannel
	SendLevelReply(&w, &result)
}

/* Write functions
 */
func AddLevel(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestLevel{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelWriteHandler"].(somaLevelWriteHandler)
	handler.input <- somaLevelRequest{
		action: "add",
		reply:  returnChannel,
		Level: somaproto.ProtoLevel{
			Name:      cReq.Level.Name,
			ShortName: cReq.Level.ShortName,
			Numeric:   cReq.Level.Numeric,
		},
	}
	result := <-returnChannel
	SendLevelReply(&w, &result)
}

func DeleteLevel(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelWriteHandler"].(somaLevelWriteHandler)
	handler.input <- somaLevelRequest{
		action: "delete",
		reply:  returnChannel,
		Level: somaproto.ProtoLevel{
			Name: params.ByName("level"),
		},
	}
	result := <-returnChannel
	SendLevelReply(&w, &result)
}

/* Utility
 */
func SendLevelReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultLevel{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Levels = make([]somaproto.ProtoLevel, 0)
	for _, i := range (*r).Levels {
		result.Levels = append(result.Levels, i.Level)
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
