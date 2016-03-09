package main

import (
	"encoding/json"
	"net/http"
	"strconv"


	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

/* Read functions
 */
func GetHostDeployment(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	var (
		err     error
		assetid int64
	)

	if _, err = uuid.FromString(params.ByName("system")); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if assetid, err = strconv.ParseInt(params.ByName("assetid"), 10, 64); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	cReq := somaproto.HostDeploymentRequest{}
	if err = DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["hostDeploymentHandler"].(somaHostDeploymentHandler)
	handler.input <- somaHostDeploymentRequest{
		action:  "get",
		reply:   returnChannel,
		system:  params.ByName("system"),
		assetid: assetid,
		idlist:  cReq.IdList,
	}
	result := <-returnChannel
	SendHostDeploymentReply(&w, &result)
}

func AssembleHostUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	var (
		err     error
		assetid int64
	)

	if _, err = uuid.FromString(params.ByName("system")); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if assetid, err = strconv.ParseInt(params.ByName("assetid"), 10, 64); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["hostDeploymentHandler"].(somaHostDeploymentHandler)
	handler.input <- somaHostDeploymentRequest{
		action:  "assemble",
		reply:   returnChannel,
		system:  params.ByName("system"),
		assetid: assetid,
	}
	result := <-returnChannel
	SendHostDeploymentReply(&w, &result)
}

/* Utility
 */
func SendHostDeploymentReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.HostDeploymentResult{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Deployments = make([]somaproto.DeploymentDetails, 0)
	for _, i := range (*r).HostDeployments {
		if i.Delete {
			result.Delete = append(result.Delete, i.DeleteId)
			continue
		}
		result.Deployments = append(result.Deployments, i.Deployment)
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
