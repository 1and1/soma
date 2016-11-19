package main

import (
	"fmt"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// RightSearch function
func RightSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `right`,
		Action:     `search`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	crq := proto.NewGrantFilter()
	if err := DecodeJsonBody(r, &crq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	mr := msg.Request{
		Section:    `right`,
		Action:     `search`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Grant: proto.Grant{
			RecipientType: crq.Filter.Grant.RecipientType,
			RecipientId:   crq.Filter.Grant.RecipientId,
			PermissionId:  crq.Filter.Grant.PermissionId,
			Category:      crq.Filter.Grant.Category,
			ObjectType:    crq.Filter.Grant.ObjectType,
			ObjectId:      crq.Filter.Grant.ObjectId,
		},
	}

	handler.input <- mr
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// RightGrant function
func RightGrant(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	if cReq.Grant.Category != params.ByName(`category`) ||
		cReq.Grant.PermissionId != params.ByName(`permission`) {
		DispatchBadRequest(&w,
			fmt.Errorf(`Category/PermissionId mismatch`))
		return
	}

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `right`,
		Action:     `grant`,
		Grant:      cReq.Grant,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `right`,
		Action:     `grant`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Grant:      *cReq.Grant,
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// RightRevoke function
func RightRevoke(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	grant := proto.Grant{
		Id:           params.ByName(`grant`),
		Category:     params.ByName(`category`),
		PermissionId: params.ByName(`permission`),
	}

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `right`,
		Action:     `revoke`,
		Grant:      &grant,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section:    `right`,
		Action:     `revoke`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Grant:      grant,
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
