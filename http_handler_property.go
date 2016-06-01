package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	prType, _ := GetPropertyTypeFromUrl(r.URL)
	pa := fmt.Sprintf("property_%s_list", prType)
	switch prType {
	case `custom`:
	case `service`:
	default:
		if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
			pa, ``, ``, ``); !ok {
			DispatchForbidden(&w, nil)
			return
		}
	}

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "list",
		reply:  returnChannel,
	}
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
	cReq := proto.NewPropertyFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Property.Type == "custom") && (cReq.Filter.Property.Name != "") &&
		(cReq.Filter.Property.RepositoryId != "") {
		filtered := make([]somaPropertyResult, 0)
		for _, i := range result.Properties {
			if (i.Custom.Name == cReq.Filter.Property.Name) &&
				(i.Custom.RepositoryId == cReq.Filter.Property.RepositoryId) {
				filtered = append(filtered, i)
			}
		}
		result.Properties = filtered
	}
	if (cReq.Filter.Property.Type == "system") && (cReq.Filter.Property.Name != "") {
		filtered := make([]somaPropertyResult, 0)
		for _, i := range result.Properties {
			if i.System.Name == cReq.Filter.Property.Name {
				filtered = append(filtered, i)
			}
		}
		result.Properties = filtered
	}
	if (cReq.Filter.Property.Type == "service") && (cReq.Filter.Property.Name != "") {
		filtered := make([]somaPropertyResult, 0)
		for _, i := range result.Properties {
			if (i.Service.Name == cReq.Filter.Property.Name) &&
				(i.Service.TeamId == params.ByName("team")) {
				filtered = append(filtered, i)
			}
		}
		result.Properties = filtered
	}
	if (cReq.Filter.Property.Type == "template") && (cReq.Filter.Property.Name != "") {
		filtered := make([]somaPropertyResult, 0)
		for _, i := range result.Properties {
			if i.Service.Name == cReq.Filter.Property.Name {
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
	prType, _ := GetPropertyTypeFromUrl(r.URL)
	pa := fmt.Sprintf("property_%s_show", prType)
	switch prType {
	case `custom`:
	case `service`:
	default:
		if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
			pa, ``, ``, ``); !ok {
			DispatchForbidden(&w, nil)
			return
		}
	}

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "show",
		reply:  returnChannel,
	}
	switch prType {
	case "native":
		req.prType = prType
		req.Native.Name = params.ByName("native")
	case "system":
		req.prType = prType
		req.System.Name = params.ByName("system")
	case "custom":
		req.prType = prType
		req.Custom.Id = params.ByName("custom")
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
	prType, _ := GetPropertyTypeFromUrl(r.URL)
	pa := fmt.Sprintf("property_%s_create", prType)
	switch prType {
	case `custom`:
	case `service`:
	default:
		if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
			pa, ``, ``, ``); !ok {
			DispatchForbidden(&w, nil)
			return
		}
	}

	cReq := proto.NewPropertyRequest()
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
	switch prType {
	case "native":
		req.prType = prType
		req.Native = *cReq.Property.Native
	case "system":
		req.prType = prType
		req.System = *cReq.Property.System
	case "custom":
		if params.ByName("repository") != cReq.Property.Custom.RepositoryId {
			DispatchBadRequest(&w, errors.New("Body and URL repositories do not match"))
			return
		}
		req.prType = prType
		req.Custom = *cReq.Property.Custom
		req.Custom.RepositoryId = params.ByName("repository")
	case "service":
		if params.ByName("team") != cReq.Property.Service.TeamId {
			DispatchBadRequest(&w, errors.New("Body and URL teams do not match"))
			return
		}
		req.prType = prType
		req.Service = *cReq.Property.Service
		req.Service.TeamId = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service = *cReq.Property.Service
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
	prType, _ := GetPropertyTypeFromUrl(r.URL)
	pa := fmt.Sprintf("property_%s_delete", prType)
	switch prType {
	case `custom`:
	case `service`:
	default:
		if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
			pa, ``, ``, ``); !ok {
			DispatchForbidden(&w, nil)
			return
		}
	}

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "delete",
		reply:  returnChannel,
	}
	switch prType {
	case "native":
		req.prType = prType
		req.Native.Name = params.ByName("native")
	case "system":
		req.prType = prType
		req.System.Name = params.ByName("system")
	case "custom":
		req.prType = prType
		req.Custom.Id = params.ByName("custom")
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
	result := proto.NewPropertyResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Properties {
		switch i.prType {
		case "system":
			*result.Properties = append(*result.Properties, proto.Property{Type: "system",
				System: &proto.PropertySystem{
					Name:  i.System.Name,
					Value: i.System.Value,
				}})
		case "native":
			*result.Properties = append(*result.Properties, proto.Property{Type: "native",
				Native: &proto.PropertyNative{
					Name:  i.Native.Name,
					Value: i.Native.Value,
				}})
		case "custom":
			*result.Properties = append(*result.Properties, proto.Property{Type: "custom",
				Custom: &proto.PropertyCustom{
					Id:           i.Custom.Id,
					Name:         i.Custom.Name,
					Value:        i.Custom.Value,
					RepositoryId: i.Custom.RepositoryId,
				}})
		case "service":
			prop := proto.Property{
				Type: "service",
				Service: &proto.PropertyService{
					Name:       i.Service.Name,
					TeamId:     i.Service.TeamId,
					Attributes: []proto.ServiceAttribute{},
				}}
			for _, a := range i.Service.Attributes {
				prop.Service.Attributes = append(prop.Service.Attributes, proto.ServiceAttribute{
					Name:  a.Name,
					Value: a.Value,
				})
			}
			*result.Properties = append(*result.Properties, prop)
		case "template":
			prop := proto.Property{
				Type: "template",
				Service: &proto.PropertyService{
					Name:       i.Service.Name,
					Attributes: []proto.ServiceAttribute{},
				}}
			for _, a := range i.Service.Attributes {
				prop.Service.Attributes = append(prop.Service.Attributes, proto.ServiceAttribute{
					Name:  a.Name,
					Value: a.Value,
				})
			}
			*result.Properties = append(*result.Properties, prop)
		}
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
