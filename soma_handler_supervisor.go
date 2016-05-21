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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/mjolnir42/scrypth64"
	"github.com/satori/go.uuid"

)

type supervisor struct {
	input           chan msg.Request
	shutdown        chan bool
	conn            *sql.DB
	seed            []byte
	key             []byte
	readonly        bool
	tokenExpiry     uint64
	kexExpiry       uint64
	root_disabled   bool
	root_restricted bool
	kex             svKexMap
	tokens          svTokenMap
	credentials     svCredMap
	stmt_FToken     *sql.Stmt
}

func (s *supervisor) run() {
	var err error

	// set library options
	auth.TokenExpirySeconds = s.tokenExpiry
	auth.KexExpirySeconds = s.kexExpiry

	// initialize maps
	s.tokens = s.newTokenMap()
	s.credentials = s.newCredentialMap()
	s.kex = s.newKexMap()

	// prepare SQL statements
	if s.stmt_FToken, err = s.conn.Prepare(stmt.SelectToken); err != nil {
		log.Fatal("supervisor/fetch-token: ", err)
	}
	defer s.stmt_FToken.Close()

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
	case `basic_auth`:
		s.validate_basic_auth(q)
	}
}

func (s *supervisor) kexInit(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `kex_reply`}
	kex := q.Super.Kex
	var err error

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
	kex.SetIPAddressString(q.Super.RemoteAddr)

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

	q.Reply <- result
}

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

	// start response timer
	timer := time.NewTimer(1 * time.Second)
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
	if !kex.IsSameSourceString(q.Super.RemoteAddr) {
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
	token.SetIPAddressString(q.Super.RemoteAddr)
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
	q.Reply <- result
}

// TODO: timer
// delete all expired key exchanges
func (s *supervisor) pruneKex() {
	s.kex.lock()
	defer s.kex.unlock()
	for kexId, kex := range s.kex.KMap {
		if kex.IsExpired() {
			delete(s.kex.KMap, kexId)
		}
	}
}

func (s *supervisor) newTokenMap() svTokenMap {
	m := svTokenMap{}
	m.TMap = make(map[string]svToken)
	return m
}

func (s *supervisor) newCredentialMap() svCredMap {
	m := svCredMap{}
	m.CMap = make(map[string]svCredential)
	return m
}

func (s *supervisor) newKexMap() svKexMap {
	m := svKexMap{}
	m.KMap = make(map[string]auth.Kex)
	return m
}

func (s *supervisor) validate_basic_auth(q *msg.Request) {
	result := msg.Result{Type: `supervisor`, Action: `authenticate`}

	tok := s.tokens.read(q.Super.BasicAuthToken)
	if tok == nil && !s.readonly {
		// rw instance knows every token
		result.ServerError(fmt.Errorf(`Unknown Token (TokenMap)`))
		goto unauthorized
	} else if tok == nil {
		if !s.fetchTokenFromDB(q.Super.BasicAuthToken) {
			result.ServerError(fmt.Errorf(`Unknown Token (pgSQL)`))
			goto unauthorized
		}
		tok = s.tokens.read(q.Super.BasicAuthToken)
	}
	if time.Now().UTC().Before(tok.validFrom.UTC()) ||
		time.Now().UTC().After(tok.expiresAt.UTC()) {
		result.Unauthorized(fmt.Errorf(`Token expired`))
		goto unauthorized
	}

	if auth.Verify(q.Super.BasicAuthUser, q.Super.RemoteAddr, tok.binToken, s.key,
		s.seed, tok.binExpiresAt, tok.salt) {
		// valid token
		result.Super = &msg.Supervisor{Verdict: 200}
		result.OK()
		q.Reply <- result
	}

unauthorized:
	result.Super = &msg.Supervisor{Verdict: 401}
	q.Reply <- result
}

func (s *supervisor) fetchTokenFromDB(token string) bool {
	var (
		err                       error
		salt, strValid, strExpire string
		validF, validU            time.Time
	)

	err = s.stmt_FToken.QueryRow(token).Scan(&salt, &validF, &validU)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		// XXX log error
		return false
	}

	strValid = validF.UTC().Format(rfc3339Milli)
	strExpire = validU.UTC().Format(rfc3339Milli)

	if err = s.tokens.insert(token, strValid, strExpire, salt); err == nil {
		return true
	}
	return false
}

func (s *supervisor) fetchRootToken() (string, error) {
	var (
		err   error
		token string
	)

	err = s.conn.QueryRow(stmt.SelectRootToken).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// the nonces used for encryption are implemented as
// a counter on top of the agreed upon IV. The first
// nonce used is IV+1.
// Check that the IV is not 0, this is likely to indicate
// a bad client. An IV of -1 would be worse, resulting in
// an initial nonce of 0 which can always lead to crypto
// swamps. Why are safe from that, since the Nonce calculation
// always takes the Abs value of the IV, stripping the sign.
func (s *supervisor) checkIV(iv string) error {
	var (
		err       error
		bIV       []byte
		iIV, zero *big.Int
	)
	zero = big.NewInt(0)

	if bIV, err = hex.DecodeString(iv); err != nil {
		return err
	}

	iIV = big.NewInt(0)
	iIV.SetBytes(bIV)
	iIV.Abs(iIV)
	if iIV.Cmp(zero) == 0 {
		return fmt.Errorf(`Invalid Initialization vector`)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
