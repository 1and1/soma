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
func ListDatacenters(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaDatacenterResult)

	handler := handlerMap["datacenterReadHandler"].(somaDatacenterReadHandler)
	handler.input <- somaDatacenterRequest{
		action: "list",
		reply:  returnChannel,
	}

	results := <-returnChannel
	datacenters := make([]string, len(results))
	for pos, res := range results {
		datacenters[pos] = res.datacenter
	}
	json, err := json.Marshal(somaproto.ProtoResultDatacenterList{Code: 200, Status: "OK", Datacenters: datacenters})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ShowDatacenter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaDatacenterResult)

	handler := handlerMap["datacenterReadHandler"].(somaDatacenterReadHandler)
	handler.input <- somaDatacenterRequest{
		action:     "show",
		datacenter: params.ByName("datacenter"),
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
	json, err := json.Marshal(somaproto.ProtoResultDatacenterDetail{
		Code:    200,
		Status:  "OK",
		Details: somaproto.ProtoDatacenterDetails{Datacenter: result.datacenter},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func AddDatacenter(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnChannel := make(chan []somaDatacenterResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestDatacenter
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["datacenterWriteHandler"].(somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action:     "add",
		datacenter: clientRequest.Datacenter,
		reply:      returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
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
		json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Added datacenter: %s", result.datacenter)
	json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func DeleteDatacenter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaDatacenterResult)

	handler := handlerMap["datacenterWriteHandler"].(somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action:     "delete",
		datacenter: params.ByName("datacenter"),
		reply:      returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
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
		json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Deleted datacenter: %s", result.datacenter)
	json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func RenameDatacenter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	returnChannel := make(chan []somaDatacenterResult)

	// read POST body
	decoder := json.NewDecoder(r.Body)
	var clientRequest somaproto.ProtoRequestDatacenter
	err := decoder.Decode(&clientRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	handler := handlerMap["datacenterWriteHandler"].(somaDatacenterWriteHandler)
	handler.input <- somaDatacenterRequest{
		action:     "rename",
		datacenter: params.ByName("datacenter"),
		rename:     clientRequest.Datacenter,
		reply:      returnChannel,
	}

	results := <-returnChannel
	if len(results) != 1 {
		json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
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
		json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
			Code:   500,
			Status: "Internal Server Error",
			Text:   []string{result.err.Error()},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		return
	}

	txt := fmt.Sprintf("Renamed datacenter: %s to %s", result.datacenter, clientRequest.Datacenter)
	json, _ := json.Marshal(somaproto.ProtoResultDatacenter{
		Code:   200,
		Status: "OK",
		Text:   []string{txt},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
