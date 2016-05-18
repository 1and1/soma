package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Ping(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	defer PanicCatcher(w)

	w.Header().Set(`X-Powered-By`, `SOMA Configuration System`)
	w.Header().Set(`X-Version`, SomaVersion)
	switch SomaCfg.ReadOnly {
	case true:
		w.Header().Set(`X-SOMA-Mode`, `ReadOnly`)
	case false:
		w.Header().Set(`X-SOMA-Mode`, `Master`)
	}
	w.WriteHeader(http.StatusOK)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
