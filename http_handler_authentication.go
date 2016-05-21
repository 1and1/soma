/*-
Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"


	"github.com/julienschmidt/httprouter"
)

/* Write functions
 */

func AuthenticationKex(w http.ResponseWriter, r *http.Request,
	_ httprouter.Params) {
	defer PanicCatcher(w)

	kex := auth.Kex{}
	err := DecodeJsonBody(r, &kex)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(supervisor)
	handler.input <- msg.Request{
		Type:   `supervisor`,
		Action: `kex_init`,
		Reply:  returnChannel,
		Super: &msg.Supervisor{
			RemoteAddr: r.RemoteAddr,
			Kex: auth.Kex{
				Public:               kex.Public,
				InitializationVector: kex.InitializationVector,
			},
		},
	}

	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func AuthenticationBootstrapRoot(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	data := make([]byte, r.ContentLength)
	io.ReadFull(r.Body, data)

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(supervisor)
	handler.input <- msg.Request{
		Type:   `supervisor`,
		Action: `bootstrap_root`,
		Reply:  returnChannel,
		Super: &msg.Supervisor{
			RemoteAddr: r.RemoteAddr,
			KexId:      params.ByName(`uuid`),
			Data:       data,
		},
	}
	result := <-returnChannel
	SendMsgResult(&w, &result)
}

func AuthenticationIssueToken(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
}

func AuthenticationActivateUser(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)
}

func AuthenticationValidate(w http.ResponseWriter, _ *http.Request,
	_ httprouter.Params) {
	w.WriteHeader(http.StatusNoContent)
}

/* Utility
 */

func SendMsgResult(w *http.ResponseWriter, r *msg.Result) {
	var (
		bjson []byte
		err   error
		k     auth.Kex
	)

	// this is central error command, proceeding to log
	if r.Error != nil {
		log.Printf(LogStrErr, `supervisor`, r.Action, r.Code, r.Error.Error())
	}

	switch r.Type {
	case `supervisor`:
		switch r.Action {
		case `kex_reply`:
			k = r.Super.Kex
			if bjson, err = json.Marshal(&k); err != nil {
				log.Printf(LogStrErr, `supervisor`, r.Action, r.Code, err)
				DispatchInternalError(w, nil)
				return
			}
			goto dispatchJSON
		case `bootstrap_root`:
			// for this request type, errors are masked in responses
			switch r.Code {
			case 200:
				if r.Super.Verdict == 200 {
					log.Printf(LogStrOK, `supervisor`, r.Action, r.Code, 200)
					goto dispatchOCTET
				}
				log.Printf(LogStrOK, `supervisor`, r.Action, r.Code, 401)
				DispatchUnauthorized(w, nil)
			case 400:
				log.Printf(LogStrOK, `supervisor`, r.Action, r.Code, 400)
				DispatchBadRequest(w, nil)
			case 404:
				log.Printf(LogStrOK, `supervisor`, r.Action, r.Code, 404)
				DispatchNotFound(w, r.Error)
			case 406:
				log.Printf(LogStrOK, `supervisor`, r.Action, r.Code, 406)
				DispatchConflict(w, r.Error)
			default:
				log.Printf(LogStrOK, `supervisor`, r.Action, r.Code, 401)
				DispatchUnauthorized(w, nil)
			}
			return
		}
	}

dispatchOCTET:
	DispatchOctetReply(w, &r.Super.Data)
	return

dispatchJSON:
	DispatchJsonReply(w, &bjson)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
