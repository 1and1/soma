package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListCluster(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["clusterReadHandler"].(somaClusterReadHandler)
	handler.input <- somaClusterRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestCluster{}
	cReq.Filter = &somaproto.ProtoClusterFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaClusterResult, 0)
		for _, i := range result.Clusters {
			if i.Cluster.Name == cReq.Filter.Name {
				filtered = append(filtered, i)
			}
		}
		result.Clusters = filtered
	}

skip:
	SendClusterReply(&w, &result)
}

func ShowCluster(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["clusterReadHandler"].(somaClusterReadHandler)
	handler.input <- somaClusterRequest{
		action: "show",
		reply:  returnChannel,
		Cluster: somaproto.ProtoCluster{
			Id: params.ByName("cluster"),
		},
	}
	result := <-returnChannel
	SendClusterReply(&w, &result)
}

/* Write functions
 */
func AddCluster(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestCluster{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
	handler.input <- treeRequest{
		RequestType: "cluster",
		Action:      "add",
		Cluster: somaClusterRequest{
			action:  "add",
			Cluster: *cReq.Cluster,
			reply:   returnChannel,
		},
	}
	result := <-returnChannel
	SendClusterReply(&w, &result)
}

/*
 * Utility
 */
func SendClusterReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultCluster{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Clusters = make([]somaproto.ProtoCluster, 0)
	for _, i := range (*r).Clusters {
		result.Clusters = append(result.Clusters, i.Cluster)
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
