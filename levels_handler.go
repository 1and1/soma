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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s\n", debug.Stack())
			msg := fmt.Sprintf("PANIC! %s", r)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}()

	returnChannel := make(chan []somaLevelResult)
	handler := handlerMap["levelReadHandler"].(somaLevelReadHandler)
	handler.input <- somaLevelRequest{
		action: "list",
		reply:  returnChannel,
	}
	results := <-returnChannel

	var res somaproto.ProtoResultLevel
	if len(results) == 0 {
		res.Code = 404
		res.Status = "NOTFOUND"
		goto submission
	} else if results[0].rErr != nil {
		res.Code = 500
		res.Status = "ERROR"
		res.Text = make([]string, 0)
		res.Text = append(res.Text, results[0].rErr.Error())
		goto submission
	}
	res.Code = 200
	res.Status = "OK"
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
