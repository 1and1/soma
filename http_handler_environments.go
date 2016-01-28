package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 * Read functions
 */
func ListEnvironments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaEnvironmentResult)

	handler := handlerMap["environmentReadHandler"].(somaEnvironmentReadHandler)
	handler.input <- somaEnvironmentRequest{
		action: "list",
		reply:  returnChannel,
	}

	results := <-returnChannel
	environments := make([]string, len(results))
	for pos, res := range results {
		environments[pos] = res.environment
	}
	json, err := json.Marshal(somaproto.ProtoResultEnvironmentList{Code: 200, Status: "OK", Environments: environments})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ShowEnvironment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaEnvironmentResult)

	handler := handlerMap["environmentReadHandler"].(somaEnvironmentReadHandler)
	handler.input <- somaEnvironmentRequest{
		action:      "show",
		environment: params.ByName("environment"),
		reply:       returnChannel,
	}

	results := <-returnChannel
	if len(results) == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if len(results) != 1 {
		http.Error(w, "Not found", http.StatusInternalServerError)
		return
	}
	result := results[0]
	json, err := json.Marshal(somaproto.ProtoResultEnvironmentDetail{
		Code:    200,
		Status:  "OK",
		Details: somaproto.ProtoEnvironmentDetails{Environment: result.environment},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func AddEnvironment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaEnvironmentResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestEnvironment
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["environmentWriteHandler"].(somaEnvironmentWriteHandler)
	handler.input <- somaEnvironmentRequest{
		action:      "add",
		environment: clientRequest.Environment,
		reply:       returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{"Database statement returned no/wrong number of results"},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	result := results[0]
	if result.err != nil {
		json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Added environment: %s", result.environment)
	json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func DeleteEnvironment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaEnvironmentResult)

	handler := handlerMap["environmentWriteHandler"].(somaEnvironmentWriteHandler)
	handler.input <- somaEnvironmentRequest{
		action:      "delete",
		environment: params.ByName("environment"),
		reply:       returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{"Database statement returned no/wrong number of results"},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	result := results[0]
	if result.err != nil {
		json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Deleted environment: %s", result.environment)
	json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func RenameEnvironment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaEnvironmentResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestEnvironment
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["environmentWriteHandler"].(somaEnvironmentWriteHandler)
	handler.input <- somaEnvironmentRequest{
		action:      "rename",
		environment: params.ByName("environment"),
		rename:      clientRequest.Environment,
		reply:       returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{"Database statement returned no/wrong number of results"},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	result := results[0]
	if result.err != nil {
		json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Renamed environment: %s to %s", result.environment, clientRequest.Environment)
	json, _ := json.Marshal(somaproto.ProtoResultEnvironment{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
