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
	"time"

	"github.com/1and1/soma/internal/stmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/scrypth64"
	uuid "github.com/satori/go.uuid"
)

func (s *supervisor) startupLoad() {

	s.startupRoot()

	s.startupUsersAndTeams()
	s.startupPermissions()

	if !s.readonly {
		s.startupCredentials()
	}

	s.startupTokens()

	s.startupGrants()
}

func (s *supervisor) startupRoot() {
	var (
		err                  error
		flag, crypt          string
		mcf                  scrypth64.Mcf
		validFrom, expiresAt time.Time
		state                bool
		rows                 *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadRootFlags)
	if err != nil {
		log.Fatal(`supervisor/load-root-flags,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&flag,
			&state,
		); err != nil {
			log.Fatal(`supervisor/load-root-flags,scan: `, err)
		}
		switch flag {
		case `disabled`:
			s.root_disabled = state
		case `restricted`:
			s.root_restricted = state
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(`supervisor/load-root-flags,next: `, err)
	}

	// only load root credentials on master instance
	if !s.readonly {
		if err = s.conn.QueryRow(stmt.LoadRootPassword).Scan(
			&crypt,
			&validFrom,
			&expiresAt,
		); err == sql.ErrNoRows {
			// root bootstrap outstanding
			return
		} else if err != nil {
			log.Fatal(`supervisor/load-root-password: `, err)
		}
		if mcf, err = scrypth64.FromString(crypt); err != nil {
			log.Fatal(`supervisor/string-to-mcf: `, err)
		}
		s.credentials.insert(`root`, uuid.Nil, validFrom.UTC(),
			PosTimeInf.UTC(), mcf)
	}
}

func (s *supervisor) startupCredentials() {
	var (
		err                  error
		rows                 *sql.Rows
		user_id, user, crypt string
		reset                bool
		validFrom, expiresAt time.Time
		id                   uuid.UUID
		mcf                  scrypth64.Mcf
	)

	rows, err = s.conn.Query(stmt.LoadAllUserCredentials)
	if err != nil {
		log.Fatal(`supervisor/load-credentials,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&user_id,
			&crypt,
			&reset,
			&validFrom,
			&expiresAt,
			&user,
		); err != nil {
			log.Fatal(`supervisor/load-credentials,scan: `, err)
		}

		if id, err = uuid.FromString(user_id); err != nil {
			log.Fatal(`supervisor/string-to-uuid: `, err)
		}
		if mcf, err = scrypth64.FromString(crypt); err != nil {
			log.Fatal(`supervisor/string-to-mcf: `, err)
		}

		s.credentials.restore(user, id, validFrom, expiresAt, mcf, reset, true)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(`supervisor/load-credentials,next: `, err)
	}
}

func (s *supervisor) startupTokens() {
	var (
		err                         error
		token, salt, valid, expires string
		validFrom, expiresAt        time.Time
		rows                        *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadAllTokens)
	if err != nil {
		log.Fatal(`supervisor/load-tokens,query: `, err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&token,
			&salt,
			&validFrom,
			&expiresAt,
		); err != nil {
			log.Fatal(`supervisor/load-tokens,scan: `, err)
		}
		valid = validFrom.Format(rfc3339Milli)
		expires = expiresAt.Format(rfc3339Milli)

		if err = s.tokens.insert(token, valid, expires, salt); err != nil {
			log.Fatal(`supervisor/load-tokens,insert: `, err)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(`supervisor/load-tokens,next: `, err)
	}
}

func (s *supervisor) startupUsersAndTeams() {
	var (
		err                                    error
		userUUID, userName, teamUUID, teamName string
		rows                                   *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadUserTeamMapping)
	if err != nil {
		log.Fatal(`supervisor/load-user-team-mapping,query: `, err)
	}
	defer rows.Close()

	// reduce lock overhead by locking here once and then using the
	// unlocked bulk interface
	s.id_user.lock()
	defer s.id_user.unlock()
	s.id_user_rev.lock()
	defer s.id_user_rev.unlock()
	s.id_team.lock()
	defer s.id_team.unlock()
	s.id_userteam.lock()
	defer s.id_userteam.unlock()

	for rows.Next() {
		if err = rows.Scan(
			&userUUID,
			&userName,
			&teamUUID,
			&teamName,
		); err != nil {
			log.Fatal(`supervisor/load-user-team-mapping,scan: `, err)
		}
		s.id_user.load(userUUID, userName)
		s.id_user_rev.load(userName, userUUID)
		s.id_team.load(teamUUID, teamName)
		s.id_userteam.load(userUUID, teamUUID)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(`supervisor/load-user-team-mapping,next: `, err)
	}
}

func (s *supervisor) startupPermissions() {
	var (
		err                error
		permUUID, permName string
		rows               *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadPermissions)
	if err != nil {
		log.Fatal(`supervisor/load-permissions,query: `, err)
	}
	defer rows.Close()

	// reduce lock overhead by locking here once and then using the
	// unlocked bulk interface
	s.id_permission.lock()
	defer s.id_permission.unlock()

	for rows.Next() {
		if err = rows.Scan(
			&permUUID,
			&permName,
		); err != nil {
			log.Fatal(`supervisor/load-permissions,scan: `, err)
		}
		s.id_permission.load(permName, permUUID)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(`supervisor/load-permissions,next: `, err)
	}
}

func (s *supervisor) startupGrants() {
	var (
		err                           error
		grantUUID, permUUID, userUUID string
		rows                          *sql.Rows
	)

	rows, err = s.conn.Query(stmt.LoadGlobalOrSystemUserGrants)
	if err != nil {
		log.Fatal(`supervisor/load-user-systemglobal-grants,query: `, err)
	}
	defer rows.Close()

	// reduce lock overhead by locking here once and then using the
	// unlocked load method
	s.global_permissions.lock()
	defer s.global_permissions.unlock()
	s.global_grants.lock()
	defer s.global_grants.unlock()

	for rows.Next() {
		if err = rows.Scan(
			&grantUUID,
			&userUUID,
			&permUUID,
		); err != nil {
			log.Fatal(`supervisor/load-user-systemglobal-grants,scan: `, err)
		}
		s.global_permissions.load(userUUID, permUUID, grantUUID)
		s.global_grants.load(userUUID, permUUID, grantUUID)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(`supervisor/load-user-systemglobal-grants,next: `, err)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
