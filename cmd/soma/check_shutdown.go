package main

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	metrics "github.com/rcrowley/go-metrics"
)

func Check(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		if !ShutdownInProgress {
			metrics.GetOrRegisterCounter(`.requests`, Metrics[`soma`]).Inc(1)
			start := time.Now()

			h(w, r, ps)

			metrics.GetOrRegisterTimer(`.requests.latency`,
				Metrics[`soma`]).UpdateSince(start)
			return
		}

		http.Error(w, `Shutdown in progress`,
			http.StatusServiceUnavailable)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
