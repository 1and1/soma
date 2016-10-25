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
	"math/big"
	"time"

	"github.com/1and1/soma/internal/msg"
	"github.com/1and1/soma/internal/stmt"
	"github.com/1and1/soma/lib/auth"
	log "github.com/Sirupsen/logrus"
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
	stmt_CheckUser      *sql.Stmt
	stmt_AddCategory    *sql.Stmt
	stmt_DelCategory    *sql.Stmt
	stmt_ListCategory   *sql.Stmt
	stmt_ShowCategory   *sql.Stmt
	stmt_AddPermission  *sql.Stmt
	stmt_DelPermission  *sql.Stmt
	stmt_ListPermission *sql.Stmt
	stmt_ShowPermission *sql.Stmt
	stmt_SearchPerm     *sql.Stmt
	stmt_GrantSysGlUser *sql.Stmt
	stmt_RevkSysGlUser  *sql.Stmt
	stmt_SrchGlSysGrant *sql.Stmt
	appLog              *log.Logger
	reqLog              *log.Logger
	errLog              *log.Logger
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

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.SelectToken:             s.stmt_FToken,
		stmt.FindUserID:              s.stmt_FindUser,
		stmt.ListPermissionCategory:  s.stmt_ListCategory,
		stmt.ShowPermissionCategory:  s.stmt_ShowCategory,
		stmt.ListPermission:          s.stmt_ListPermission,
		stmt.ShowPermission:          s.stmt_ShowPermission,
		stmt.SearchPermissionByName:  s.stmt_SearchPerm,
		stmt.SearchGlobalSystemGrant: s.stmt_SrchGlSysGrant,
	} {
		if prepStmt, err = s.conn.Prepare(statement); err != nil {
			s.errLog.Fatal(`supervisor`, err, statement)
		}
		defer prepStmt.Close()
	}

	if !s.readonly {
		for statement, prepStmt := range map[string]*sql.Stmt{
			stmt.AddPermissionCategory:        s.stmt_AddCategory,
			stmt.DeletePermissionCategory:     s.stmt_DelCategory,
			stmt.AddPermission:                s.stmt_AddPermission,
			stmt.DeletePermission:             s.stmt_DelPermission,
			stmt.GrantGlobalOrSystemToUser:    s.stmt_GrantSysGlUser,
			stmt.RevokeGlobalOrSystemFromUser: s.stmt_RevkSysGlUser,
			stmt.CheckUserActive:              s.stmt_CheckUser,
		} {
			if prepStmt, err = s.conn.Prepare(statement); err != nil {
				s.errLog.Fatal(`supervisor`, err, statement)
			}
			defer prepStmt.Close()
		}
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
	case `reset_user_password`, `change_user_password`:
		go func() { s.userPassword(q) }()
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

/* Ops Access
 */
func (s *supervisor) shutdownNow() {
	s.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
