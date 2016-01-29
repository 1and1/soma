package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverReadHandler"].(somaServerReadHandler)
	handler.input <- somaServerRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestServer{}
	cReq.Filter = &somaproto.ProtoServerFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaServerResult, 0)
		for _, i := range result.Servers {
			if i.Server.Name == cReq.Filter.Name {
				filtered = append(filtered, i)
			}
		}
		result.Servers = filtered
	}

skip:
	SendServerReply(&w, &result)
}

/* Utility
 */
func SendServerReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultServer{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Servers = make([]somaproto.ProtoServer, 0)
	for _, i := range (*r).Servers {
		result.Servers = append(result.Servers, i.Server)
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
