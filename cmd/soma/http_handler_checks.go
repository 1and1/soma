package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

// CheckConfigurationList function
func CheckConfigurationList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `check`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["checkConfigurationReadHandler"].(*somaCheckConfigurationReadHandler)
	handler.input <- somaCheckConfigRequest{
		action: "list",
		reply:  returnChannel,
		CheckConfig: proto.CheckConfig{
			RepositoryId: params.ByName("repository"),
		},
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.Request{}
	cReq.Filter = &proto.Filter{}
	cReq.Filter.CheckConfig = &proto.CheckConfigFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.CheckConfig.Name != "" {
		filtered := []somaCheckConfigResult{}
		for _, i := range result.CheckConfigs {
			if i.CheckConfig.Name == cReq.Filter.CheckConfig.Name {
				filtered = append(filtered, i)
			}
		}
		result.CheckConfigs = filtered
	}

skip:
	SendCheckConfigurationReply(&w, &result)
}

// CheckConfigurationShow function
func CheckConfigurationShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:     params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `check`,
		Action:       `show`,
		RepositoryID: params.ByName(`repository`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["checkConfigurationReadHandler"].(*somaCheckConfigurationReadHandler)
	handler.input <- somaCheckConfigRequest{
		action: "show",
		reply:  returnChannel,
		CheckConfig: proto.CheckConfig{
			Id:           params.ByName("check"),
			RepositoryId: params.ByName("repository"),
		},
	}
	result := <-returnChannel

	SendCheckConfigurationReply(&w, &result)
}

// CheckConfigurationCreate function
func CheckConfigurationCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	cReq.CheckConfig.Id = uuid.Nil.String()

	if !IsAuthorized(&msg.Authorization{
		AuthUser:     params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `check`,
		Action:       `create`,
		RepositoryID: cReq.CheckConfig.RepositoryId,
		BucketID:     cReq.CheckConfig.BucketId,
		CapabilityID: cReq.CheckConfig.CapabilityId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "check",
		Action:      fmt.Sprintf("add_check_to_%s", cReq.CheckConfig.ObjectType),
		AuthUser:    params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		CheckConfig: somaCheckConfigRequest{
			action:      "check_configuration_new",
			CheckConfig: *cReq.CheckConfig,
		},
	}
	result := <-returnChannel
	SendCheckConfigurationReply(&w, &result)
}

// CheckConfigurationDestroy function
func CheckConfigurationDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:     params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `check`,
		Action:       `destroy`,
		RepositoryID: params.ByName(`repository`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: `check`,
		Action:      `remove_check`,
		AuthUser:    params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		CheckConfig: somaCheckConfigRequest{
			action: `check_configuration_delete`,
			CheckConfig: proto.CheckConfig{
				Id:           params.ByName(`check`),
				RepositoryId: params.ByName(`repository`),
			},
		},
	}
	result := <-returnChannel
	SendCheckConfigurationReply(&w, &result)
}

// SendCheckConfigurationReply function
func SendCheckConfigurationReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.Result{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	if result.Errors == nil {
		result.Errors = &[]string{}
	}
	result.CheckConfigs = &[]proto.CheckConfig{}
	for _, i := range (*r).CheckConfigs {
		*result.CheckConfigs = append(*result.CheckConfigs, i.CheckConfig)
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
		}
	}

dispatch:
	result.Clean()
	json, err := json.Marshal(result)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
