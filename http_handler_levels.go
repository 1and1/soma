package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListLevel(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`levels_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declase here since goto does not jump over declarations
	cReq := proto.NewLevelFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Level.Name != "") || (cReq.Filter.Level.ShortName != "") {
		filtered := make([]somaLevelResult, 0)
		for _, i := range result.Levels {
			if ((cReq.Filter.Level.Name != "") && (cReq.Filter.Level.Name == i.Level.Name)) ||
				((cReq.Filter.Level.ShortName != "") && (cReq.Filter.Level.ShortName == i.Level.ShortName)) {
				filtered = append(filtered, i)
			}
		}
		result.Levels = filtered
	}

skip:
	SendLevelReply(&w, &result)
}

func ShowLevel(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`levels_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
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

/* Write functions
 */
func AddLevel(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`levels_create`, ``, ``, ``); !ok {
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
	handler := handlerMap["levelWriteHandler"].(somaLevelWriteHandler)
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

func DeleteLevel(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`levels_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["levelWriteHandler"].(somaLevelWriteHandler)
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

/* Utility
 */
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
