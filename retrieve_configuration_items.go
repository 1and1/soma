package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RetrieveConfigurationItems(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		err            error
		reply          ConfigurationData
		jsonb          []byte
		lookup, config string
		retrieve       *sql.Stmt
		rows           *sql.Rows
	)

	lookup = params.ByName("lookup")
	// lookup is supposed to be a sha256 hash
	if len(lookup) != 64 {
		http.Error(w, "Invalid lookup id format", http.StatusBadRequest)
		return
	}

	reply.Configurations = []ConfigurationItem{}

	if retrieve, err = Eye.run.conn.Prepare(stmtRetrieveConfigurationsByLookup); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer retrieve.Close()

	if rows, err = retrieve.Query(lookup); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			http.Error(w, "No items found", http.StatusNotFound)
			return
		} else {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&config); err != nil {
			if err == sql.ErrNoRows {
				log.Println(err)
				http.Error(w, "No items found", http.StatusNotFound)
				return
			} else {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		c := ConfigurationItem{}
		if err = json.Unmarshal([]byte(config), c); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reply.Configurations = append(reply.Configurations, c)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if jsonb, err = json.Marshal(reply); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonb)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
