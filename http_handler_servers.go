package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListServer(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
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

func ShowServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverReadHandler"].(somaServerReadHandler)
	handler.input <- somaServerRequest{
		action: "show",
		reply:  returnChannel,
		Server: somaproto.ProtoServer{
			Id: params.ByName("server"),
		},
	}
	result := <-returnChannel
	SendServerReply(&w, &result)
}

/* Write functions
 */
func AddServer(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestServer{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverWriteHandler"].(somaServerWriteHandler)
	handler.input <- somaServerRequest{
		action: "add",
		reply:  returnChannel,
		Server: somaproto.ProtoServer{
			AssetId:    cReq.Server.AssetId,
			Datacenter: cReq.Server.Datacenter,
			Location:   cReq.Server.Location,
			Name:       cReq.Server.Name,
			IsOnline:   cReq.Server.IsOnline,
			IsDeleted:  false,
		},
	}
	result := <-returnChannel
	SendServerReply(&w, &result)
}

func DeleteServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	action := "delete"

	cReq := somaproto.ProtoRequestServer{}
	_ = DecodeJsonBody(r, &cReq)
	if cReq.Purge {
		action = "purge"
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverWriteHandler"].(somaServerWriteHandler)
	handler.input <- somaServerRequest{
		action: action,
		reply:  returnChannel,
		Server: somaproto.ProtoServer{
			Id: params.ByName("server"),
		},
	}
	result := <-returnChannel
	SendServerReply(&w, &result)
}

func InsertNullServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestServer{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if cReq.Server.Id != "00000000-0000-0000-0000-000000000000" ||
		params.ByName("server") != "null" {
		DispatchBadRequest(&w, errors.New("not null server"))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverWriteHandler"].(somaServerWriteHandler)
	handler.input <- somaServerRequest{
		action: "insert-null",
		reply:  returnChannel,
		Server: somaproto.ProtoServer{
			Datacenter: cReq.Server.Datacenter,
		},
	}
	result := <-returnChannel
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
