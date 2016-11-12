/*-
Copyright (c) 2013 Julien Schmidt. All rights reserved.
Copyright (c) 2016 Jörg Pernfuß <joerg.pernfuss@1und1.de>


Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * The names of the contributors may not be used to endorse or promote
      products derived from this software without specific prior written
      permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL JULIEN SCHMIDT BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


The following code is nearly verbatim the example code from the httprouter
distribution. Therefor copyright is set to the license text of that distribution.
*/

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/1and1/soma/internal/msg"
	log "github.com/Sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

func BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {
		const basicAuthPrefix string = "Basic "

		// if the supervisor is not available, no requests are accepted
		if _, ok := handlerMap[`supervisor`]; !ok {
			http.Error(w, `Authentication supervisor not available`,
				http.StatusServiceUnavailable)
			return
		}

		// disable authentication much?
		if SomaCfg.OpenInstance {
			ps = append(ps, httprouter.Param{
				Key:   `AuthenticatedUser`,
				Value: `AnonymousCoward`,
			})
			h(w, r, ps)
			return
		}

		// Get credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 {
					returnChannel := make(chan msg.Result)
					super := handlerMap[`supervisor`].(*supervisor)
					super.input <- msg.Request{
						Type:    `supervisor`,
						Section: `authenticate`,
						Action:  `basic`,
						Reply:   returnChannel,
						Super: &msg.Supervisor{
							RemoteAddr:     extractAddress(r.RemoteAddr),
							Restricted:     false,
							BasicAuthUser:  string(pair[0]),
							BasicAuthToken: string(pair[1]),
						},
					}
					result := <-returnChannel
					if result.Error != nil {
						// log authentication errors
						log.Printf(LogStrErr, fmt.Sprintf("BasicAuth:%s", string(pair[0])), result.Action, result.Code, result.Error.Error())
					}
					if result.Super.Verdict == 200 {
						// record the authenticated user
						ps = append(ps, httprouter.Param{
							Key:   `AuthenticatedUser`,
							Value: string(pair[0]),
						})
						// Delegate request to given handle
						h(w, r, ps)
						return

					}
				}
			}
		}

		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
