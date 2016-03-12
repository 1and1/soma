package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ListConfigurationItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		get_items *sql.Stmt
		items     *sql.Rows
		list      ConfigurationList
		item      string
		err       error
		jsonb     []byte
	)

	if get_items, err = Eye.conn.Prepare(stmtGetConfigurationItemIds); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if items, err = get_items.Query(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for items.Next() {
		if err = items.Scan(&item); err != nil {
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
		i := item
		list.ConfigurationItemIdList = append(list.ConfigurationItemIdList, i)
	}
	if len(list.ConfigurationItemIdList) == 0 {
		log.Println(err)
		http.Error(w, "No items found", http.StatusNotFound)
		return
	}

	if jsonb, err = json.Marshal(list); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonb)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
