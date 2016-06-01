package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 * Read functions
 */
func ListObjectStates(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`states_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectStateResult)

	handler := handlerMap["objectStateReadHandler"].(somaObjectStateReadHandler)
	handler.input <- somaObjectStateRequest{
		action: "list",
		reply:  returnChannel,
	}

	results := <-returnChannel
	pres := proto.NewStateResult()
	for _, res := range results {
		*pres.States = append(*pres.States, proto.State{Name: res.state})
	}
	pres.OK()
	json, err := json.Marshal(pres)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ShowObjectState(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`states_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectStateResult)

	handler := handlerMap["objectStateReadHandler"].(somaObjectStateReadHandler)
	handler.input <- somaObjectStateRequest{
		action: "show",
		state:  params.ByName("state"),
		reply:  returnChannel,
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
	res := proto.NewStateResult()
	result := results[0]
	res.States = &[]proto.State{proto.State{Name: result.state}}
	res.OK()
	json, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func AddObjectState(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`states_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectStateResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest proto.Request
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["objectStateWriteHandler"].(somaObjectStateWriteHandler)
	handler.input <- somaObjectStateRequest{
		action: "add",
		state:  clientRequest.State.Name,
		reply:  returnChannel,
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

	if results[0].err != nil {
		json, _ := json.Marshal(proto.Result{
			StatusCode: 500,
			StatusText: "Internal Server Error",
			Errors:     &[]string{results[0].err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	json, _ := json.Marshal(proto.Result{
		StatusCode: 200,
		StatusText: "OK",
		States:     &[]proto.State{proto.State{Name: results[0].state}},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func DeleteObjectState(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`states_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectStateResult)

	handler := handlerMap["objectStateWriteHandler"].(somaObjectStateWriteHandler)
	handler.input <- somaObjectStateRequest{
		action: "delete",
		state:  params.ByName("state"),
		reply:  returnChannel,
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

func RenameObjectState(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`states_rename`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectStateResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest proto.Request
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["objectStateWriteHandler"].(somaObjectStateWriteHandler)
	handler.input <- somaObjectStateRequest{
		action: "rename",
		state:  params.ByName("state"),
		rename: clientRequest.State.Name,
		reply:  returnChannel,
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
