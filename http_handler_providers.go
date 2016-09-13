package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListProvider(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`providers_list`, ``, ``, ``); !ok {
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

func ShowProvider(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`providers_show`, ``, ``, ``); !ok {
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

/* Write functions
 */
func AddProvider(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`providers_create`, ``, ``, ``); !ok {
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

func DeleteProvider(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`providers_delete`, ``, ``, ``); !ok {
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

/* Utility
 */
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
