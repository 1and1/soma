package main

import (
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListPermission(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`permission_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Action:     `permission`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `list`,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func ShowPermission(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`permission_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Action:     `permission`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `show`,
		},
		Permission: proto.Permission{
			Name: params.ByName(`permission`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func SearchPermission(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`permission_search`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	crq := proto.NewPermissionFilter()
	_ = DecodeJsonBody(r, &crq)
	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	mr := msg.Request{
		Type:       `supervisor`,
		Action:     `permission`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `search/name`,
		},
		Permission: proto.Permission{
			Name: crq.Filter.Permission.Name,
		},
	}

	handler.input <- mr
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

/* Write functions
 */

func AddPermission(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`permission_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewPermissionRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Action:     `permission`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `add`,
		},
		Permission: proto.Permission{
			Name:     cReq.Permission.Name,
			Category: cReq.Permission.Category,
			Grants:   cReq.Permission.Grants,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func DeletePermission(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`permission_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:       `supervisor`,
		Action:     `permission`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		Super: &msg.Supervisor{
			Action: `delete`,
		},
		Permission: proto.Permission{
			Name: params.ByName(`permission`),
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
