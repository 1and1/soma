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

	"github.com/mjolnir42/scrypth64"
	"github.com/satori/go.uuid"
)

func (s *supervisor) activate_user(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `activate_user`, Super: &msg.Supervisor{Action: ``}}

	var (
		timer                               *time.Timer
		plain                               []byte
		err                                 error
		kex                                 *auth.Kex
		validFrom, expiresAt, credExpiresAt time.Time
		token                               auth.Token
		userId                              string
		userUUID                            uuid.UUID
		ok                                  bool
		mcf                                 scrypth64.Mcf
		tx                                  *sql.Tx
	)
	data := q.Super.Data

	if s.readonly {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto conflict
	}

	// start response timer
	timer = time.NewTimer(1 * time.Second)
	defer timer.Stop()

	// -> get kex
	if kex = s.kex.read(q.Super.KexId); kex == nil {
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
	s.kex.remove(q.Super.KexId)
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
	// request has been decrypted, log it
	log.Printf(LogStrReq, q.Type, q.Action, token.UserName, q.Super.RemoteAddr)

	// -> check token.UserName != `root`
	if token.UserName == `root` {
		//    --> reply 401
		result.Unauthorized(fmt.Errorf(`Cannot activate root`))
		goto dispatch
	}

	// check we have the user
	if err = s.stmt_FindUser.QueryRow(token.UserName).Scan(&userId); err == sql.ErrNoRows {
		result.Unauthorized(fmt.Errorf("Unknown user: %s", token.UserName))
		goto dispatch
	} else if err != nil {
		result.ServerError(err)
	}
	userUUID, _ = uuid.FromString(userId)

	// no account ownership verification in open mode
	if !SomaCfg.OpenInstance {
		switch s.activation {
		case `ldap`:
			if ok, err = validateLdapCredentials(token.UserName, token.Token); err != nil {
				result.ServerError(err)
				goto dispatch
			} else if !ok {
				result.Unauthorized(fmt.Errorf(`Invalid LDAP credentials`))
				goto dispatch
			}
			// fail activation if local password is the same as the
			// upstream password
			if token.Token == token.Password {
				result.Unauthorized(fmt.Errorf("User %s denied: matching local/upstream passwords", token.UserName))
				goto dispatch
			}
		case `token`: // TODO
			result.ServerError(fmt.Errorf(`Not implemented`))
			goto dispatch
		default:
			result.ServerError(fmt.Errorf("Unknown activation: %s",
				SomaCfg.Auth.Activation))
			goto dispatch
		}
	}
	// OK: validation success

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
	credExpiresAt = validFrom.Add(time.Duration(s.credExpiry) * time.Hour * 24).UTC()

	// -> open transaction
	if tx, err = s.conn.Begin(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	defer tx.Rollback()
	// -> DB Insert: password data
	if _, err = tx.Exec(
		stmt.SetUserCredential,
		userUUID,
		mcf.String(),
		validFrom.UTC(),
		credExpiresAt.UTC(),
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> DB Update: activate user
	if _, err = tx.Exec(
		stmt.ActivateUser,
		userUUID,
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
	s.credentials.insert(token.UserName, userUUID, validFrom.UTC(),
		credExpiresAt.UTC(), mcf)
	// -> s.tokens Update
	if err = s.tokens.insert(token.Token, token.ValidFrom, token.ExpiresAt,
		token.Salt); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// commit transaction
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
	result.Super.Verdict = 200
	result.Super.Data = data
	result.OK()

dispatch:
	<-timer.C

conflict:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
