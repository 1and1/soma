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
	"log"
	"time"

	"github.com/mjolnir42/scrypth64"
	"github.com/satori/go.uuid"
)

func (s *supervisor) startupLoad() {

	s.startupRoot()

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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
