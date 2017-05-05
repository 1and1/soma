/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest

import (
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/soma"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// NodeAdd function
func (x *Rest) NodeAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	if !x.isAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node-mgmt`,
		Action:     `add`,
	}) {
		dispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	var serverID string
	if cReq.Node.ServerId != `` {
		serverID = cReq.Node.ServerId
	} else {
		serverID = `00000000-0000-0000-0000-000000000000`
	}

	returnChannel := make(chan msg.Result)
	handler := x.handlerMap.Get(`node_w`).(*soma.NodeWrite)
	handler.Input <- msg.Request{
		Section:    `node-mgmt`,
		Action:     `add`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Node: proto.Node{
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			TeamId:    cReq.Node.TeamId,
			ServerId:  serverID,
			State:     `unassigned`,
			IsOnline:  cReq.Node.IsOnline,
			IsDeleted: false,
		},
	}
	result := <-returnChannel
	sendMsgResult(&w, &result)
}

// NodeSync function
func (x *Rest) NodeSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	if !x.isAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node-mgmt`,
		Action:     `sync`,
	}) {
		dispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := x.handlerMap.Get(`node_r`).(*soma.NodeRead)
	handler.Input <- msg.Request{
		Section:    `node-mgmt`,
		Action:     `sync`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
	}
	result := <-returnChannel
	sendMsgResult(&w, &result)
}

// NodeUpdate function
func (x *Rest) NodeUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	if !x.isAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node-mgmt`,
		Action:     `update`,
		NodeID:     params.ByName(`nodeID`),
	}) {
		dispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewNodeRequest()
	err := decodeJSONBody(r, &cReq)
	if err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := x.handlerMap.Get(`node_w`).(*soma.NodeWrite)
	handler.Input <- msg.Request{
		Section:    `node-mgmt`,
		Action:     `update`,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Node: proto.Node{
			Id:        cReq.Node.Id,
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			TeamId:    cReq.Node.TeamId,
			ServerId:  cReq.Node.ServerId,
			IsOnline:  cReq.Node.IsOnline,
			IsDeleted: cReq.Node.IsDeleted,
		},
	}
	result := <-returnChannel
	sendMsgResult(&w, &result)
}

// NodeRemove function
func (x *Rest) NodeRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	action := `remove`
	cReq := proto.NewNodeRequest()
	err := decodeJSONBody(r, &cReq)
	if err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if cReq.Flags.Purge {
		action = `purge`
	}

	if !x.isAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `node-mgmt`,
		Action:     action,
		NodeID:     params.ByName(`nodeID`),
	}) {
		dispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := x.handlerMap.Get(`node_w`).(*soma.NodeWrite)
	handler.Input <- msg.Request{
		Section:    `node-mgmt`,
		Action:     action,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Node: proto.Node{
			Id: params.ByName(`nodeID`),
		},
	}
	result := <-returnChannel
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
