package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "list",
		reply:  returnChannel,
	}
	prType, _ := GetPropertyTypeFromUrl(r.URL)
	switch prType {
	case "native":
		req.prType = prType
	case "system":
		req.prType = prType
	case "custom":
		req.prType = prType
		req.Custom.Repository = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.Team = params.ByName("team")
	case "template":
		req.prType = prType
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyReadHandler"].(somaPropertyReadHandler)
	handler.input <- req
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestProperty{}
	cReq.Filter = &somaproto.ProtoPropertyFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Type == "custom") && (cReq.Filter.Property != "") &&
		(cReq.Filter.Repository != "") {
		filtered := make([]somaPropertyResult, 0)
		for _, i := range result.Properties {
			if (i.Custom.Property == cReq.Filter.Property) &&
				(i.Custom.Repository == cReq.Filter.Repository) {
				filtered = append(filtered, i)
			}
		}
		result.Properties = filtered
	}

skip:
	SendPropertyReply(&w, &result)
}

func ShowProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "show",
		reply:  returnChannel,
	}
	prType, _ := GetPropertyTypeFromUrl(r.URL)
	switch prType {
	case "native":
		req.prType = prType
		req.Native.Property = params.ByName("native")
	case "system":
		req.prType = prType
		req.System.Property = params.ByName("system")
	case "custom":
		req.prType = prType
		req.Custom.Id = params.ByName("custom")
		req.Custom.Repository = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.Property = params.ByName("service")
		req.Service.Team = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service.Property = params.ByName("service")
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyReadHandler"].(somaPropertyReadHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

/* Write functions
 */
func AddProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestProperty{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "add",
		reply:  returnChannel,
	}
	prType, _ := GetPropertyTypeFromUrl(r.URL)
	switch prType {
	case "native":
		req.prType = prType
		req.Native = *cReq.Native
	case "system":
		req.prType = prType
		req.System = *cReq.System
	case "custom":
		if params.ByName("repository") != cReq.Custom.Repository {
			DispatchBadRequest(&w, errors.New("Body and URL repositories do not match"))
			return
		}
		req.prType = prType
		req.Custom = *cReq.Custom
		req.Custom.Repository = params.ByName("repository")
	case "service":
		if params.ByName("team") != cReq.Service.Team {
			DispatchBadRequest(&w, errors.New("Body and URL teams do not match"))
			return
		}
		req.prType = prType
		req.Service = *cReq.Service
		req.Service.Team = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service = *cReq.Service
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyWriteHandler"].(somaPropertyWriteHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

func DeleteProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "delete",
		reply:  returnChannel,
	}
	prType, _ := GetPropertyTypeFromUrl(r.URL)
	switch prType {
	case "native":
		req.prType = prType
		req.Native.Property = params.ByName("native")
	case "system":
		req.prType = prType
		req.System.Property = params.ByName("system")
	case "custom":
		req.prType = prType
		req.Custom.Id = params.ByName("custom")
		req.Custom.Repository = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.Property = params.ByName("service")
		req.Service.Team = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service.Property = params.ByName("service")
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyReadHandler"].(somaPropertyReadHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

/* Utility
 */
func SendPropertyReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultProperty{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Custom = make([]somaproto.ProtoPropertyCustom, 0)
	result.Native = make([]somaproto.ProtoPropertyNative, 0)
	result.Service = make([]somaproto.ProtoPropertyService, 0)
	result.System = make([]somaproto.ProtoPropertySystem, 0)
	for _, i := range (*r).Properties {
		switch i.prType {
		case "system":
			result.System = append(result.System, i.System)
		case "native":
			result.Native = append(result.Native, i.Native)
		case "custom":
			result.Custom = append(result.Custom, i.Custom)
		default:
			result.Service = append(result.Service, i.Service)
		}
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
