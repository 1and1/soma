package main

import (
	"encoding/json"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
)

// MetricList function
func MetricList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `metric`,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["metricReadHandler"].(*somaMetricReadHandler)
	handler.input <- somaMetricRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendMetricReply(&w, &result)
}

// MetricShow function
func MetricShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `metric`,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["metricReadHandler"].(*somaMetricReadHandler)
	handler.input <- somaMetricRequest{
		action: "show",
		reply:  returnChannel,
		Metric: proto.Metric{
			Path: params.ByName("metric"),
		},
	}
	result := <-returnChannel
	SendMetricReply(&w, &result)
}

// MetricAdd function
func MetricAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `metric`,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.Request{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["metricWriteHandler"].(*somaMetricWriteHandler)
	handler.input <- somaMetricRequest{
		action: "add",
		reply:  returnChannel,
		Metric: *cReq.Metric,
	}
	result := <-returnChannel
	SendMetricReply(&w, &result)
}

// MetricRemove function
func MetricRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	if !IsAuthorizedd(&msg.Authorization{
		User:       params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    `metric`,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["metricWriteHandler"].(*somaMetricWriteHandler)
	handler.input <- somaMetricRequest{
		action: "delete",
		reply:  returnChannel,
		Metric: proto.Metric{
			Path: params.ByName("metric"),
		},
	}
	result := <-returnChannel
	SendMetricReply(&w, &result)
}

//SendMetricReply function
func SendMetricReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewMetricResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Metrics {
		*result.Metrics = append(*result.Metrics, i.Metric)
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
