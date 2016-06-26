package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListNode(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeReadHandler"].(somaNodeReadHandler)
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
		filtered := make([]somaNodeResult, 0)
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

func ShowNode(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeReadHandler"].(somaNodeReadHandler)
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

func ShowNodeConfig(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeReadHandler"].(somaNodeReadHandler)
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

func SyncNode(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`node_sync`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeReadHandler"].(somaNodeReadHandler)
	handler.input <- somaNodeRequest{
		action: `sync`,
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

/* Write functions
 */
func AddNode(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`node_create`, ``, ``, ``); !ok {
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
	handler := handlerMap["nodeWriteHandler"].(somaNodeWriteHandler)
	handler.input <- somaNodeRequest{
		action: "add",
		reply:  returnChannel,
		// TODO: assign default server if no server information provided
		Node: proto.Node{
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			TeamId:    cReq.Node.TeamId,
			ServerId:  cReq.Node.ServerId,
			State:     "unassigned",
			IsOnline:  cReq.Node.IsOnline,
			IsDeleted: false,
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

func UpdateNode(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`node_update`, ``, ``, ``); !ok {
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
	handler := handlerMap["nodeWriteHandler"].(somaNodeWriteHandler)
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

func AssignNode(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
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

func DeleteNode(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`node_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	action := "delete"

	cReq := proto.NewNodeRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if cReq.Flags.Purge {
		action = "purge"
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeWriteHandler"].(somaNodeWriteHandler)
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

func AddPropertyToNode(w http.ResponseWriter, r *http.Request,
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

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
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

func DeletePropertyFromNode(w http.ResponseWriter, r *http.Request,
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
	handler := handlerMap[`guidePost`].(guidePost)
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

/* Utility
 */
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
