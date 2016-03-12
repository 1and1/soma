package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

func GetConfigurationItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		err     error
		jConfig string
		reply   ConfigurationData
		config  ConfigurationItem
		jsonb   []byte
	)

	if _, err = uuid.FromString(params.ByName("item")); err != nil {
		dispatchBadRequest(&w, err.Error())
		return
	}

	if err = Eye.run.get_config.QueryRow(params.ByName("item")).Scan(&jConfig); err != nil {
		if err == sql.ErrNoRows {
			dispatchNotFound(&w)
		} else {
			dispatchInternalServerError(&w, err.Error())
		}
		return
	}
	reply.Configurations = make([]ConfigurationItem, 1)

	if err = json.Unmarshal([]byte(jConfig), config); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}

	reply.Configurations[0] = config
	if jsonb, err = json.Marshal(reply); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}

	dispatchJsonOK(&w, &jsonb)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
