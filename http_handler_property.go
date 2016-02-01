package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListProperty(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "list",
		reply:  returnChannel,
	}
	// strip surrounding / and skip first path element `property`
	el := strings.Split(strings.Trim(r.URL.Path, "/"), "/")[1:]
	switch el[0] {
	case "native":
		req.prType = "native"
	case "system":
		req.prType = "system"
	case "custom":
		req.prType = "custom"
		req.Custom.Repository = params.ByName("repository")
	case "service":
		switch el[1] {
		case "team":
			req.prType = "service"
			req.Service.Team = params.ByName("team")
		case "global":
			req.prType = "template"
		default:
			SendPropertyReply(&w, &somaResult{})
			return
		}
	default:
		SendPropertyReply(&w, &somaResult{})
		return
	}

	handler := handlerMap["propertyReadHandler"].(somaPropertyReadHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

func ShowMetric(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "show",
		reply:  returnChannel,
	}
	el := strings.Split(strings.Trim(r.URL.Path, "/"), "/")[1:]
	switch el[0] {
	case "native":
		req.prType = "native"
		req.Native.Property = params.ByName("native")
	case "system":
		req.prType = "system"
		req.System.Property = params.ByName("system")
	case "custom":
		req.prType = "custom"
		req.Custom.Id = params.ByName("custom")
		req.Custom.Repository = params.ByName("repository")
	case "service":
		switch el[1] {
		case "team":
			req.prType = "service"
			req.Service.Property = params.ByName("service")
			req.Service.Team = params.ByName("team")
		case "global":
			req.prType = "template"
			req.Service.Property = params.ByName("service")
		default:
			SendPropertyReply(&w, &somaResult{})
			return
		}
	default:
		SendPropertyReply(&w, &somaResult{})
		return
	}

	handler := handlerMap["propertyReadHandler"].(somaPropertyReadHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

/* Write functions
 */
func AddMetric(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestMetric{}
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

	returnChannel := make(chan somaResult)
	handler := handlerMap["metricWriteHandler"].(somaMetricWriteHandler)
	handler.input <- somaMetricRequest{
		action: "delete",
		reply:  returnChannel,
		Metric: somaproto.ProtoMetric{
			Metric: params.ByName("metric"),
		},
	}
	result := <-returnChannel
	SendMetricReply(&w, &result)
}

/* Utility
 */
func SendMetricReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultMetric{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Metrics = make([]somaproto.ProtoMetric, 0)
	for _, i := range (*r).Metrics {
		result.Metrics = append(result.Metrics, i.Metric)
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
