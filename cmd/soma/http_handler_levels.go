package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// LevelList function
func LevelList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `level`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelReadHandler"].(*somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declared here since goto does not jump over declarations
	cReq := proto.NewLevelFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Level.Name != "") ||
		(cReq.Filter.Level.ShortName != "") {
		filtered := []somaLevelResult{}
		for _, i := range result.Levels {
			if ((cReq.Filter.Level.Name != "") &&
				(cReq.Filter.Level.Name == i.Level.Name)) ||
				((cReq.Filter.Level.ShortName != "") &&
					(cReq.Filter.Level.ShortName == i.Level.ShortName)) {
				filtered = append(filtered, i)
			}
		}
		result.Levels = filtered
	}

skip:
	SendLevelReply(&w, &result)
}

// LevelShow function
func LevelShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `level`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelReadHandler"].(*somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "show",
		reply:  returnChannel,
		Level: proto.Level{
			Name: params.ByName("level"),
		},
	}
	result := <-returnChannel
	SendLevelReply(&w, &result)
}

// LevelAdd function
func LevelAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `level`,
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
	handler := handlerMap["levelWriteHandler"].(*somaLevelWriteHandler)
	handler.input <- somaLevelRequest{
		action: "add",
		reply:  returnChannel,
		Level: proto.Level{
			Name:      cReq.Level.Name,
			ShortName: cReq.Level.ShortName,
			Numeric:   cReq.Level.Numeric,
		},
	}
	result := <-returnChannel
	SendLevelReply(&w, &result)
}

// LevelRemove functions
func LevelRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `level`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelWriteHandler"].(*somaLevelWriteHandler)
	handler.input <- somaLevelRequest{
		action: "delete",
		reply:  returnChannel,
		Level: proto.Level{
			Name: params.ByName("level"),
		},
	}
	result := <-returnChannel
	SendLevelReply(&w, &result)
}

// SendLevelReply function
func SendLevelReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewLevelResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Levels {
		*result.Levels = append(*result.Levels, i.Level)
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
