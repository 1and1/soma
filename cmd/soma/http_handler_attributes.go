package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

/*Read functions
 */
func ListAttribute(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `attribute`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["attributeReadHandler"].(*somaAttributeReadHandler)
	handler.input <- somaAttributeRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendAttributeReply(&w, &result)
}

func ShowAttribute(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `attribute`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["attributeReadHandler"].(*somaAttributeReadHandler)
	handler.input <- somaAttributeRequest{
		action: "show",
		reply:  returnChannel,
		Attribute: proto.Attribute{
			Name: params.ByName("attribute"),
		},
	}
	result := <-returnChannel
	SendAttributeReply(&w, &result)
}

/* Write functions
 */
func AddAttribute(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `attribute`,
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
	handler := handlerMap["attributeWriteHandler"].(*somaAttributeWriteHandler)
	handler.input <- somaAttributeRequest{
		action: "add",
		reply:  returnChannel,
		Attribute: proto.Attribute{
			Name:        cReq.Attribute.Name,
			Cardinality: cReq.Attribute.Cardinality,
		},
	}
	result := <-returnChannel
	SendAttributeReply(&w, &result)
}

func DeleteAttribute(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `attribute`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["attributeWriteHandler"].(*somaAttributeWriteHandler)
	handler.input <- somaAttributeRequest{
		action: "delete",
		reply:  returnChannel,
		Attribute: proto.Attribute{
			Name: params.ByName("attribute"),
		},
	}
	result := <-returnChannel
	SendAttributeReply(&w, &result)
}

/* Utility
 */
func SendAttributeReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.Result{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	if result.Errors == nil {
		result.Errors = &[]string{}
	}
	result.Attributes = &[]proto.Attribute{}
	for _, i := range (*r).Attributes {
		*result.Attributes = append(*result.Attributes, i.Attribute)
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
