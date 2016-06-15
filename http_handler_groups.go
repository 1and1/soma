package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

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
	cReq := proto.NewGroupFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Group.Name != "" {
		filtered := make([]somaGroupResult, 0)
		for _, i := range result.Groups {
			if i.Group.Name == cReq.Filter.Group.Name {
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
		Group: proto.Group{
			Id: params.ByName("group"),
		},
	}
	result := <-returnChannel
	SendGroupReply(&w, &result)
}

func ListGroupMembers(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["groupReadHandler"].(somaGroupReadHandler)
	handler.input <- somaGroupRequest{
		action: "member_list",
		reply:  returnChannel,
		Group: proto.Group{
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

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Group.Name)
	if nameLen < 4 || nameLen > 256 {
		DispatchBadRequest(&w, fmt.Errorf(`Illegal group name length (4 < x <= 256)`))
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

	cReq := proto.Request{}
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

func AddPropertyToGroup(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	switch {
	case params.ByName("group") != cReq.Group.Id:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched group ids: %s, %s",
				params.ByName("group"),
				cReq.Group.Id))
		return
	case len(*cReq.Group.Properties) != 1:
		DispatchBadRequest(&w,
			fmt.Errorf("Expected property count 1, actual count: %d",
				len(*cReq.Group.Properties)))
		return
	case params.ByName("type") != (*cReq.Group.Properties)[0].Type:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched property types: %s, %s",
				params.ByName("type"),
				(*cReq.Group.Properties)[0].Type))
		return
	case (params.ByName("type") == "service") && (*cReq.Group.Properties)[0].Service.Name == "":
		DispatchBadRequest(&w,
			fmt.Errorf("Empty service name is invalid"))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(guidePost)
	handler.input <- treeRequest{
		RequestType: "group",
		Action:      fmt.Sprintf("add_%s_property_to_group", params.ByName("type")),
		reply:       returnChannel,
		Group: somaGroupRequest{
			action: fmt.Sprintf("%s_property_new", params.ByName("type")),
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
	result := proto.NewGroupResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Groups {
		*result.Groups = append(*result.Groups, i.Group)
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
		}
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
