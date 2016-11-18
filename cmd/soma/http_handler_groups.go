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

/*
 * Read functions
 */
func ListGroup(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `group`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["groupReadHandler"].(*somaGroupReadHandler)
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
			if i.Group.Name == cReq.Filter.Group.Name &&
				i.Group.BucketId == cReq.Filter.Group.BucketId {
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

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `group`,
		Action:     `show`,
		GroupID:    params.ByName(`group`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["groupReadHandler"].(*somaGroupReadHandler)
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

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `group`,
		Action:     `list_member`,
		GroupID:    params.ByName(`group`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["groupReadHandler"].(*somaGroupReadHandler)
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
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `group`,
		Action:     `create`,
		GroupID:    cReq.Group.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Group.Name)
	if nameLen < 4 || nameLen > 256 {
		DispatchBadRequest(&w, fmt.Errorf(`Illegal group name length (4 < x <= 256)`))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "group",
		Action:      "create_group",
		User:        params.ByName(`AuthenticatedUser`),
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

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `group`,
		Action:     `add_member`,
		GroupID:    cReq.Group.Id,
		BucketID:   cReq.Group.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	var rAct string
	switch {
	case len(*cReq.Group.MemberGroups) > 0:
		rAct = "add_group_to_group"
	case len(*cReq.Group.MemberClusters) > 0:
		rAct = "add_cluster_to_group"
	case len(*cReq.Group.MemberNodes) > 0:
		rAct = "add_node_to_group"
	}
	handler.input <- treeRequest{
		RequestType: "group",
		Action:      rAct,
		User:        params.ByName(`AuthenticatedUser`),
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

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `group`,
		Action:     `add_property`,
		GroupID:    params.ByName(`group`),
		BucketID:   cReq.Group.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["guidePost"].(*guidePost)
	handler.input <- treeRequest{
		RequestType: "group",
		Action:      fmt.Sprintf("add_%s_property_to_group", params.ByName("type")),
		User:        params.ByName(`AuthenticatedUser`),
		reply:       returnChannel,
		Group: somaGroupRequest{
			action: fmt.Sprintf("%s_property_new", params.ByName("type")),
			Group:  *cReq.Group,
		},
	}
	result := <-returnChannel
	SendGroupReply(&w, &result)
}

func DeletePropertyFromGroup(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	switch {
	case params.ByName(`group`) != cReq.Group.Id:
		DispatchBadRequest(&w,
			fmt.Errorf("Mismatched group ids: %s, %s",
				params.ByName(`group`),
				cReq.Group.Id))
		return
	case cReq.Group.BucketId == ``:
		DispatchBadRequest(&w,
			fmt.Errorf(`Missing bucketId in group delete request`))
		return
	}

	group := proto.Group{
		Id: params.ByName(`group`),
		Properties: &[]proto.Property{
			proto.Property{
				Type:             params.ByName(`type`),
				BucketId:         cReq.Group.BucketId,
				SourceInstanceId: params.ByName(`source`),
			},
		},
	}

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `group`,
		Action:     `remove_property`,
		GroupID:    params.ByName(`group`),
		BucketID:   cReq.Group.BucketId,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap[`guidePost`].(*guidePost)
	handler.input <- treeRequest{
		RequestType: `group`,
		Action: fmt.Sprintf("delete_%s_property_from_group",
			params.ByName(`type`)),
		User:  params.ByName(`AuthenticatedUser`),
		reply: returnChannel,
		Group: somaGroupRequest{
			action: fmt.Sprintf("%s_property_remove",
				params.ByName(`type`)),
			Group: group,
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
