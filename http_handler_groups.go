package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListGroup(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["groupReadHandler"].(somaGroupReadHandler)
	handler.input <- somaGroupRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestGroup{}
	cReq.Filter = &somaproto.ProtoGroupFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaGroupResult, 0)
		for _, i := range result.Groups {
			if i.Group.Name == cReq.Filter.Name {
				filtered = append(filtered, i)
			}
		}
		result.Groups = filtered
	}

skip:
	SendGroupReply(&w, &result)
}

func ShowGroup(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["groupReadHandler"].(somaGroupReadHandler)
	handler.input <- somaGroupRequest{
		action: "show",
		reply:  returnChannel,
		Group: somaproto.ProtoGroup{
			Id: params.ByName("group"),
		},
	}
	result := <-returnChannel
	SendGroupReply(&w, &result)
}

/* Write functions
 */
func AddGroup(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestGroup{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
	handler.input <- treeRequest{
		RequestType: "group",
		Action:      "create_group",
		reply:       returnChannel,
		Group: somaGroupRequest{
			action: "add",
			Group:  *cReq.Group,
		},
	}
	result := <-returnChannel
	SendGroupReply(&w, &result)
}

func AddMemberToGroup(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestGroup{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
	var rAct string
	switch {
	case len(cReq.Group.MemberGroups) > 0:
		rAct = "add_group_to_group"
	case len(cReq.Group.MemberClusters) > 0:
		rAct = "add_cluster_to_group"
	case len(cReq.Group.MemberNodes) > 0:
		rAct = "add_node_to_group"
	}
	handler.input <- treeRequest{
		RequestType: "group",
		Action:      rAct,
		reply:       returnChannel,
		Group: somaGroupRequest{
			action: "member",
			Group:  *cReq.Group,
		},
	}
	result := <-returnChannel
	SendGroupReply(&w, &result)
}

/*
 * Utility
 */
func SendGroupReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultGroup{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Groups = make([]somaproto.ProtoGroup, 0)
	for _, i := range (*r).Groups {
		result.Groups = append(result.Groups, i.Group)
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
