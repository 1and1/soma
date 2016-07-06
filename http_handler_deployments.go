package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

/* Write functions
 */
func DeliverDeploymentDetails(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if _, err := uuid.FromString(params.ByName("uuid")); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["deploymentHandler"].(somaDeploymentHandler)
	handler.input <- somaDeploymentRequest{
		action:     "get",
		reply:      returnChannel,
		Deployment: params.ByName("uuid"),
	}
	result := <-returnChannel
	SendDeploymentReply(&w, &result)
}

func DeliverMonitoringDeployments(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	var action string

	if _, err := uuid.FromString(params.ByName("uuid")); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if params.ByName("all") == "" {
		action = "list"
	} else {
		action = "listall"
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["deploymentHandler"].(somaDeploymentHandler)
	handler.input <- somaDeploymentRequest{
		action:     action,
		reply:      returnChannel,
		Deployment: params.ByName("uuid"),
	}
	result := <-returnChannel
	SendDeploymentReply(&w, &result)
}

func UpdateDeploymentDetails(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if _, err := uuid.FromString(params.ByName("uuid")); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	switch params.ByName("result") {
	case "success":
	case "failed":
	default:
		DispatchBadRequest(&w, fmt.Errorf("Unknown result: %s", params.ByName("result")))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["deploymentHandler"].(somaDeploymentHandler)
	handler.input <- somaDeploymentRequest{
		action:     fmt.Sprintf("update/%s", params.ByName("result")),
		reply:      returnChannel,
		Deployment: params.ByName("uuid"),
	}
	result := <-returnChannel
	SendDeploymentReply(&w, &result)
}

/* Utility
 */
func SendDeploymentReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewDeploymentResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Deployments = &[]proto.Deployment{}
	result.DeploymentsList = &[]string{}
	for _, i := range (*r).Deployments {
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
		}
		if i.ListEntry != "" {
			*result.DeploymentsList = append(*result.DeploymentsList, i.ListEntry)
			continue
		}
		*result.Deployments = append(*result.Deployments, i.Deployment)
	}
	if len(*result.Deployments) > 0 && len(*result.DeploymentsList) == 0 {
		result.DeploymentsList = nil
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
