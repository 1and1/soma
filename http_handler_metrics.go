package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListMetric(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`metrics_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["metricReadHandler"].(somaMetricReadHandler)
	handler.input <- somaMetricRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel
	SendMetricReply(&w, &result)
}

func ShowMetric(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`metrics_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["metricReadHandler"].(somaMetricReadHandler)
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

/* Write functions
 */
func AddMetric(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`metric_create`, ``, ``, ``); !ok {
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
	handler := handlerMap["metricWriteHandler"].(somaMetricWriteHandler)
	handler.input <- somaMetricRequest{
		action: "add",
		reply:  returnChannel,
		Metric: *cReq.Metric,
	}
	result := <-returnChannel
	SendMetricReply(&w, &result)
}

func DeleteMetric(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`metric_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["metricWriteHandler"].(somaMetricWriteHandler)
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

/* Utility
 */
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
