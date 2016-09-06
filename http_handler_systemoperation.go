package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func SystemOperation(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`system_operation`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewSystemOperationRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	// check the system operation is valid
	var sys *proto.SystemOperation
	switch cReq.SystemOperation.Request {
	case `stop_repository`:
		sys = &proto.SystemOperation{
			Request:      cReq.SystemOperation.Request,
			RepositoryId: cReq.SystemOperation.RepositoryId,
		}
	default:
		DispatchBadRequest(&w, fmt.Errorf("%s %s",
			`Unknown system operation:`, cReq.SystemOperation.Request))
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`guidePost`].(guidePost)
	handler.system <- msg.Request{
		Type:       `guidepost`,
		Action:     `systemoperation`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		User:       params.ByName(`AuthenticatedUser`),
		System:     *sys,
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
