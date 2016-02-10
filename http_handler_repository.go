package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListRepository(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["repositoryReadHandler"].(somaRepositoryReadHandler)
	handler.input <- somaRepositoryRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestRepository{}
	cReq.Filter = &somaproto.ProtoRepositoryFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Name != "" {
		filtered := make([]somaRepositoryResult, 0)
		for _, i := range result.Repositories {
			if i.Repository.Name == cReq.Filter.Name {
				filtered = append(filtered, i)
			}
		}
		result.Repositories = filtered
	}

skip:
	SendRepositoryReply(&w, &result)
}

func ShowRepository(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["repositoryReadHandler"].(somaRepositoryReadHandler)
	handler.input <- somaRepositoryRequest{
		action: "show",
		reply:  returnChannel,
		Repository: somaproto.ProtoRepository{
			Id: params.ByName("repository"),
		},
	}
	result := <-returnChannel
	SendRepositoryReply(&w, &result)
}

/* Write functions
 */
func AddRepository(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestRepository{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["forestCustodian"].(forestCustodian)
	handler.input <- somaRepositoryRequest{
		action: "add",
		reply:  returnChannel,
		Repository: somaproto.ProtoRepository{
			Name:      cReq.Repository.Name,
			Team:      cReq.Repository.Team,
			IsDeleted: cReq.Repository.IsDeleted,
			IsActive:  cReq.Repository.IsActive,
		},
	}
	result := <-returnChannel
	SendRepositoryReply(&w, &result)
}

/*
 * Utility
 */
func SendRepositoryReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultRepository{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Repositories = make([]somaproto.ProtoRepository, 0)
	for _, i := range (*r).Repositories {
		result.Repositories = append(result.Repositories, i.Repository)
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
