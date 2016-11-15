/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/auth"
)

// TODO: timer, delete all expired key exchanges
func (s *supervisor) pruneKex() {
	s.kex.lock()
	defer s.kex.unlock()
	for kexId, kex := range s.kex.KMap {
		if kex.IsExpired() {
			delete(s.kex.KMap, kexId)
		}
	}
}

func (s *supervisor) newTokenMap() *svTokenMap {
	m := svTokenMap{}
	m.TMap = make(map[string]svToken)
	return &m
}

func (s *supervisor) newCredentialMap() *svCredMap {
	m := svCredMap{}
	m.CMap = make(map[string]svCredential)
	return &m
}

func (s *supervisor) newKexMap() *svKexMap {
	m := svKexMap{}
	m.KMap = make(map[string]auth.Kex)
	return &m
}

func (s *supervisor) newLockMap() *svLockMap {
	l := svLockMap{}
	l.LockMap = make(map[string]string)
	return &l
}

func (s *supervisor) newGlobalPermMap() *svPermMapGlobal {
	g := svPermMapGlobal{}
	g.GMap = make(map[string]map[string]string)
	return &g
}

func (s *supervisor) newGlobalGrantMap() *svGrantMapGlobal {
	g := svGrantMapGlobal{}
	g.GMap = make(map[string][]string)
	return &g
}

func (s *supervisor) newLimitedPermMap() *svPermMapLimited {
	l := svPermMapLimited{}
	l.LMap = make(map[string]map[string][]string)
	return &l
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

func (s *supervisor) requestLog(q *msg.Request) {
	s.reqLog.Printf(LogStrSRq,
		fmt.Sprintf("supervisor/%s",
			q.Section,
		),
		q.Action,
		q.User,
		q.RemoteAddr,
	)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
