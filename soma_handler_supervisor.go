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
	"fmt"
	"log"
	"math/big"
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
	kex         svKexMap
	tokens      svTokenMap
	credentials svCredMap
	stmt_FToken *sql.Stmt
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
		Kex: auth.Kex{
			Public:               kex.Public,
			InitializationVector: kex.InitializationVector,
			Request:              kex.Request,
		},
	}
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

func (s *supervisor) Validate(account, token, addr string) uint16 {
	tok := s.tokens.read(token)
	if tok == nil && !s.readonly {
		// rw instance knows every token
		return 401
	} else if tok == nil {
		if !s.fetchTokenFromDB(token) {
			return 401
		}
		tok = s.tokens.read(token)
	}
	if time.Now().UTC().Before(tok.validFrom.UTC()) ||
		time.Now().UTC().After(tok.expiresAt.UTC()) {
		return 401
	}

	if auth.Verify(account, addr, tok.binToken, s.key,
		s.seed, tok.binExpiresAt, tok.salt) {
		return 200
	}

	return 401
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
