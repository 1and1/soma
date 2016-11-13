/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main

import (
	"database/sql"

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
	kex                 *svKexMap
	tokens              *svTokenMap
	credentials         *svCredMap
	global_permissions  *svPermMapGlobal
	global_grants       *svGrantMapGlobal
	limited_permissions *svPermMapLimited
	id_user             *svLockMap
	id_user_rev         *svLockMap
	id_team             *svLockMap
	id_permission       *svLockMap
	id_userteam         *svLockMap
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
	stmt_SectionList    *sql.Stmt
	stmt_SectionShow    *sql.Stmt
	stmt_SectionSearch  *sql.Stmt
	stmt_SectionAdd     *sql.Stmt
	stmt_ActionList     *sql.Stmt
	stmt_ActionShow     *sql.Stmt
	stmt_ActionSearch   *sql.Stmt
	stmt_ActionAdd      *sql.Stmt
	stmt_RevokeGlobal   *sql.Stmt
	stmt_RevokeRepo     *sql.Stmt
	stmt_RevokeTeam     *sql.Stmt
	stmt_RevokeMonitor  *sql.Stmt
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
		stmt.SectionList:             s.stmt_SectionList,
		stmt.SectionShow:             s.stmt_SectionShow,
		stmt.SectionSearch:           s.stmt_SectionSearch,
		stmt.ActionList:              s.stmt_ActionList,
		stmt.ActionShow:              s.stmt_ActionShow,
		stmt.ActionSearch:            s.stmt_ActionSearch,
	} {
		if prepStmt, err = s.conn.Prepare(statement); err != nil {
			s.errLog.Fatal(`supervisor`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

	if !s.readonly {
		for statement, prepStmt := range map[string]*sql.Stmt{
			stmt.AddPermissionCategory:         s.stmt_AddCategory,
			stmt.DeletePermissionCategory:      s.stmt_DelCategory,
			stmt.AddPermission:                 s.stmt_AddPermission,
			stmt.DeletePermission:              s.stmt_DelPermission,
			stmt.GrantGlobalOrSystemToUser:     s.stmt_GrantSysGlUser,
			stmt.RevokeGlobalOrSystemFromUser:  s.stmt_RevkSysGlUser,
			stmt.CheckUserActive:               s.stmt_CheckUser,
			stmt.SectionAdd:                    s.stmt_SectionAdd,
			stmt.ActionAdd:                     s.stmt_ActionAdd,
			stmt.RevokeGlobalAuthorization:     s.stmt_RevokeGlobal,
			stmt.RevokeRepositoryAuthorization: s.stmt_RevokeRepo,
			stmt.RevokeTeamAuthorization:       s.stmt_RevokeTeam,
			stmt.RevokeMonitoringAuthorization: s.stmt_RevokeMonitor,
		} {
			if prepStmt, err = s.conn.Prepare(statement); err != nil {
				s.errLog.Fatal(`supervisor`, err, stmt.Name(statement))
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
	switch q.Section {
	case `kex`:
		go func() { s.kexInit(q) }()

	case `bootstrap`:
		s.bootstrapRoot(q)

	case `authenticate`:
		go func() { s.validate_basic_auth(q) }()

	case `token`:
		go func() { s.issue_token(q) }()

	case `activate`:
		go func() { s.activate_user(q) }()

	case `password`:
		go func() { s.userPassword(q) }()

	case `authorize`:
		go func() { s.authorize(q) }()

	case `map`:
		go func() { s.update_map(q) }()

	case `category`:
		switch q.Action {
		case `add`, `remove`:
			s.permission_category(q)
		default:
			go func() { s.permission_category(q) }()
		}

	case `permission`:
		switch q.Action {
		case `add`, `delete`:
			s.permission(q)
		default:
			go func() { s.permission(q) }()
		}

	case `right`:
		s.right(q)

	case `section`:
		s.section(q)

	case `action`:
		s.action(q)
	}
}

func (s *supervisor) shutdownNow() {
	s.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
