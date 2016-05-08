package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListView(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewReadHandler"].(somaViewReadHandler)
	handler.input <- somaViewRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

func ShowView(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewReadHandler"].(somaViewReadHandler)
	handler.input <- somaViewRequest{
		action: "show",
		reply:  returnChannel,
		View: proto.View{
			Name: params.ByName("view"),
		},
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

/* Write functions
 */
func AddView(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.NewViewRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewWriteHandler"].(somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "add",
		reply:  returnChannel,
		View: proto.View{
			Name: cReq.View.Name,
		},
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

func DeleteView(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewWriteHandler"].(somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "delete",
		reply:  returnChannel,
		View: proto.View{
			Name: params.ByName("view"),
		},
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

func RenameView(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.NewViewRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["viewWriteHandler"].(somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "rename",
		reply:  returnChannel,
		name:   params.ByName("view"),
		View: proto.View{
			Name: cReq.View.Name,
		},
	}
	result := <-returnChannel
	SendViewReply(&w, &result)
}

/* Utility
 */
func SendViewReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewViewResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Views {
		*result.Views = append(*result.Views, i.View)
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
