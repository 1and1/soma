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
func ListEnvironments(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`environments_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaEnvironmentResult)

	handler := handlerMap["environmentReadHandler"].(somaEnvironmentReadHandler)
	handler.input <- somaEnvironmentRequest{
		action: "list",
		reply:  returnChannel,
	}

	results := <-returnChannel
	res := proto.NewEnvironmentResult()
	for _, env := range results {
		if res.Error(env.err) {
			goto dispatch
		}
		*res.Environments = append(*res.Environments, proto.Environment{Name: env.environment})
	}
	if len(*res.Environments) == 0 {
		res.NotFound()
	} else {
		res.OK()
	}

dispatch:
	jsonB, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonB)
}

func ShowEnvironment(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`environments_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
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
		res.NotFound()
		goto dispatch
	}
	if len(results) != 1 {
		res.NotFound()
		goto dispatch
	}
	if results[0].environment == `` {
		res.NotFound()
		goto dispatch
	}
	*res.Environments = append(*res.Environments, proto.Environment{
		Name: results[0].environment,
	})
	res.OK()

dispatch:
	jsonB, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonB)
}

func AddEnvironment(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`environments_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaEnvironmentResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	clientRequest := proto.NewEnvironmentRequest()
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
		jsonB, _ := json.Marshal(&proto.Result{
			StatusCode:   500,
			StatusText:   "Internal Server Error",
			Errors:       &[]string{result.err.Error()},
			Environments: &[]proto.Environment{},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonB)
		return
	}

	jsonB, _ := json.Marshal(&proto.Result{
		StatusCode:   200,
		StatusText:   "OK",
		Errors:       &[]string{},
		Environments: &[]proto.Environment{proto.Environment{Name: clientRequest.Environment.Name}},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonB)
}

func DeleteEnvironment(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`environments_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
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

	jsonB, _ := json.Marshal(proto.Result{
		StatusCode: 200,
		StatusText: "OK",
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonB)
}

func RenameEnvironment(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`environments_rename`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
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
		jsonB, _ := json.Marshal(proto.Result{
			StatusCode: 500,
			StatusText: "Internal Server Error",
			Errors:     &[]string{"Database statement returned no/wrong number of results"},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonB)
		return
	}

	result := results[0]
	if result.err != nil {
		jsonB, _ := json.Marshal(proto.Result{
			StatusCode: 500,
			StatusText: "Internal Server Error",
			Errors:     &[]string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonB)
		return
	}

	jsonB, _ := json.Marshal(proto.Result{
		StatusCode: 200,
		StatusText: "OK",
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonB)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
