package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

func GetConfigurationItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		err        error
		get_config *sql.Stmt
		jConfig    string
		reply      ConfigurationData
		config     ConfigurationItem
		jsonb      []byte
	)

	if _, err = uuid.FromString(params.ByName("item")); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if get_config, err = Eye.conn.Prepare(stmtGetSingleConfiguration); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer get_config.Close()

	if err = get_config.QueryRow(params.ByName("item")).Scan(&jConfig); err != nil {
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
	reply.Configurations = make([]ConfigurationItem, 1)

	if err = json.Unmarshal([]byte(jConfig), config); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reply.Configurations[0] = config
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
