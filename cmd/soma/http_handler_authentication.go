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
	"io"
	"net/http"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/lib/auth"

	"github.com/julienschmidt/httprouter"
)

/* Read functions
 */

func AuthenticationValidate(w http.ResponseWriter, _ *http.Request,
	_ httprouter.Params) {
	w.WriteHeader(http.StatusNoContent)
}

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
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:   `supervisor`,
		Action: `kex_init`,
		Reply:  returnChannel,
		Super: &msg.Supervisor{
			RemoteAddr: extractAddress(r.RemoteAddr),
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

	AuthenticationEncryptedData(&w, r, &params, `bootstrap_root`)
}

func AuthenticationIssueToken(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `request_token`)
}

func AuthenticationActivateUser(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `activate_user`)
}

func AuthenticationChangeUserPassword(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `change_user_password`)
}

func AuthenticationResetUserPassword(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `reset_user_password`)
}

func AuthenticationEncryptedData(w *http.ResponseWriter, r *http.Request,
	params *httprouter.Params, action string) {
	defer PanicCatcher(*w)

	data := make([]byte, r.ContentLength)
	io.ReadFull(r.Body, data)

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Type:   `supervisor`,
		Action: action,
		Reply:  returnChannel,
		Super: &msg.Supervisor{
			RemoteAddr: extractAddress(r.RemoteAddr),
			KexId:      params.ByName(`uuid`),
			Data:       data,
			Restricted: false,
		},
	}
	result := <-returnChannel
	SendMsgResult(w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
