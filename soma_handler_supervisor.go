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
	input               chan msg.Request
	shutdown            chan bool
	conn                *sql.DB
	seed                []byte
	key                 []byte
	readonly            bool
	tokenExpiry         uint64
	kexExpiry           uint64
	credExpiry          uint64
	activation          string
	root_disabled       bool
	root_restricted     bool
	kex                 svKexMap
	tokens              svTokenMap
	credentials         svCredMap
	global_permissions  svPermMapGlobal
	global_grants       svGrantMapGlobal
	limited_permissions svPermMapLimited
	id_user             svLockMap
	id_user_rev         svLockMap
	id_team             svLockMap
	id_permission       svLockMap
	id_userteam         svLockMap
	stmt_FToken         *sql.Stmt
	stmt_FindUser       *sql.Stmt
	stmt_AddCategory    *sql.Stmt
	stmt_DelCategory    *sql.Stmt
	stmt_ListCategory   *sql.Stmt
	stmt_ShowCategory   *sql.Stmt
	stmt_AddPermission  *sql.Stmt
	stmt_DelPermission  *sql.Stmt
	stmt_ListPermission *sql.Stmt
	stmt_ShowPermission *sql.Stmt
}

func (s *supervisor) run() {
	var err error

	// set library options
	auth.TokenExpirySeconds = s.tokenExpiry
	auth.KexExpirySeconds = s.kexExpiry

	// initialize maps
	s.id_user = s.newLockMap()
	s.id_user_rev = s.newLockMap()
	s.id_team = s.newLockMap()
	s.id_permission = s.newLockMap()
	s.id_userteam = s.newLockMap()
	s.tokens = s.newTokenMap()
	s.credentials = s.newCredentialMap()
	s.kex = s.newKexMap()
	s.global_permissions = s.newGlobalPermMap()
	s.global_grants = s.newGlobalGrantMap()
	s.limited_permissions = s.newLimitedPermMap()

	// load from datbase
	s.startupLoad()

	// prepare SQL statements
	if s.stmt_FToken, err = s.conn.Prepare(stmt.SelectToken); err != nil {
		log.Fatal("supervisor/fetch-token: ", err)
	}
	defer s.stmt_FToken.Close()

	if s.stmt_FindUser, err = s.conn.Prepare(stmt.FindUserID); err != nil {
		log.Fatal(`supervisor/find-userid: `, err)
	}
	defer s.stmt_FindUser.Close()

	if s.stmt_ListCategory, err = s.conn.Prepare(stmt.ListPermissionCategory); err != nil {
		log.Fatal(`supervisor/list-category: `, err)
	}
	defer s.stmt_ListCategory.Close()

	if s.stmt_ShowCategory, err = s.conn.Prepare(stmt.ShowPermissionCategory); err != nil {
		log.Fatal(`supervisor/show-category: `, err)
	}
	defer s.stmt_ShowCategory.Close()

	if s.stmt_ListPermission, err = s.conn.Prepare(stmt.ListPermission); err != nil {
		log.Fatal(`supervisor/list-permission: `, err)
	}
	defer s.stmt_ListPermission.Close()

	if s.stmt_ShowPermission, err = s.conn.Prepare(stmt.ShowPermission); err != nil {
		log.Fatal(`supervisor/show-permission: `, err)
	}
	defer s.stmt_ShowPermission.Close()

	if !s.readonly {
		if s.stmt_AddCategory, err = s.conn.Prepare(stmt.AddPermissionCategory); err != nil {
			log.Fatal(`supervisor/add-category: `, err)
		}
		defer s.stmt_AddCategory.Close()

		if s.stmt_DelCategory, err = s.conn.Prepare(stmt.DeletePermissionCategory); err != nil {
			log.Fatal(`supervisor/delete-category: `, err)
		}
		defer s.stmt_DelCategory.Close()

		if s.stmt_AddPermission, err = s.conn.Prepare(stmt.AddPermission); err != nil {
			log.Fatal(`supervisor/add-permission: `, err)
		}
		defer s.stmt_AddPermission.Close()

		if s.stmt_DelPermission, err = s.conn.Prepare(stmt.DeletePermission); err != nil {
			log.Fatal(`supervisor/delete-permission: `, err)
		}
		defer s.stmt_DelPermission.Close()
	}

runloop:
	for {
		select {
		case <-s.shutdown:
			break runloop
		case req := <-s.input:
			s.process(&req)
		}
	}
}

func (s *supervisor) process(q *msg.Request) {
	switch q.Action {
	case `kex_init`:
		go func() { s.kexInit(q) }()
	case `bootstrap_root`:
		s.bootstrapRoot(q)
	case `basic_auth`:
		go func() { s.validate_basic_auth(q) }()
	case `request_token`:
		go func() { s.issue_token(q) }()
	case `activate_user`:
		go func() { s.activate_user(q) }()
	case `authorize`:
		go func() { s.authorize(q) }()
	case `update_map`:
		go func() { s.update_map(q) }()
	case `category`:
		if q.Super.Action == `add` || q.Super.Action == `delete` {
			s.permission_category(q)
		} else {
			go func() { s.permission_category(q) }()
		}
	case `permission`:
		if q.Super.Action == `add` || q.Super.Action == `delete` {
			s.permission(q)
		} else {
			go func() { s.permission(q) }()
		}
	case `right`:
		if q.Super.Action == `grant` || q.Super.Action == `revoke` {
			s.right(q)
		} else {
			go func() { s.right(q) }()
		}
	}
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

func (s *supervisor) newLockMap() svLockMap {
	l := svLockMap{}
	l.LockMap = make(map[string]string)
	return l
}

func (s *supervisor) newGlobalPermMap() svPermMapGlobal {
	g := svPermMapGlobal{}
	g.GMap = make(map[string]map[string]string)
	return g
}

func (s *supervisor) newGlobalGrantMap() svGrantMapGlobal {
	g := svGrantMapGlobal{}
	g.GMap = make(map[string][]string)
	return g
}

func (s *supervisor) newLimitedPermMap() svPermMapLimited {
	l := svPermMapLimited{}
	l.LMap = make(map[string]map[string][]string)
	return l
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
