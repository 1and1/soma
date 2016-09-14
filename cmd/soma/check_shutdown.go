package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Check(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		if !ShutdownInProgress {
			h(w, r, ps)
			return
		}

		http.Error(w, `Shutdown in progress`,
			http.StatusServiceUnavailable)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
