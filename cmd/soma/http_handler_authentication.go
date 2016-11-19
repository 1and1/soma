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

// AuthenticationValidate is a noop function wrapped in HTTP basic
// authentication that can be used to verify one's credentials
func AuthenticationValidate(w http.ResponseWriter, _ *http.Request,
	_ httprouter.Params) {
	w.WriteHeader(http.StatusNoContent)
}

// AuthenticationKex is used by the client to initiate a key exchange
// that can the be used for one of the encrypted endpoints
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
		Section: `kex`,
		Action:  `init`,
		Reply:   returnChannel,
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

// AuthenticationBootstrapRoot is the encrypted endpoint used during
// service setup to access the builtin root account
func AuthenticationBootstrapRoot(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `bootstrap_root`)
}

// AuthenticationIssueToken is the encrypted endpoint used to
// request a password token
func AuthenticationIssueToken(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `request_token`)
}

// AuthenticationActivateUser is the encrypted endpoint used to
// activate a user account using external ownership verification
func AuthenticationActivateUser(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `activate_user`)
}

// AuthenticationChangeUserPassword is the encrypted endpoint used
// to change the account password using the current one.
func AuthenticationChangeUserPassword(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `change_user_password`)
}

// AuthenticationResetUserPassword is the encrypted endpoint used
// to change the account password using external ownership
// verification
func AuthenticationResetUserPassword(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer PanicCatcher(w)

	AuthenticationEncryptedData(&w, r, &params, `reset_user_password`)
}

// AuthenticationEncryptedData is the generic function for
// encrypted endpoints
func AuthenticationEncryptedData(w *http.ResponseWriter,
	r *http.Request, params *httprouter.Params, request string) {
	defer PanicCatcher(*w)

	data := make([]byte, r.ContentLength)
	io.ReadFull(r.Body, data)

	var section, action string
	switch request {
	case `reset_user_password`:
		section = `password`
		action = `reset`
	case `change_user_password`:
		section = `password`
		action = `change`
	case `bootstrap_root`:
		section = `bootstrap`
		action = `root`
	case `request_token`:
		section = `token`
		action = `request`
	case `activate_user`:
		section = `activate`
		action = `user`
	}

	returnChannel := make(chan msg.Result)
	handler := handlerMap[`supervisor`].(*supervisor)
	handler.input <- msg.Request{
		Section: section,
		Action:  action,
		Reply:   returnChannel,
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
