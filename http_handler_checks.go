package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */
func ListCheckConfiguration(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	result := <-returnChannel

	SendCheckConfigurationReply(&w, &result)
}

func ShowCheckConfiguration(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	result := <-returnChannel

	SendCheckConfigurationReply(&w, &result)
}

/* Write functions
 */
func AddCheckConfiguration(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	returnChannel := make(chan somaResult)
	result := <-returnChannel

	SendCheckConfigurationReply(&w, &result)
}

/* Utility
 */
func SendCheckConfigurationReply(w *http.ResponseWriter, r *somaResult) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
