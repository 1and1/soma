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
	"encoding/json"
	"fmt"
	"time"

	"github.com/mjolnir42/scrypth64"
	"github.com/satori/go.uuid"
)

func (s *supervisor) bootstrapRoot(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `bootstrap_root`}
	kexId := q.Super.KexId
	data := q.Super.Data
	var kex *auth.Kex
	var err error
	var plain []byte
	var token auth.Token
	var rootToken string
	var mcf scrypth64.Mcf
	var tx *sql.Tx
	var validFrom, expiresAt time.Time
	var timer *time.Timer

	// bootstrapRoot is a master instance function
	if s.readonly {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto conflict
	}

	// start response timer
	timer = time.NewTimer(1 * time.Second)
	defer timer.Stop()

	// -> check if root is not already active
	if s.credentials.read(`root`) != nil {
		result.BadRequest(fmt.Errorf(`Root account is already active`))
		//    --> delete kex
		s.kex.remove(kexId)
		goto dispatch
	}
	// -> get kex
	if kex = s.kex.read(kexId); kex == nil {
		//    --> reply 404 if not found
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// -> check kex.SameSource
	if !kex.IsSameSourceExtractedString(q.Super.RemoteAddr) {
		//    --> reply 404 if !SameSource
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// -> delete kex from s.kex (kex is now used)
	s.kex.remove(kexId)
	// -> rdata = kex.DecodeAndDecrypt(data)
	if err = kex.DecodeAndDecrypt(&data, &plain); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> json.Unmarshal(rdata, &token)
	if err = json.Unmarshal(plain, &token); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> check token.UserName == `root`
	if token.UserName != `root` {
		//    --> reply 401
		result.Unauthorized(nil)
		goto dispatch
	}
	if token.UserName == `root` && s.root_restricted && !q.Super.Restricted {
		result.ServerError(
			fmt.Errorf(`Root bootstrap requested on unrestricted endpoint`))
		goto dispatch
	}
	// -> check token.Token is correct bearer token
	if rootToken, err = s.fetchRootToken(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if token.Token != rootToken || len(token.Password) == 0 {
		//    --> reply 401
		result.Unauthorized(nil)
		goto dispatch
	}
	// -> scrypth64.Digest(Password, nil)
	if mcf, err = scrypth64.Digest(token.Password, nil); err != nil {
		result.Unauthorized(nil)
		goto dispatch
	}
	// -> generate token
	token.SetIPAddressExtractedString(q.Super.RemoteAddr)
	if err = token.Generate(mcf, s.key, s.seed); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	validFrom, _ = time.Parse(rfc3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(rfc3339Milli, token.ExpiresAt)

	// -> DB Insert: root password data
	if tx, err = s.conn.Begin(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	defer tx.Rollback()
	if _, err = tx.Exec(
		stmt.SetRootCredentials,
		uuid.Nil,
		mcf.String(),
		validFrom.UTC(),
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> DB Insert: token data
	if _, err = tx.Exec(
		stmt.InsertToken,
		token.Token,
		token.Salt,
		validFrom.UTC(),
		expiresAt.UTC(),
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> s.credentials Update
	s.credentials.insert(`root`, uuid.Nil, validFrom.UTC(),
		PosTimeInf.UTC(), mcf)
	// -> s.tokens Update
	if err = s.tokens.insert(token.Token, token.ValidFrom, token.ExpiresAt,
		token.Salt); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if err = tx.Commit(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> sdata = kex.EncryptAndEncode(&token)
	plain = []byte{}
	data = []byte{}
	if plain, err = json.Marshal(token); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if err = kex.EncryptAndEncode(&plain, &data); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> send sdata reply
	result.Super = &msg.Supervisor{
		Verdict: 200,
		Data:    data,
	}
	result.OK()

dispatch:
	<-timer.C

conflict:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
