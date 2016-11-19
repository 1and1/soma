package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// ProviderList function
func ProviderList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `provider`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["providerReadHandler"].(*somaProviderReadHandler)
	handler.input <- somaProviderRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendProviderReply(&w, &result)
}

// ProviderShow function
func ProviderShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `provider`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["providerReadHandler"].(*somaProviderReadHandler)
	handler.input <- somaProviderRequest{
		action: "show",
		reply:  returnChannel,
		Provider: proto.Provider{
			Name: params.ByName("provider"),
		},
	}
	result := <-returnChannel
	SendProviderReply(&w, &result)
}

// ProviderAdd function
func ProviderAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `provider`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["providerWriteHandler"].(*somaProviderWriteHandler)
	handler.input <- somaProviderRequest{
		action: "add",
		reply:  returnChannel,
		Provider: proto.Provider{
			Name: cReq.Provider.Name,
		},
	}
	result := <-returnChannel
	SendProviderReply(&w, &result)
}

// ProviderRemove function
func ProviderRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `provider`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["providerWriteHandler"].(*somaProviderWriteHandler)
	handler.input <- somaProviderRequest{
		action: "delete",
		reply:  returnChannel,
		Provider: proto.Provider{
			Name: params.ByName("provider"),
		},
	}
	result := <-returnChannel
	SendProviderReply(&w, &result)
}

// SendProviderReply function
func SendProviderReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewProviderResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Providers {
		*result.Providers = append(*result.Providers, i.Provider)
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
