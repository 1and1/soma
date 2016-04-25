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
		req.Custom.Name = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.TeamId = params.ByName("team")
	case "template":
		req.prType = prType
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyReadHandler"].(somaPropertyReadHandler)
	handler.input <- req
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.PropertyRequest{}
	cReq.Filter = &somaproto.PropertyFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Type == "custom") && (cReq.Filter.Name != "") &&
		(cReq.Filter.Repository != "") {
		filtered := make([]somaPropertyResult, 0)
		for _, i := range result.Properties {
			if (i.Custom.Name == cReq.Filter.Name) &&
				(i.Custom.RepositoryId == cReq.Filter.Repository) {
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
		req.Native.Name = params.ByName("native")
	case "system":
		req.prType = prType
		req.System.Name = params.ByName("system")
	case "custom":
		req.prType = prType
		req.Custom.CustomId = params.ByName("custom")
		req.Custom.RepositoryId = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.Name = params.ByName("service")
		req.Service.TeamId = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service.Name = params.ByName("service")
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

	cReq := somaproto.PropertyRequest{}
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
		if params.ByName("repository") != cReq.Custom.RepositoryId {
			DispatchBadRequest(&w, errors.New("Body and URL repositories do not match"))
			return
		}
		req.prType = prType
		req.Custom = *cReq.Custom
		req.Custom.RepositoryId = params.ByName("repository")
	case "service":
		if params.ByName("team") != cReq.Service.TeamId {
			DispatchBadRequest(&w, errors.New("Body and URL teams do not match"))
			return
		}
		req.prType = prType
		req.Service = *cReq.Service
		req.Service.TeamId = params.ByName("team")
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
		req.Native.Name = params.ByName("native")
	case "system":
		req.prType = prType
		req.System.Name = params.ByName("system")
	case "custom":
		req.prType = prType
		req.Custom.CustomId = params.ByName("custom")
		req.Custom.RepositoryId = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.Name = params.ByName("service")
		req.Service.TeamId = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service.Name = params.ByName("service")
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyWriteHandler"].(somaPropertyWriteHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

/* Utility
 */
func SendPropertyReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.PropertyResult{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Custom = make([]somaproto.TreePropertyCustom, 0)
	result.Native = make([]somaproto.TreePropertyNative, 0)
	result.Service = make([]somaproto.TreePropertyService, 0)
	result.System = make([]somaproto.TreePropertySystem, 0)
	for _, i := range (*r).Properties {
		switch i.prType {
		case "system":
			result.System = append(result.System, i.System)
		case "native":
			result.Native = append(result.Native, i.Native)
		case "custom":
			result.Custom = append(result.Custom, i.Custom)
		case "service":
			result.Service = append(result.Service, i.Service)
		case "template":
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
