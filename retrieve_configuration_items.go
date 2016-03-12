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
		dispatchInternalServerError(&w, err.Error())
		return
	}

	if jsonb, err = json.Marshal(reply); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}

	dispatchJsonOK(&w, &jsonb)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
