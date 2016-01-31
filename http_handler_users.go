package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListUser(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["userReadHandler"].(somaUserReadHandler)
	handler.input <- somaUserRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestUser{}
	cReq.Filter = &somaproto.ProtoUserFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.UserName != "" {
		filtered := make([]somaUserResult, 0)
		for _, i := range result.Users {
			if i.User.UserName == cReq.Filter.UserName {
				filtered = append(filtered, i)
			}
		}
		result.Users = filtered
	}

skip:
	SendUserReply(&w, &result)
}

func ShowUser(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["userReadHandler"].(somaUserReadHandler)
	handler.input <- somaUserRequest{
		action: "show",
		reply:  returnChannel,
		User: somaproto.ProtoUser{
			Id: params.ByName("user"),
		},
	}
	result := <-returnChannel
	SendUserReply(&w, &result)
}

/* Write functions
 */
func AddUser(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestUser{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["userWriteHandler"].(somaUserWriteHandler)
	handler.input <- somaUserRequest{
		action: "add",
		reply:  returnChannel,
		User: somaproto.ProtoUser{
			UserName:       cReq.User.UserName,
			FirstName:      cReq.User.FirstName,
			LastName:       cReq.User.LastName,
			EmployeeNumber: cReq.User.EmployeeNumber,
			MailAddress:    cReq.User.MailAddress,
			IsActive:       false,
			IsSystem:       cReq.User.IsSystem,
			IsDeleted:      false,
			Team:           cReq.User.Team,
		},
	}
	result := <-returnChannel
	SendUserReply(&w, &result)
}

func DeleteUser(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	action := "delete"

	cReq := somaproto.ProtoRequestUser{}
	_ = DecodeJsonBody(r, &cReq)
	if cReq.Purge {
		action = "purge"
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["userWriteHandler"].(somaUserWriteHandler)
	handler.input <- somaUserRequest{
		action: action,
		reply:  returnChannel,
		User: somaproto.ProtoUser{
			Id: params.ByName("user"),
		},
	}
	result := <-returnChannel
	SendUserReply(&w, &result)
}

/* Utility
 */
func SendUserReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultUser{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Users = make([]somaproto.ProtoUser, 0)
	for _, i := range (*r).Users {
		result.Users = append(result.Users, i.User)
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
