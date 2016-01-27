package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

/*
 * Read functions
 */
func ListLevels(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaLevelResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "list",
		reply:  returnChannel,
	}
	results := <-returnChannel

	var res somaproto.ProtoResultLevel
	var res somaproto.ProtoResultLevel
	dispatchError := CheckErrorHandler(&results, &res)
	if dispatchError {
		goto submission
	}
	res.Text = make([]string, 0)
	res.Levels = make([]somaproto.ProtoLevel, 0)
	for _, l := range results {
		res.Levels = append(res.Levels, l.level)
		if l.lErr != nil {
			res.Text = append(res.Text, l.lErr.Error())
		}
	}

submission:
	json, jErr := json.Marshal(res)
	if jErr != nil {
		http.Error(w, jErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func ShowLevels(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan []somaLevelResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "show",
		reply:  returnChannel,
		level: somaproto.ProtoLevel{
			Name: params.ByName("level"),
		},
	}
	results := <-returnChannel

	var res somaproto.ProtoResultLevel
	dispatchError := CheckErrorHandler(&results, &res)
	if dispatchError {
		goto submission
	}
	res.Text = make([]string, 0)
	res.Levels = make([]somaproto.ProtoLevel, 0)
	for _, l := range results {
		res.Levels = append(res.Levels, l.level)
		if l.lErr != nil {
			res.Text = append(res.Text, l.lErr.Error())
		}
	}

submission:
	json, jErr := json.Marshal(res)
	if jErr != nil {
		http.Error(w, jErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
