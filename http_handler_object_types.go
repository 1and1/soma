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
func ListObjectTypes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaObjectTypeResult)

	handler := handlerMap["objectTypeReadHandler"].(somaObjectTypeReadHandler)
	handler.input <- somaObjectTypeRequest{
		action: "list",
		reply:  returnChannel,
	}

	results := <-returnChannel
	objectTypes := make([]string, len(results))
	for pos, res := range results {
		objectTypes[pos] = res.objectType
	}
	json, err := json.Marshal(somaproto.ProtoResultObjectTypeList{Code: 200, Status: "OK", Types: objectTypes})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ShowObjectType(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaObjectTypeResult)

	handler := handlerMap["objectTypeReadHandler"].(somaObjectTypeReadHandler)
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
	json, err := json.Marshal(somaproto.ProtoResultObjectTypeDetail{
		Code:    200,
		Status:  "OK",
		Details: somaproto.ProtoObjectTypeDetails{Type: result.objectType},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func AddObjectType(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaObjectTypeResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestObjectType
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["objectTypeWriteHandler"].(somaObjectTypeWriteHandler)
	handler.input <- somaObjectTypeRequest{
		action:     "add",
		objectType: clientRequest.Type,
		reply:      returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultObjectType{
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
		json, _ := json.Marshal(somaproto.ProtoResultObjectType{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Added objectType: %s", result.objectType)
	json, _ := json.Marshal(somaproto.ProtoResultObjectType{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func DeleteObjectType(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaObjectTypeResult)

	handler := handlerMap["objectTypeWriteHandler"].(somaObjectTypeWriteHandler)
	handler.input <- somaObjectTypeRequest{
		action:     "delete",
		objectType: params.ByName("objectType"),
		reply:      returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultObjectType{
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
		json, _ := json.Marshal(somaproto.ProtoResultObjectType{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Deleted objectType: %s", result.objectType)
	json, _ := json.Marshal(somaproto.ProtoResultObjectType{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func RenameObjectType(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaObjectTypeResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestObjectType
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["objectTypeWriteHandler"].(somaObjectTypeWriteHandler)
	handler.input <- somaObjectTypeRequest{
		action:     "rename",
		objectType: params.ByName("objectType"),
		rename:     clientRequest.Type,
		reply:      returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultObjectType{
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
		json, _ := json.Marshal(somaproto.ProtoResultObjectType{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Renamed objectType: %s to %s", result.objectType, clientRequest.Type)
	json, _ := json.Marshal(somaproto.ProtoResultObjectType{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
