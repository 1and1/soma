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
func ListViews(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaViewResult)

	handler := handlerMap["viewReadHandler"].(somaViewReadHandler)
	handler.input <- somaViewRequest{
		action: "list",
		reply:  returnChannel,
	}

	results := <-returnChannel
	views := make([]string, len(results))
	for pos, res := range results {
		views[pos] = res.view
	}
	json, err := json.Marshal(somaproto.ProtoResultViewList{Code: 200, Status: "OK", Views: views})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ShowView(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaViewResult)

	handler := handlerMap["viewReadHandler"].(somaViewReadHandler)
	handler.input <- somaViewRequest{
		action: "show",
		view:   params.ByName("view"),
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
	json, err := json.Marshal(somaproto.ProtoResultViewDetail{
		Code:    200,
		Status:  "OK",
		Details: somaproto.ProtoViewDetails{View: result.view},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func AddView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaViewResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestView
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["viewWriteHandler"].(somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "add",
		view:   clientRequest.View,
		reply:  returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultView{
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
		json, _ := json.Marshal(somaproto.ProtoResultView{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Added view: %s", result.view)
	json, _ := json.Marshal(somaproto.ProtoResultView{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func DeleteView(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaViewResult)

	handler := handlerMap["viewWriteHandler"].(somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "delete",
		view:   params.ByName("view"),
		reply:  returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultView{
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
		json, _ := json.Marshal(somaproto.ProtoResultView{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Deleted view: %s", result.view)
	json, _ := json.Marshal(somaproto.ProtoResultView{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func RenameView(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaViewResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestView
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["viewWriteHandler"].(somaViewWriteHandler)
	handler.input <- somaViewRequest{
		action: "rename",
		view:   params.ByName("view"),
		rename: clientRequest.View,
		reply:  returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultView{
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
		json, _ := json.Marshal(somaproto.ProtoResultView{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Renamed view: %s to %s", result.view, clientRequest.View)
	json, _ := json.Marshal(somaproto.ProtoResultView{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
