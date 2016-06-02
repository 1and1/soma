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
	"log"
	"time"

)

func (s *supervisor) issue_token(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `issue_token`}
	var (
		cred                 *svCredential
		err                  error
		kex                  *auth.Kex
		plain                []byte
		timer                *time.Timer
		token                auth.Token
		tx                   *sql.Tx
		validFrom, expiresAt time.Time
	)
	data := q.Super.Data

	// issue_token is a master instance function
	if s.readonly {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto conflict
	}
	// start response timer
	timer = time.NewTimer(1 * time.Second)
	defer timer.Stop()

	// -> get kex
	if kex = s.kex.read(q.Super.KexId); kex == nil {
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// check kex.SameSource
	if !kex.IsSameSourceString(q.Super.RemoteAddr) {
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// delete kex from s.kex (kex is now used)
	s.kex.remove(q.Super.KexId)
	// decrypt request
	if err = kex.DecodeAndDecrypt(&data, &plain); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> json.Unmarshal(rdata, &token)
	if err = json.Unmarshal(plain, &token); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if token.UserName == `root` && s.root_restricted && !q.Super.Restricted {
		result.ServerError(
			fmt.Errorf(`Root token requested on unrestricted endpoint`))
		goto dispatch
	}

	log.Printf(LogStrReq, q.Type, fmt.Sprintf("%s/%s", `authenticate`, q.Action), token.UserName, q.Super.RemoteAddr)

	if cred = s.credentials.read(token.UserName); cred == nil {
		result.Unauthorized(fmt.Errorf("Unknown user: %s", token.UserName))
		goto dispatch
	}
	if !cred.isActive {
		result.Unauthorized(fmt.Errorf("Inactive user: %s", token.UserName))
		goto dispatch
	}
	if time.Now().UTC().Before(cred.validFrom.UTC()) ||
		time.Now().UTC().After(cred.expiresAt.UTC()) {
		result.Unauthorized(fmt.Errorf("Expired: %s", token.UserName))
		goto dispatch
	}
	// generate token
	token.SetIPAddressString(q.Super.RemoteAddr)
	if err = token.Generate(cred.cryptMCF, s.key, s.seed); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	validFrom, _ = time.Parse(rfc3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(rfc3339Milli, token.ExpiresAt)
	// -> DB Insert: token data
	if tx, err = s.conn.Begin(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	defer tx.Rollback()
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
