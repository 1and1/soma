package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListTeam(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaTeamResult)
	handler := handlerMap["teamReadHandler"].(somaTeamReadHandler)
	handler.input <- somaTeamRequest{
		action: "list",
		reply:  returnChannel,
	}
	results := <-returnChannel

	cReq := somaproto.ProtoRequestTeam{}
	cReq.Filter = &somaproto.ProtoTeamFilter{}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaTeamResult, 0)
	filterloop:
		for _, iterTeam := range results {
			if iterTeam.rErr != nil {
				filtered = append(filtered, iterTeam)
				break filterloop
			}
			if iterTeam.team.Name == cReq.Filter.Name {
				filtered = append(filtered, iterTeam)
			}
		}
		results = filtered
	}

	SendTeamReply(&w, &results)
}

func ShowTeam(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaTeamResult)
	handler := handlerMap["teamReadHandler"].(somaTeamReadHandler)
	handler.input <- somaTeamRequest{
		action: "show",
		reply:  returnChannel,
		team: somaproto.ProtoTeam{
			Id: params.ByName("team"),
		},
	}
	results := <-returnChannel
	SendTeamReply(&w, &results)
}

/*
 * Write functions
 */
func AddTeam(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestTeam{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan []somaTeamResult)
	handler := handlerMap["teamWriteHandler"].(somaTeamWriteHandler)
	handler.input <- somaTeamRequest{
		action: "add",
		reply:  returnChannel,
		team: somaproto.ProtoTeam{
			Name:   cReq.Team.Name,
			Ldap:   cReq.Team.Ldap,
			System: cReq.Team.System,
		},
	}
	results := <-returnChannel
	SendTeamReply(&w, &results)
}

func DeleteTeam(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaTeamResult)
	handler := handlerMap["teamWriteHandler"].(somaTeamWriteHandler)
	handler.input <- somaTeamRequest{
		action: "delete",
		reply:  returnChannel,
		team: somaproto.ProtoTeam{
			Id: params.ByName("team"),
		},
	}
	results := <-returnChannel
	SendTeamReply(&w, &results)
}

/*
 * Utility
 */
func SendTeamReply(w *http.ResponseWriter, r *[]somaTeamResult) {
	var res somaproto.ProtoResultTeam
	dispatchError := CheckErrorHandler(r, &res)
	if dispatchError {
		goto dispatch
	}
	res.Text = make([]string, 0)
	res.Teams = make([]somaproto.ProtoTeam, 0)
	for _, l := range *r {
		res.Teams = append(res.Teams, l.team)
		if l.lErr != nil {
			res.Text = append(res.Text, l.lErr.Error())
		}
	}

dispatch:
	json, err := json.Marshal(res)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
