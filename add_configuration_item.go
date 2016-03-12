package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"


	"github.com/julienschmidt/httprouter"
)

func AddConfigurationItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		dec      *json.Decoder
		item     *ConfigurationItem
		lookupID string
		details  *somaproto.DeploymentDetails
		err      error
	)

	dec = json.NewDecoder(r.Body)
	if err = dec.Decode(details); err != nil {
		dispatchUnprocessable(&w, err.Error())
		return
	}

	if lookupID, item, err = Itemize(details); err != nil {
		dispatchUnprocessable(&w, err.Error())
		return
	}

	if err = addItem(item, lookupID); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}
	dispatchNoContent(&w)
}

func addItem(item *ConfigurationItem, lookupID string) error {
	var (
		hostid int64
		err    error
		look   string
		jsonb  []byte
	)

	// string was generated from uint64, we need int now
	if hostid, err = strconv.ParseInt(item.HostId, 10, 64); err != nil {
		return err
	}
	if jsonb, err = json.Marshal(item); err != nil {
		return err
	}

	err = Eye.run.check_lookup.QueryRow(lookupID).Scan(&look)
	if err == sql.ErrNoRows {
		if _, err = Eye.run.insert_lookup.Exec(
			lookupID,
			int(hostid),
			item.Metric,
		); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	if lookupID != look {
		panic(`Database corrupted`)
	}

	_, err = Eye.run.insert_item.Exec(
		item.ConfigurationItemId.String(),
		lookupID,
		jsonb,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
