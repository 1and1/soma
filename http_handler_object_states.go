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
func ListObjectStates(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaObjectStateResult)

	handler := handlerMap["objectStateReadHandler"].(somaObjectStateReadHandler)
	handler.input <- somaObjectStateRequest{
		action: "list",
		reply:  returnChannel,
	}

	results := <-returnChannel
	objectStates := make([]string, len(results))
	for pos, res := range results {
		objectStates[pos] = res.state
	}
	json, err := json.Marshal(somaproto.ProtoResultObjectStateList{Code: 200, Status: "OK", States: objectStates})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ShowObjectState(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	result := results[0]
	json, err := json.Marshal(somaproto.ProtoResultObjectStateDetail{
		Code:    200,
		Status:  "OK",
		Details: somaproto.ProtoObjectStateDetails{State: result.state},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func AddObjectState(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaObjectStateResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestObjectState
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["objectStateWriteHandler"].(somaObjectStateWriteHandler)
	handler.input <- somaObjectStateRequest{
		action: "add",
		state:  clientRequest.State,
		reply:  returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultObjectState{
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
		json, _ := json.Marshal(somaproto.ProtoResultObjectState{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Added objectState: %s", result.state)
	json, _ := json.Marshal(somaproto.ProtoResultObjectState{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func DeleteObjectState(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaObjectStateResult)

	handler := handlerMap["objectStateWriteHandler"].(somaObjectStateWriteHandler)
	handler.input <- somaObjectStateRequest{
		action: "delete",
		state:  params.ByName("state"),
		reply:  returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultObjectState{
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
		json, _ := json.Marshal(somaproto.ProtoResultObjectState{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Deleted objectState: %s", result.state)
	json, _ := json.Marshal(somaproto.ProtoResultObjectState{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func RenameObjectState(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaObjectStateResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestObjectState
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["objectStateWriteHandler"].(somaObjectStateWriteHandler)
	handler.input <- somaObjectStateRequest{
		action: "rename",
		state:  params.ByName("state"),
		rename: clientRequest.State,
		reply:  returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultObjectState{
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
		json, _ := json.Marshal(somaproto.ProtoResultObjectState{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Renamed objectState: %s to %s", result.state, clientRequest.State)
	json, _ := json.Marshal(somaproto.ProtoResultObjectState{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
