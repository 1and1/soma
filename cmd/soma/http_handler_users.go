package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// UserList function
func UserList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `user`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["userReadHandler"].(*somaUserReadHandler)
	handler.input <- somaUserRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.NewUserFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.User.UserName != "" {
		filtered := []somaUserResult{}
		for _, i := range result.Users {
			if i.User.UserName == cReq.Filter.User.UserName {
				filtered = append(filtered, i)
			}
		}
		result.Users = filtered
	}

skip:
	SendUserReply(&w, &result)
}

// UserShow function
func UserShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `user`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["userReadHandler"].(*somaUserReadHandler)
	handler.input <- somaUserRequest{
		action: "show",
		reply:  returnChannel,
		User: proto.User{
			Id: params.ByName("user"),
		},
	}
	result := <-returnChannel
	SendUserReply(&w, &result)
}

// UserSync function
func UserSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `user`,
		Action:     `sync`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["userReadHandler"].(*somaUserReadHandler)
	handler.input <- somaUserRequest{
		action: `sync`,
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendUserReply(&w, &result)
}

// UserAdd function
func UserAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `user`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewUserRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.User.UserName, `:`) {
		DispatchBadRequest(&w, fmt.Errorf(`Invalid username containing : character`))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["userWriteHandler"].(*somaUserWriteHandler)
	handler.input <- somaUserRequest{
		action: "add",
		reply:  returnChannel,
		User: proto.User{
			UserName:       cReq.User.UserName,
			FirstName:      cReq.User.FirstName,
			LastName:       cReq.User.LastName,
			EmployeeNumber: cReq.User.EmployeeNumber,
			MailAddress:    cReq.User.MailAddress,
			IsActive:       false,
			IsSystem:       cReq.User.IsSystem,
			IsDeleted:      false,
			TeamId:         cReq.User.TeamId,
		},
	}
	result := <-returnChannel
	SendUserReply(&w, &result)
}

// UserUpdate function
func UserUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `user`,
		Action:     `update`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewUserRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.User.UserName, `:`) {
		DispatchBadRequest(&w, fmt.Errorf(`Invalid username containing : character`))
		return
	}
	if params.ByName(`user`) != cReq.User.Id {
		DispatchBadRequest(&w, fmt.Errorf(`Mismatching user UUIDs in body and URL`))
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["userWriteHandler"].(*somaUserWriteHandler)
	handler.input <- somaUserRequest{
		action: "update",
		reply:  returnChannel,
		User: proto.User{
			Id:             cReq.User.Id,
			UserName:       cReq.User.UserName,
			FirstName:      cReq.User.FirstName,
			LastName:       cReq.User.LastName,
			EmployeeNumber: cReq.User.EmployeeNumber,
			MailAddress:    cReq.User.MailAddress,
			IsDeleted:      cReq.User.IsDeleted,
			TeamId:         cReq.User.TeamId,
		},
	}
	result := <-returnChannel
	SendUserReply(&w, &result)
}

// UserRemove function
func UserRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	cReq := proto.NewUserRequest()
	if err := DecodeJsonBody(r, &cReq); err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	action := `remove`
	if cReq.Flags.Purge {
		action = `purge`
	}

	if !IsAuthorized(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `user`,
		Action:     action,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["userWriteHandler"].(*somaUserWriteHandler)
	handler.input <- somaUserRequest{
		action: action,
		reply:  returnChannel,
		User: proto.User{
			Id: params.ByName("user"),
		},
	}
	result := <-returnChannel
	SendUserReply(&w, &result)
}

// SendUserReply function
func SendUserReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewUserResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Users {
		*result.Users = append(*result.Users, i.User)
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
