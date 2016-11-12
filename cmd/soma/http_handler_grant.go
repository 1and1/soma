package main

import (
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

/* GLOBAL RIGHTS
 */
func GrantGlobalRight(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`grant_global_right`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	crq := proto.NewGrantRequest()
	err := DecodeJsonBody(r, &crq)
	// check body is consistent with URI
	if err != nil || crq.Grant.RecipientType != params.ByName(`rtyp`) ||
		crq.Grant.RecipientId != params.ByName(`rid`) ||
		crq.Grant.Category != `global` {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Section:    `right`,
		Action:     `grant`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Grant: proto.Grant{
			RecipientType: crq.Grant.RecipientType,
			RecipientId:   crq.Grant.RecipientId,
			PermissionId:  crq.Grant.PermissionId,
			Category:      crq.Grant.Category,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func RevokeGlobalRight(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`revoke_global_right`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Section:    `right`,
		Action:     `revoke`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			GrantId: params.ByName(`grant`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

/* LIMITED RIGHTS
 */
func GrantLimitedRight(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	scope := params.ByName(`scope`)
	obj := params.ByName(`uuid`)
	switch scope {
	case `repository`:
	default:
		// only implement repository for now
		DispatchNotImplemented(&w, nil)
		return
	}

	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`grant_limited_right`, obj, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	crq := proto.NewGrantRequest()
	err := DecodeJsonBody(r, &crq)
	// check body is consistent with URI
	if err != nil || crq.Grant.RecipientType != params.ByName(`rtyp`) ||
		crq.Grant.RecipientId != params.ByName(`rid`) ||
		crq.Grant.Category != `limited` {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Section:    `right`,
		Action:     `grant`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Grant: proto.Grant{
			RecipientType: crq.Grant.RecipientType,
			RecipientId:   crq.Grant.RecipientId,
			PermissionId:  crq.Grant.PermissionId,
			Category:      crq.Grant.Category,
			ObjectType:    crq.Grant.ObjectType,
			ObjectId:      crq.Grant.ObjectId,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func RevokeLimitedRight(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	scope := params.ByName(`scope`)
	obj := params.ByName(`uuid`)
	switch scope {
	case `repository`:
	default:
		// only implement repository for now
		DispatchNotImplemented(&w, nil)
		return
	}

	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`revoke_limited_right`, obj, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Section:    `right`,
		Action:     `revoke`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			GrantId: params.ByName(`grant`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

/* SYSTEM RIGHTS
 */
func GrantSystemRight(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`grant_system_right`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	crq := proto.NewGrantRequest()
	err := DecodeJsonBody(r, &crq)
	// check body is consistent with URI
	if err != nil || crq.Grant.RecipientType != params.ByName(`rtyp`) ||
		crq.Grant.RecipientId != params.ByName(`rid`) ||
		crq.Grant.Category != `system` {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Section:    `right`,
		Action:     `grant`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Grant: proto.Grant{
			RecipientType: crq.Grant.RecipientType,
			RecipientId:   crq.Grant.RecipientId,
			PermissionId:  crq.Grant.PermissionId,
			Category:      crq.Grant.Category,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func RevokeSystemRight(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`revoke_system_right`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Section:    `right`,
		Action:     `revoke`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			GrantId: params.ByName(`grant`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
