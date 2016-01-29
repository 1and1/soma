package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListNode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaNodeResult)
	handler := handlerMap["nodeReadHandler"].(somaNodeReadHandler)
	handler.input <- somaNodeRequest{
		action: "list",
		reply:  returnChannel,
	}
	results := <-returnChannel

	cReq := somaproto.ProtoRequestNode{}
	cReq.Filter = &somaproto.ProtoNodeFilter{}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaNodeResult, 0)
	filterloop:
		for _, iterNode := range results {
			if iterNode.rErr != nil {
				filtered = append(filtered, iterNode)
				break filterloop
			}
			if iterNode.node.Name == cReq.Filter.Name {
				filtered = append(filtered, iterNode)
			}
		}
		results = filtered
	}

	SendNodeReply(&w, &results)
}

func ShowNode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaNodeResult)
	handler := handlerMap["nodeReadHandler"].(somaNodeReadHandler)
	handler.input <- somaNodeRequest{
		action: "show",
		reply:  returnChannel,
		node: somaproto.ProtoNode{
			Id: params.ByName("node"),
		},
	}
	results := <-returnChannel
	SendNodeReply(&w, &results)
}

/* Write functions
 */
func AddNode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestNode{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan []somaNodeResult)
	handler := handlerMap["nodeWriteHandler"].(somaNodeWriteHandler)
	handler.input <- somaNodeRequest{
		action: "add",
		reply:  returnChannel,
		// TODO: assign default server if no server information provided
		node: somaproto.ProtoNode{
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			Team:      cReq.Node.Team,
			Server:    cReq.Node.Server,
			State:     "standalone",
			IsOnline:  true,
			IsDeleted: false,
		},
	}
	results := <-returnChannel
	SendNodeReply(&w, &results)
}

/* Utility
 */
func SendNodeReply(w *http.ResponseWriter, r *[]somaNodeResult) {
	res := somaproto.ProtoResultNode{}
	dispatchError := CheckErrorHandler(r, &res)
	if dispatchError {
		goto dispatch
	}
	res.Text = make([]string, 0)
	res.Nodes = make([]somaproto.ProtoNode, 0)
	for _, l := range *r {
		res.Nodes = append(res.Nodes, l.node)
		if l.lErr != nil {
			res.Text = append(res.Text, l.lErr.Error())
		}
	}

dispatch:
	json, err := json.Marshal(res)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
