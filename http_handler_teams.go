package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListTeam(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamReadHandler"].(somaTeamReadHandler)
	handler.input <- somaTeamRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestTeam{}
	cReq.Filter = &somaproto.ProtoTeamFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaTeamResult, 0)
		for _, i := range result.Teams {
			if i.Team.Name == cReq.Filter.Name {
				filtered = append(filtered, i)
			}
		}
		result.Teams = filtered
	}

skip:
	SendTeamReply(&w, &result)
}

func ShowTeam(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamReadHandler"].(somaTeamReadHandler)
	handler.input <- somaTeamRequest{
		action: "show",
		reply:  returnChannel,
		Team: somaproto.ProtoTeam{
			Id: params.ByName("team"),
		},
	}
	result := <-returnChannel
	SendTeamReply(&w, &result)
}

/* Write functions
 */
func AddTeam(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestTeam{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamWriteHandler"].(somaTeamWriteHandler)
	handler.input <- somaTeamRequest{
		action: "add",
		reply:  returnChannel,
		Team: somaproto.ProtoTeam{
			Name:   cReq.Team.Name,
			Ldap:   cReq.Team.Ldap,
			System: cReq.Team.System,
		},
	}
	result := <-returnChannel
	SendTeamReply(&w, &result)
}

func DeleteTeam(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["teamWriteHandler"].(somaTeamWriteHandler)
	handler.input <- somaTeamRequest{
		action: "delete",
		reply:  returnChannel,
		Team: somaproto.ProtoTeam{
			Id: params.ByName("team"),
		},
	}
	result := <-returnChannel
	SendTeamReply(&w, &result)
}

/*
 * Utility
 */
func SendTeamReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultTeam{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Teams = make([]somaproto.ProtoTeam, 0)
	for _, i := range (*r).Teams {
		result.Teams = append(result.Teams, i.Team)
		if i.ResultError != nil {
			result.Text = append(result.Text, i.ResultError.Error())
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