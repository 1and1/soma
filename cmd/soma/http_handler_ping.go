package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Ping is the function for HEAD requests on the base API that
// reports facts about the running application
func Ping(w http.ResponseWriter, _ *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	w.Header().Set(`X-Powered-By`, `SOMA Configuration System`)
	w.Header().Set(`X-Version`, somaVersion)
	switch {
	case SomaCfg.Observer == true:
		w.Header().Set(`X-SOMA-Mode`, `Observer`)
	case SomaCfg.ReadOnly == true:
		w.Header().Set(`X-SOMA-Mode`, `ReadOnly`)
	case SomaCfg.ReadOnly == false:
		w.Header().Set(`X-SOMA-Mode`, `Master`)
	}
	w.WriteHeader(http.StatusNoContent)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
