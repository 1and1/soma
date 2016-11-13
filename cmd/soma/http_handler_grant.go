package main

import (
	"fmt"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListRights(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
}

func ListUserRights(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
}

func ListUserRepoRights(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
}

func SearchGrant(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`grant_search`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	crq := proto.NewGrantFilter()
	_ = DecodeJsonBody(r, &crq)
	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	mr := msg.Request{
		Type:       `supervisor`,
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
		},
	}

	handler.input <- mr
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

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
		Type:       `supervisor`,
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
		Type:       `supervisor`,
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
