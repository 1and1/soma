package main

import (
	"encoding/json"
	"github.com/1and1/soma/lib/proto"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 * Read functions
 */
func ListObjectTypes(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`types_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectTypeResult)

	handler := handlerMap["objectTypeReadHandler"].(*somaObjectTypeReadHandler)
	handler.input <- somaObjectTypeRequest{
		action: "list",
		reply:  returnChannel,
	}

	results := <-returnChannel
	pres := proto.NewEntityResult()
	for _, res := range results {
		*pres.Entities = append(*pres.Entities, proto.Entity{Name: res.objectType})
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

func ShowObjectType(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`types_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectTypeResult)

	handler := handlerMap["objectTypeReadHandler"].(*somaObjectTypeReadHandler)
	handler.input <- somaObjectTypeRequest{
		action:     "show",
		objectType: params.ByName("objectType"),
		reply:      returnChannel,
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
	json, err := json.Marshal(proto.Result{
		StatusCode: 200,
		StatusText: "OK",
		Entities:   &[]proto.Entity{proto.Entity{Name: result.objectType}},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func AddObjectType(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`types_create`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectTypeResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest proto.Request
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["objectTypeWriteHandler"].(*somaObjectTypeWriteHandler)
	handler.input <- somaObjectTypeRequest{
		action:     "add",
		objectType: clientRequest.Entity.Name,
		reply:      returnChannel,
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

func DeleteObjectType(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`types_delete`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectTypeResult)

	handler := handlerMap["objectTypeWriteHandler"].(*somaObjectTypeWriteHandler)
	handler.input <- somaObjectTypeRequest{
		action:     "delete",
		objectType: params.ByName("objectType"),
		reply:      returnChannel,
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

func RenameObjectType(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`types_rename`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
	returnChannel := make(chan []somaObjectTypeResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest proto.Request
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["objectTypeWriteHandler"].(*somaObjectTypeWriteHandler)
	handler.input <- somaObjectTypeRequest{
		action:     "rename",
		objectType: params.ByName("objectType"),
		rename:     clientRequest.Entity.Name,
		reply:      returnChannel,
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
