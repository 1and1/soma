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
	res := proto.NewEnvironmentResult()
	for _, env := range results {
		*res.Environments = append(*res.Environments, proto.Environment{Name: env.environment})
	}
	res.OK()
	json, err := json.Marshal(res)
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
	res := proto.NewEnvironmentResult()
	if len(results) == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if len(results) != 1 {
		http.Error(w, "Not found", http.StatusInternalServerError)
		return
	}
	*res.Environments = append(*res.Environments, proto.Environment{
		Name: results[0].environment,
	})
	res.OK()
	json, err := json.Marshal(res)
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
	var clientRequest proto.Request
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["environmentWriteHandler"].(somaEnvironmentWriteHandler)
	handler.input <- somaEnvironmentRequest{
		action:      "add",
		environment: clientRequest.Environment.Name,
		reply:       returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(proto.Result{
			StatusCode: 500,
			StatusText: "Internal Server Error",
			Errors:     &[]string{"Database statement returned no/wrong number of results"},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	result := results[0]
	if result.err != nil {
		json, _ := json.Marshal(proto.Result{
			StatusCode: 500,
			StatusText: "Internal Server Error",
			Errors:     &[]string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	json, _ := json.Marshal(proto.Result{
		StatusCode: 200,
		StatusText: "OK",
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
		DispatchInternalError(&w, fmt.Errorf("Database statement returned no/wrong number of results"))
		return
	}

	result := results[0]
	if result.err != nil {
		DispatchInternalError(&w, result.err)
		return
	}

	json, _ := json.Marshal(proto.Result{
		StatusCode: 200,
		StatusText: "OK",
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func RenameEnvironment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaEnvironmentResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest proto.Request
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["environmentWriteHandler"].(somaEnvironmentWriteHandler)
	handler.input <- somaEnvironmentRequest{
		action:      "rename",
		environment: params.ByName("environment"),
		rename:      clientRequest.Environment.Name,
		reply:       returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(proto.Result{
			StatusCode: 500,
			StatusText: "Internal Server Error",
			Errors:     &[]string{"Database statement returned no/wrong number of results"},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	result := results[0]
	if result.err != nil {
		json, _ := json.Marshal(proto.Result{
			StatusCode: 500,
			StatusText: "Internal Server Error",
			Errors:     &[]string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	json, _ := json.Marshal(proto.Result{
		StatusCode: 200,
		StatusText: "OK",
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
