package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// NodeList function
func NodeList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeReadHandler"].(*somaNodeReadHandler)
	handler.input <- somaNodeRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.NewNodeFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Node.Name != "" {
		filtered := []somaNodeResult{}
		for _, i := range result.Nodes {
			if i.Node.Name == cReq.Filter.Node.Name {
				filtered = append(filtered, i)
			}
		}
		result.Nodes = filtered
	}

skip:
	SendNodeReply(&w, &result)
}

// NodeShow function
func NodeShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node`,
		Action:     `show`,
		NodeID:     params.ByName(`node`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeReadHandler"].(*somaNodeReadHandler)
	handler.input <- somaNodeRequest{
		action: "show",
		reply:  returnChannel,
		Node: proto.Node{
			Id: params.ByName("node"),
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// NodeShowConfig function
func NodeShowConfig(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node`,
		Action:     `config`,
		NodeID:     params.ByName(`node`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeReadHandler"].(*somaNodeReadHandler)
	handler.input <- somaNodeRequest{
		action: "get_config",
		reply:  returnChannel,
		Node: proto.Node{
			Id: params.ByName("node"),
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// NodeSync function
func NodeSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node`,
		Action:     `sync`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeReadHandler"].(*somaNodeReadHandler)
	handler.input <- somaNodeRequest{
		action: `sync`,
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// NodeAdd function
func NodeAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewNodeRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	var serverID string
	if cReq.Node.ServerId != `` {
		serverID = cReq.Node.ServerId
	} else {
		serverID = `00000000-0000-0000-0000-000000000000`
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeWriteHandler"].(*somaNodeWriteHandler)
	handler.input <- somaNodeRequest{
		action: "add",
		reply:  returnChannel,
		user:   params.ByName(`AuthenticatedUser`),
		Node: proto.Node{
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			TeamId:    cReq.Node.TeamId,
			ServerId:  serverID,
			State:     "unassigned",
			IsOnline:  cReq.Node.IsOnline,
			IsDeleted: false,
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// NodeUpdate function
func NodeUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node`,
		Action:     `update`,
		NodeID:     params.ByName(`node`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewNodeRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeWriteHandler"].(*somaNodeWriteHandler)
	handler.input <- somaNodeRequest{
		action: `update`,
		reply:  returnChannel,
		Node: proto.Node{
			Id:        cReq.Node.Id,
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			TeamId:    cReq.Node.TeamId,
			ServerId:  cReq.Node.ServerId,
			IsOnline:  cReq.Node.IsOnline,
			IsDeleted: cReq.Node.IsDeleted,
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// NodeAssign function
func NodeAssign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if !IsAuthorized(&msg.Authorization{
		User:         params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `node`,
		Action:       `assign`,
		NodeID:       params.ByName(`node`),
		RepositoryID: cReq.Node.Config.RepositoryId,
		BucketID:     cReq.Node.Config.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "node",
		Action:      "assign_node",
		User:        params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		Node: somaNodeRequest{
			action: "assign",
			Node:   *cReq.Node,
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// NodeRemove function
func NodeRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	action := `remove`
	cReq := proto.NewNodeRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if cReq.Flags.Purge {
		action = `purge`
	}

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node`,
		Action:     action,
		NodeID:     params.ByName(`node`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeWriteHandler"].(*somaNodeWriteHandler)
	handler.input <- somaNodeRequest{
		action: action,
		reply:  returnChannel,
		Node: proto.Node{
			Id: params.ByName("node"),
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// NodeAddProperty function
func NodeAddProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	switch {
	case params.ByName("node") != cReq.Node.Id:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched node ids: %s, %s",
				params.ByName("node"),
				cReq.Node.Id))
		return
	case len(*cReq.Node.Properties) != 1:
		DispatchBadRequest(&w,
			fmt.Errorf("Expected property count 1, actual count: %d",
				len(*cReq.Node.Properties)))
		return
	case params.ByName("type") != (*cReq.Node.Properties)[0].Type:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched property types: %s, %s",
				params.ByName("type"),
				(*cReq.Node.Properties)[0].Type))
		return
	case (params.ByName("type") == "service") && (*cReq.Node.Properties)[0].Service.Name == "":
		DispatchBadRequest(&w,
			fmt.Errorf("Empty service name is invalid"))
		return
	}

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node`,
		Action:     `add_property`,
		NodeID:     params.ByName(`node`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "node",
		Action:      fmt.Sprintf("add_%s_property_to_node", params.ByName("type")),
		User:        params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		Node: somaNodeRequest{
			action: fmt.Sprintf("%s_property_new", params.ByName("type")),
			Node:   *cReq.Node,
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// NodeRemoveProperty function
func NodeRemoveProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	switch {
	case params.ByName(`node`) != cReq.Node.Id:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched node ids: %s, %s",
				params.ByName(`node`),
				cReq.Node.Id))
		return
	case cReq.Node.Config == nil:
		DispatchBadRequest(&w,
			fmt.Errorf(`Node configuration data missing`))
		return
	}
	// outside switch: _after_ nil test
	if cReq.Node.Config.RepositoryId == `` ||
		cReq.Node.Config.BucketId == `` {
		DispatchBadRequest(&w,
			fmt.Errorf(`Node configuration data incomplete`))
		return
	}

	if !IsAuthorized(&msg.Authorization{
		User:         params.ByName(`AuthenticatedUser`),
		RemoteAddr:   extractAddress(r.RemoteAddr),
		Section:      `node`,
		Action:       `remove_property`,
		NodeID:       params.ByName(`node`),
		RepositoryID: cReq.Node.Config.RepositoryId,
		BucketID:     cReq.Node.Config.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	node := proto.Node{
		Id: params.ByName(`node`),
		Config: &proto.NodeConfig{
			RepositoryId: cReq.Node.Config.RepositoryId,
			BucketId:     cReq.Node.Config.BucketId,
		},
		Properties: &[]proto.Property{
			proto.Property{
				Type:             params.ByName(`type`),
				RepositoryId:     cReq.Node.Config.RepositoryId,
				BucketId:         cReq.Node.Config.BucketId,
				SourceInstanceId: params.ByName(`source`),
			},
		},
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap[`guidePost`].(*guidePost)
	handler.input <- treeRequest{
		RequestType: `node`,
		Action: fmt.Sprintf("delete_%s_property_from_node",
			params.ByName(`type`)),
		User:  params.ByName(`AuthenticatedUser`),
		reply: returnChannel,
		Node: somaNodeRequest{
			action: fmt.Sprintf("%s_property_remove",
				params.ByName(`type`)),
			Node: node,
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

// SendNodeReply function
func SendNodeReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewNodeResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Nodes {
		*result.Nodes = append(*result.Nodes, i.Node)
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
