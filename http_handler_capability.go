package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListCapability(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["capabilityReadHandler"].(somaCapabilityReadHandler)
	handler.input <- somaCapabilityRequest{
		action: "list",
		reply:  returnChannel,
	}
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := somaproto.ProtoRequestCapability{}
	cReq.Filter = &somaproto.ProtoCapabilityFilter{}
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if cReq.Filter.Monitoring != "" {
		filtered := make([]somaCapabilityResult, 0)
		for _, i := range result.Capabilities {
			if i.Capability.Monitoring == cReq.Filter.Monitoring &&
				i.Capability.Metric == cReq.Filter.Metric &&
				i.Capability.View == cReq.Filter.View {
				filtered = append(filtered, i)
			}
		}
		result.Capabilities = filtered
	}

	// cleanup reply
	for i, _ := range result.Capabilities {
		result.Capabilities[i].Capability.Monitoring = ""
		result.Capabilities[i].Capability.Metric = ""
		result.Capabilities[i].Capability.View = ""
	}

skip:
	SendCapabilityReply(&w, &result)
}

func ShowCapability(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["capabilityReadHandler"].(somaCapabilityReadHandler)
	handler.input <- somaCapabilityRequest{
		action: "show",
		reply:  returnChannel,
		Capability: somaproto.ProtoCapability{
			Id: params.ByName("capability"),
		},
	}
	result := <-returnChannel
	SendCapabilityReply(&w, &result)
}

/* Write functions
 */
func AddCapability(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	cReq := somaproto.ProtoRequestCapability{}
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan somaResult)
	handler := handlerMap["capabilityWriteHandler"].(somaCapabilityWriteHandler)
	handler.input <- somaCapabilityRequest{
		action: "add",
		reply:  returnChannel,
		Capability: somaproto.ProtoCapability{
			Monitoring: cReq.Capability.Monitoring,
			Metric:     cReq.Capability.Metric,
			View:       cReq.Capability.View,
			Thresholds: cReq.Capability.Thresholds,
		},
	}
	result := <-returnChannel
	SendCapabilityReply(&w, &result)
}

func DeleteCapability(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	handler := handlerMap["capabilityWriteHandler"].(somaCapabilityWriteHandler)
	handler.input <- somaCapabilityRequest{
		action: "delete",
		reply:  returnChannel,
		Capability: somaproto.ProtoCapability{
			Id: params.ByName("capability"),
		},
	}
	result := <-returnChannel
	SendCapabilityReply(&w, &result)
}

/* Utility
 */
func SendCapabilityReply(w *http.ResponseWriter, r *somaResult) {
	result := somaproto.ProtoResultCapability{}
	if r.MarkErrors(&result) {
		goto dispatch
	}
	result.Text = make([]string, 0)
	result.Capabilities = make([]somaproto.ProtoCapability, 0)
	for _, i := range (*r).Capabilities {
		result.Capabilities = append(result.Capabilities, i.Capability)
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
