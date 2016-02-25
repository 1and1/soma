package main

import (
	"encoding/json"
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
	cReq := somaproto.ProtoRequestNode{}
	cReq.Filter = &somaproto.ProtoNodeFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaNodeResult, 0)
		for _, i := range result.Nodes {
			if i.Node.Name == cReq.Filter.Name {
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
		Node: somaproto.ProtoNode{
			Id: params.ByName("node"),
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

/* Write functions
 */
func AddNode(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestNode{}
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
		Node: somaproto.ProtoNode{
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			Team:      cReq.Node.Team,
			Server:    cReq.Node.Server,
			State:     "unassigned",
			IsOnline:  true,
			IsDeleted: false,
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

func AssignNode(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestNode{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
	handler.input <- treeRequest{
		RequestType: "node",
		Action:      "assign_node",
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
	action := "delete"

	cReq := somaproto.ProtoRequestNode{}
	_ = DecodeJsonBody(r, &cReq)
	if cReq.Purge {
		action = "purge"
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["nodeWriteHandler"].(somaNodeWriteHandler)
	handler.input <- somaNodeRequest{
		action: action,
		reply:  returnChannel,
		Node: somaproto.ProtoNode{
			Id: params.ByName("node"),
		},
	}
	result := <-returnChannel
	SendNodeReply(&w, &result)
}

/* Utility
 */
func SendNodeReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultNode{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Nodes = make([]somaproto.ProtoNode, 0)
	for _, i := range (*r).Nodes {
		result.Nodes = append(result.Nodes, i.Node)
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
