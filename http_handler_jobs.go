package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func JobDelay(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan bool)
	handler := handlerMap[`jobDelay`].(jobDelay)
	handler.input <- waitSpec{
		JobId: params.ByName(`jobid`),
		Reply: returnChannel,
	}
	<-returnChannel
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

/* Read functions
 */
func ListJobs(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`jobs_list`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
}

func ShowJob(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`jobs_show`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
}

func SearchJob(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
	if ok, _ := IsAuthorized(params.ByName(`AuthenticatedUser`),
		`jobs_search`, ``, ``, ``); !ok {
		DispatchForbidden(&w, nil)
		return
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
