package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ListConfigurationItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		items *sql.Rows
		list  ConfigurationList
		item  string
		err   error
		jsonb []byte
	)

	if items, err = Eye.run.get_items.Query(); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}
	defer items.Close()

	for items.Next() {
		if err = items.Scan(&item); err != nil {
			if err == sql.ErrNoRows {
				dispatchNotFound(&w)
			} else {
				dispatchInternalServerError(&w, err.Error())
			}
			return
		}
		i := item
		list.ConfigurationItemIdList = append(list.ConfigurationItemIdList, i)
	}
	if err = items.Err(); err != nil {
		if err == sql.ErrNoRows {
			dispatchNotFound(&w)
		} else {
			dispatchInternalServerError(&w, err.Error())
		}
		return
	}
	if len(list.ConfigurationItemIdList) == 0 {
		dispatchNotFound(&w)
		return
	}

	if jsonb, err = json.Marshal(list); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}

	dispatchJsonOK(&w, &jsonb)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
