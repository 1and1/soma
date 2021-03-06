package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`servers_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverReadHandler"].(*somaServerReadHandler)
	handler.input <- somaServerRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	SendServerReply(&w, &result)
}

func SyncServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`servers_sync`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverReadHandler"].(*somaServerReadHandler)
	handler.input <- somaServerRequest{
		action: "sync",
		reply:  returnChannel,
	}
	result := <-returnChannel

	SendServerReply(&w, &result)
}

func ShowServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`servers_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverReadHandler"].(*somaServerReadHandler)
	handler.input <- somaServerRequest{
		action: "show",
		reply:  returnChannel,
		Server: proto.Server{
			Id: params.ByName("server"),
		},
	}
	result := <-returnChannel
	SendServerReply(&w, &result)
}

func SearchServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`servers_search`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverReadHandler"].(*somaServerReadHandler)
	ssr := somaServerRequest{
		reply: returnChannel,
	}
	cReq := proto.NewServerFilter()
	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Server.Name != "" {
		ssr.action = "search/name"
		ssr.Filter = proto.Filter{Server: &proto.ServerFilter{
			Name: cReq.Filter.Server.Name,
		}}
	}
	if cReq.Filter.Server.AssetId != 0 {
		ssr.action = "search/asset"
		ssr.Filter = proto.Filter{Server: &proto.ServerFilter{
			AssetId: cReq.Filter.Server.AssetId,
		}}
	}

	handler.input <- ssr
	result := <-returnChannel

	SendServerReply(&w, &result)
}

/* Write functions
 */
func AddServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`servers_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewServerRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverWriteHandler"].(*somaServerWriteHandler)
	handler.input <- somaServerRequest{
		action: "add",
		reply:  returnChannel,
		Server: proto.Server{
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
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`servers_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewServerRequest()
	_ = DecodeJsonBody(r, &cReq)
	if cReq.Flags.Purge {
		action = "purge"
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverWriteHandler"].(*somaServerWriteHandler)
	handler.input <- somaServerRequest{
		action: action,
		reply:  returnChannel,
		Server: proto.Server{
			Id: params.ByName("server"),
		},
	}
	result := <-returnChannel
	SendServerReply(&w, &result)
}

func UpdateServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`servers_update`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewServerRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if cReq.Server.Id != params.ByName(`server`) {
		DispatchBadRequest(&w, errors.New(`Mismatching server UUIDs`))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["serverWriteHandler"].(*somaServerWriteHandler)
	handler.input <- somaServerRequest{
		action: "update",
		reply:  returnChannel,
		Server: proto.Server{
			Id:         cReq.Server.Id,
			AssetId:    cReq.Server.AssetId,
			Datacenter: cReq.Server.Datacenter,
			Location:   cReq.Server.Location,
			Name:       cReq.Server.Name,
			IsOnline:   cReq.Server.IsOnline,
			IsDeleted:  cReq.Server.IsDeleted,
		},
	}
	result := <-returnChannel
	SendServerReply(&w, &result)
}

func InsertNullServer(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`servers_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewServerRequest()
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
	handler := handlerMap["serverWriteHandler"].(*somaServerWriteHandler)
	handler.input <- somaServerRequest{
		action: "insert-null",
		reply:  returnChannel,
		Server: proto.Server{
			Datacenter: cReq.Server.Datacenter,
		},
	}
	result := <-returnChannel
	SendServerReply(&w, &result)
}

/* Utility
 */
func SendServerReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewServerResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Servers {
		*result.Servers = append(*result.Servers, i.Server)
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
