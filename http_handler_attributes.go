package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*Read functions
 */
func ListAttribute(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["attributeReadHandler"].(somaAttributeReadHandler)
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

	returnChannel := make(chan somaResult)
	handler := handlerMap["attributeReadHandler"].(somaAttributeReadHandler)
	handler.input <- somaAttributeRequest{
		action: "show",
		reply:  returnChannel,
		Attribute: somaproto.ProtoAttribute{
			Attribute: params.ByName("attribute"),
		},
	}
	result := <-returnChannel
	SendAttributeReply(&w, &result)
}

/* Write functions
 */
func AddAttribute(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestAttribute{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["attributeWriteHandler"].(somaAttributeWriteHandler)
	handler.input <- somaAttributeRequest{
		action: "add",
		reply:  returnChannel,
		Attribute: somaproto.ProtoAttribute{
			Attribute: cReq.Attribute.Attribute,
		},
	}
	result := <-returnChannel
	SendAttributeReply(&w, &result)
}

func DeleteAttribute(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["attributeWriteHandler"].(somaAttributeWriteHandler)
	handler.input <- somaAttributeRequest{
		action: "delete",
		reply:  returnChannel,
		Attribute: somaproto.ProtoAttribute{
			Attribute: params.ByName("attribute"),
		},
	}
	result := <-returnChannel
	SendAttributeReply(&w, &result)
}

/* Utility
 */
func SendAttributeReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultAttribute{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Attributes = make([]somaproto.ProtoAttribute, 0)
	for _, i := range (*r).Attributes {
		result.Attributes = append(result.Attributes, i.Attribute)
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
