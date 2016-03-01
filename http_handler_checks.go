package main

import (
	"encoding/json"
	"fmt"
	"net/http"


	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListCheckConfiguration(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["checkConfigurationReadHandler"].(somaCheckConfigurationReadHandler)
	handler.input <- somaCheckConfigRequest{
		action: "list",
		reply:  returnChannel,
		CheckConfig: somaproto.CheckConfiguration{
			RepositoryId: params.ByName("repository"),
		},
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.CheckConfigurationRequest{}
	cReq.Filter = &somaproto.CheckConfigurationFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaCheckConfigResult, 0)
		for _, i := range result.CheckConfigs {
			if i.CheckConfig.Name == cReq.Filter.Name {
				filtered = append(filtered, i)
			}
		}
		result.CheckConfigs = filtered
	}

skip:
	SendCheckConfigurationReply(&w, &result)
}

func ShowCheckConfiguration(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["checkConfigurationReadHandler"].(somaCheckConfigurationReadHandler)
	handler.input <- somaCheckConfigRequest{
		action: "show",
		reply:  returnChannel,
		CheckConfig: somaproto.CheckConfiguration{
			Id:           params.ByName("check"),
			RepositoryId: params.ByName("repository"),
		},
	}
	result := <-returnChannel

	SendCheckConfigurationReply(&w, &result)
}

/* Write functions
 */
func AddCheckConfiguration(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.CheckConfigurationRequest{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
	handler.input <- treeRequest{
		RequestType: "check",
		Action:      fmt.Sprintf("add_check_to_%s", cReq.CheckConfiguration.ObjectType),
		reply:       returnChannel,
		CheckConfig: somaCheckConfigRequest{
			action:      "check_configuration_new",
			CheckConfig: *cReq.CheckConfiguration,
		},
	}
	result := <-returnChannel
	SendCheckConfigurationReply(&w, &result)
}

/* Utility
 */
func SendCheckConfigurationReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.CheckConfigurationResult{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.CheckConfigurations = make([]somaproto.CheckConfiguration, 0)
	for _, i := range (*r).CheckConfigs {
		result.CheckConfigurations = append(result.CheckConfigurations, i.CheckConfig)
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
