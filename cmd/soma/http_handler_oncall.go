package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListOncall(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`oncall_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["oncallReadHandler"].(*somaOncallReadHandler)
	handler.input <- somaOncallRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.NewOncallFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Oncall.Name != "" {
		filtered := make([]somaOncallResult, 0)
		for _, i := range result.Oncall {
			if i.Oncall.Name == cReq.Filter.Oncall.Name {
				filtered = append(filtered, i)
			}
		}
		result.Oncall = filtered
	}

skip:
	SendOncallReply(&w, &result)
}

func ShowOncall(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`oncall_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["oncallReadHandler"].(*somaOncallReadHandler)
	handler.input <- somaOncallRequest{
		action: "show",
		reply:  returnChannel,
		Oncall: proto.Oncall{
			Id: params.ByName("oncall"),
		},
	}
	result := <-returnChannel
	SendOncallReply(&w, &result)
}

/* Write functions
 */
func AddOncall(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`oncall_create`, ``, ``, ``); !ok {
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
	handler := handlerMap["oncallWriteHandler"].(*somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "add",
		reply:  returnChannel,
		Oncall: proto.Oncall{
			Name:   cReq.Oncall.Name,
			Number: cReq.Oncall.Number,
		},
	}
	result := <-returnChannel
	SendOncallReply(&w, &result)
}

func UpdateOncall(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`oncall_update`, ``, ``, ``); !ok {
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
	handler := handlerMap["oncallWriteHandler"].(*somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "update",
		reply:  returnChannel,
		Oncall: proto.Oncall{
			Id:     params.ByName("oncall"),
			Name:   cReq.Oncall.Name,
			Number: cReq.Oncall.Number,
		},
	}
	result := <-returnChannel
	SendOncallReply(&w, &result)
}

func DeleteOncall(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`oncall_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["oncallWriteHandler"].(*somaOncallWriteHandler)
	handler.input <- somaOncallRequest{
		action: "delete",
		reply:  returnChannel,
		Oncall: proto.Oncall{
			Id: params.ByName("oncall"),
		},
	}
	result := <-returnChannel
	SendOncallReply(&w, &result)
}

/* Utility
 */
func SendOncallReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewOncallResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Oncall {
		*result.Oncalls = append(*result.Oncalls, i.Oncall)
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
