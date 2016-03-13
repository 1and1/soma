/*
Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
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
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RetrieveConfigurationItems(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		err            error
		reply          ConfigurationData
		jsonb          []byte
		lookup, config string
		rows           *sql.Rows
	)

	lookup = params.ByName("lookup")
	// lookup is supposed to be a sha256 hash
	if len(lookup) != 64 {
		dispatchBadRequest(&w, "Invalid lookup id format")
		return
	}

	reply.Configurations = []ConfigurationItem{}

	if rows, err = Eye.run.retrieve.Query(lookup); err != nil {
		if err == sql.ErrNoRows {
			dispatchNotFound(&w)
		} else {
			dispatchInternalServerError(&w, err.Error())
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&config); err != nil {
			if err == sql.ErrNoRows {
				dispatchNotFound(&w)
			} else {
				dispatchInternalServerError(&w, err.Error())
			}
			return
		}
		c := ConfigurationItem{}
		if err = json.Unmarshal([]byte(config), c); err != nil {
			dispatchInternalServerError(&w, err.Error())
			return
		}
		reply.Configurations = append(reply.Configurations, c)
	}
	if err = rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			dispatchNotFound(&w)
		} else {
			dispatchInternalServerError(&w, err.Error())
		}
		return
	}

	if jsonb, err = json.Marshal(reply); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}

	dispatchJsonOK(&w, &jsonb)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
