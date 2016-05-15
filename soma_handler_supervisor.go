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
	"database/sql"
	"time"

)

type supervisor struct {
	input       chan msg.Request
	shutdown    chan bool
	conn        *sql.DB
	seed        []byte
	key         []byte
	readonly    bool
	tokenExpiry uint64
	kexExpiry   uint64
	kexMap      map[string]auth.Kex
	tokenMap    map[string]svToken
}

type svToken struct {
	binToken     []byte
	validFrom    time.Time
	binExpiresAt []byte
	salt         []byte
}

func (s *supervisor) run() {
	//var err error

	auth.TokenExpirySeconds = s.tokenExpiry
	auth.KexExpirySeconds = s.kexExpiry

runloop:
	for {
		select {
		case <-s.shutdown:
			break runloop
		case req := <-s.input:
			go func() {
				s.process(&req)
			}()
		}
	}
}

func (s *supervisor) process(q *msg.Request) {
	switch q.Action {
	case `kex_init`:
		s.kexInit(q)
	case `bootstrap_root`:
		s.bootstrapRoot(q)
	}
}

func (s *supervisor) kexInit(q *msg.Request) {
	result := msg.Result{Type: `supervisor`}
	// TODO
	q.Reply <- result
}

func (s *supervisor) bootstrapRoot(q *msg.Request) {
	result := msg.Result{Type: `supervisor`}
	// TODO
	q.Reply <- result
}

// TODO: timer
// delete all expired key exchanges
func (s *supervisor) pruneKex() {
	for kexId, kex := range s.kexMap {
		if kex.IsExpired() {
			delete(s.kexMap, kexId)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
