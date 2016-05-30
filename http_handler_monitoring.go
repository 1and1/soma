package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListMonitoring(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["monitoringReadHandler"].(somaMonitoringReadHandler)
	handler.input <- somaMonitoringRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.NewMonitoringFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Monitoring.Name != "" {
		filtered := make([]somaMonitoringResult, 0)
		for _, i := range result.Systems {
			if i.Monitoring.Name == cReq.Filter.Monitoring.Name {
				filtered = append(filtered, i)
			}
		}
		result.Systems = filtered
	}

skip:
	SendMonitoringReply(&w, &result)
}

func ShowMonitoring(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["monitoringReadHandler"].(somaMonitoringReadHandler)
	handler.input <- somaMonitoringRequest{
		action: "show",
		reply:  returnChannel,
		Monitoring: proto.Monitoring{
			Id: params.ByName("monitoring"),
		},
	}
	result := <-returnChannel
	SendMonitoringReply(&w, &result)
}

/* Write functions
 */
func AddMonitoring(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`monitoring_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewMonitoringRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["monitoringWriteHandler"].(somaMonitoringWriteHandler)
	handler.input <- somaMonitoringRequest{
		action: "add",
		reply:  returnChannel,
		Monitoring: proto.Monitoring{
			Name:     cReq.Monitoring.Name,
			Mode:     cReq.Monitoring.Mode,
			Contact:  cReq.Monitoring.Contact,
			TeamId:   cReq.Monitoring.TeamId,
			Callback: cReq.Monitoring.Callback,
		},
	}
	result := <-returnChannel
	SendMonitoringReply(&w, &result)
}

func DeleteMonitoring(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`monitoring_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["monitoringWriteHandler"].(somaMonitoringWriteHandler)
	handler.input <- somaMonitoringRequest{
		action: "delete",
		reply:  returnChannel,
		Monitoring: proto.Monitoring{
			Id: params.ByName("monitoring"),
		},
	}
	result := <-returnChannel
	SendMonitoringReply(&w, &result)
}

/* Utility
 */
func SendMonitoringReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewMonitoringResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Systems {
		*result.Monitorings = append(*result.Monitorings, i.Monitoring)
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
