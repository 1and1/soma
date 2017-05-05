package main

import (
	"fmt"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// SystemOperation function
func SystemOperation(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.NewSystemOperationRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	// check the system operation is valid
	var sys *proto.SystemOperation
	switch cReq.SystemOperation.Request {
	case `repository_stop`:
		sys = &proto.SystemOperation{
			Request:      cReq.SystemOperation.Request,
			RepositoryId: cReq.SystemOperation.RepositoryId,
		}
	case `repository_rebuild`:
		sys = &proto.SystemOperation{
			Request:      cReq.SystemOperation.Request,
			RepositoryId: cReq.SystemOperation.RepositoryId,
			RebuildLevel: cReq.SystemOperation.RebuildLevel,
		}
	case `repository_restart`:
		sys = &proto.SystemOperation{
			Request:      cReq.SystemOperation.Request,
			RepositoryId: cReq.SystemOperation.RepositoryId,
		}
	case `shutdown`:
	default:
		DispatchBadRequest(&w, fmt.Errorf("%s %s",
			`Unknown system operation:`, cReq.SystemOperation.Request))
		return
	}

	// late authorization after Request check
	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `runtime`,
		Action:     cReq.SystemOperation.Request,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	switch cReq.SystemOperation.Request {
	case `repository_stop`:
		handler := handlerMap[`guidePost`].(*guidePost)
		handler.system <- msg.Request{
			Section:    `runtime`,
			Action:     cReq.SystemOperation.Request,
			Reply:      returnChannel,
			RemoteAddr: extractAddress(r.RemoteAddr),
			AuthUser:   params.ByName(`AuthenticatedUser`),
			System:     *sys,
		}
	case `repository_rebuild`, `repository_restart`:
		handler := handlerMap[`forestCustodian`].(*forestCustodian)
		handler.system <- msg.Request{
			Section:    `runtime`,
			Action:     cReq.SystemOperation.Request,
			Reply:      returnChannel,
			RemoteAddr: extractAddress(r.RemoteAddr),
			AuthUser:   params.ByName(`AuthenticatedUser`),
			System:     *sys,
		}
	case `shutdown`:
		handler := handlerMap[`grimReaper`].(*grimReaper)
		handler.system <- msg.Request{
			Section:    `runtime`,
			Action:     cReq.SystemOperation.Request,
			Reply:      returnChannel,
			RemoteAddr: extractAddress(r.RemoteAddr),
			AuthUser:   params.ByName(`AuthenticatedUser`),
		}
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
