package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// ClusterList function
func ClusterList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `cluster`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["clusterReadHandler"].(*somaClusterReadHandler)
	handler.input <- somaClusterRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.Request{}
	cReq.Filter = &proto.Filter{}
	cReq.Filter.Cluster = &proto.ClusterFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Cluster.Name != "" {
		filtered := []somaClusterResult{}
		for _, i := range result.Clusters {
			if i.Cluster.Name == cReq.Filter.Cluster.Name &&
				i.Cluster.BucketId == cReq.Filter.Cluster.BucketId {
				filtered = append(filtered, i)
			}
		}
		result.Clusters = filtered
	}

skip:
	SendClusterReply(&w, &result)
}

// ClusterShow function
func ClusterShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `cluster`,
		Action:     `show`,
		ClusterID:  params.ByName(`cluster`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["clusterReadHandler"].(*somaClusterReadHandler)
	handler.input <- somaClusterRequest{
		action: "show",
		reply:  returnChannel,
		Cluster: proto.Cluster{
			Id: params.ByName("cluster"),
		},
	}
	result := <-returnChannel
	SendClusterReply(&w, &result)
}

// ClusterListMember function
func ClusterListMember(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `cluster`,
		Action:     `list_member`,
		ClusterID:  params.ByName(`cluster`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["clusterReadHandler"].(*somaClusterReadHandler)
	handler.input <- somaClusterRequest{
		action: "member_list",
		reply:  returnChannel,
		Cluster: proto.Cluster{
			Id: params.ByName("cluster"),
		},
	}
	result := <-returnChannel
	SendClusterReply(&w, &result)
}

// ClusterCreate function
func ClusterCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `cluster`,
		Action:     `create`,
		BucketID:   cReq.Cluster.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Cluster.Name)
	if nameLen < 4 || nameLen > 256 {
		DispatchBadRequest(&w, fmt.Errorf(`Illegal cluster name length (4 < x <= 256)`))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "cluster",
		Action:      "create_cluster",
		AuthUser:    params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		Cluster: somaClusterRequest{
			action:  "add",
			Cluster: *cReq.Cluster,
		},
	}
	result := <-returnChannel
	SendClusterReply(&w, &result)
}

// ClusterAddMember function
func ClusterAddMember(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `cluster`,
		Action:     `add_member`,
		ClusterID:  cReq.Cluster.Id,
		BucketID:   cReq.Cluster.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "cluster",
		Action:      "add_node_to_cluster",
		AuthUser:    params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		Cluster: somaClusterRequest{
			action:  "member",
			Cluster: *cReq.Cluster,
		},
	}
	result := <-returnChannel
	SendClusterReply(&w, &result)
}

// ClusterAddProperty function
func ClusterAddProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	switch {
	case params.ByName("cluster") != cReq.Cluster.Id:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched cluster ids: %s, %s",
				params.ByName("cluster"),
				cReq.Cluster.Id))
		return
	case len(*cReq.Cluster.Properties) != 1:
		DispatchBadRequest(&w,
			fmt.Errorf("Expected property count 1, actual count: %d",
				len(*cReq.Cluster.Properties)))
		return
	case params.ByName("type") != (*cReq.Cluster.Properties)[0].Type:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched property types: %s, %s",
				params.ByName("type"),
				(*cReq.Cluster.Properties)[0].Type))
		return
	case (params.ByName("type") == "service") && (*cReq.Cluster.Properties)[0].Service.Name == "":
		DispatchBadRequest(&w,
			fmt.Errorf("Empty service name is invalid"))
		return
	}

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `cluster`,
		Action:     `add_property`,
		ClusterID:  params.ByName(`cluster`),
		BucketID:   cReq.Cluster.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "cluster",
		Action:      fmt.Sprintf("add_%s_property_to_cluster", params.ByName("type")),
		AuthUser:    params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		Cluster: somaClusterRequest{
			action:  fmt.Sprintf("%s_property_new", params.ByName("type")),
			Cluster: *cReq.Cluster,
		},
	}
	result := <-returnChannel
	SendClusterReply(&w, &result)
}

// ClusterRemoveProperty function
func ClusterRemoveProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	switch {
	case params.ByName(`cluster`) != cReq.Cluster.Id:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched cluster ids: %s, %s",
				params.ByName(`cluster`),
				cReq.Cluster.Id))
		return
	case cReq.Cluster.BucketId == ``:
		DispatchBadRequest(&w,
			fmt.Errorf(`Missing bucketId in bucket property delete request`))
		return
	}

	cluster := proto.Cluster{
		Id: params.ByName(`cluster`),
		Properties: &[]proto.Property{
			proto.Property{
				Type:             params.ByName(`type`),
				BucketId:         cReq.Cluster.BucketId,
				SourceInstanceId: params.ByName(`source`),
			},
		},
	}

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `cluster`,
		Action:     `remove_property`,
		ClusterID:  params.ByName(`cluster`),
		BucketID:   cReq.Cluster.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap[`guidePost`].(*guidePost)
	handler.input <- treeRequest{
		RequestType: `cluster`,
		Action: fmt.Sprintf("delete_%s_property_from_cluster",
			params.ByName(`type`)),
		AuthUser: params.ByName(`AuthenticatedUser`),
		reply:    returnChannel,
		Cluster: somaClusterRequest{
			action: fmt.Sprintf("%s_property_remove",
				params.ByName(`type`)),
			Cluster: cluster,
		},
	}
	result := <-returnChannel
	SendClusterReply(&w, &result)
}

// SendClusterReply function
func SendClusterReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.Result{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	if result.Errors == nil {
		result.Errors = &[]string{}
	}
	result.Clusters = &[]proto.Cluster{}
	for _, i := range (*r).Clusters {
		*result.Clusters = append(*result.Clusters, i.Cluster)
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
