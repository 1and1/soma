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
	"fmt"

	"github.com/1and1/soma/lib/msg"
	"github.com/1and1/soma/lib/auth"
)

func (s *supervisor) kexInit(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `kex_reply`}
	kex := q.Super.Kex
	var err error

	// kexInit is a master instance function
	if s.readonly {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto dispatch
	}

	// check the client submitted IV for fishyness
	err = s.checkIV(kex.InitializationVector)
	for err != nil {
		if err = kex.GenerateNewVector(); err != nil {
			continue
		}
		err = s.checkIV(kex.InitializationVector)
	}

	// record the kex submission time
	kex.SetTimeUTC()

	// record the client ip address
	kex.SetIPAddressExtractedString(q.Super.RemoteAddr)

	// generate a request ID
	kex.GenerateNewRequestID()

	// set the client submitted public key as peer key
	kex.SetPeerKey(kex.PublicKey())

	// generate our own keypair
	kex.GenerateNewKeypair()

	// save kex
	s.kex.insert(kex)

	// send out reply
	result.Super = &msg.Supervisor{
		Verdict: 200,
		Kex: auth.Kex{
			Public:               kex.Public,
			InitializationVector: kex.InitializationVector,
			Request:              kex.Request,
		},
	}
	result.OK()

dispatch:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
