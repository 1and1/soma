package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// TeamList function
func TeamList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `team`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamReadHandler"].(*somaTeamReadHandler)
	handler.input <- somaTeamRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.NewTeamFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Team.Name != "" {
		filtered := []somaTeamResult{}
		for _, i := range result.Teams {
			if i.Team.Name == cReq.Filter.Team.Name {
				filtered = append(filtered, i)
			}
		}
		result.Teams = filtered
	}

skip:
	SendTeamReply(&w, &result)
}

// TeamShow function
func TeamShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `team`,
		Action:     `show`,
		TeamID:     params.ByName(`team`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamReadHandler"].(*somaTeamReadHandler)
	handler.input <- somaTeamRequest{
		action: "show",
		reply:  returnChannel,
		Team: proto.Team{
			Id: params.ByName("team"),
		},
	}
	result := <-returnChannel
	SendTeamReply(&w, &result)
}

// TeamSync function
func TeamSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `team`,
		Action:     `sync`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamReadHandler"].(*somaTeamReadHandler)
	handler.input <- somaTeamRequest{
		action: "sync",
		reply:  returnChannel,
	}
	result := <-returnChannel

	SendTeamReply(&w, &result)
}

// TeamAdd function
func TeamAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `team`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewTeamRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamWriteHandler"].(*somaTeamWriteHandler)
	handler.input <- somaTeamRequest{
		action: "add",
		reply:  returnChannel,
		Team: proto.Team{
			Name:     cReq.Team.Name,
			LdapId:   cReq.Team.LdapId,
			IsSystem: cReq.Team.IsSystem,
		},
	}
	result := <-returnChannel
	SendTeamReply(&w, &result)
}

// TeamUpdate function
func TeamUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `team`,
		Action:     `update`,
		TeamID:     params.ByName(`team`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewTeamRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamWriteHandler"].(*somaTeamWriteHandler)
	handler.input <- somaTeamRequest{
		action: `update`,
		reply:  returnChannel,
		Team: proto.Team{
			Id:       params.ByName(`team`),
			Name:     cReq.Team.Name,
			LdapId:   cReq.Team.LdapId,
			IsSystem: cReq.Team.IsSystem,
		},
	}
	result := <-returnChannel
	SendTeamReply(&w, &result)
}

// TeamRemove function
func TeamRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `team`,
		Action:     `remove`,
		TeamID:     params.ByName(`team`),
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamWriteHandler"].(*somaTeamWriteHandler)
	handler.input <- somaTeamRequest{
		action: "delete",
		reply:  returnChannel,
		Team: proto.Team{
			Id: params.ByName("team"),
		},
	}
	result := <-returnChannel
	SendTeamReply(&w, &result)
}

// SendTeamReply function
func SendTeamReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewTeamResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Teams {
		*result.Teams = append(*result.Teams, i.Team)
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
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
