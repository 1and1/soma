package main

import (
	"database/sql"
	"encoding/json"
	"log"
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
		log.Println(err)
		http.Error(w, err.Error(), 422)
		return
	}

	if lookupID, item, err = Itemize(details); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 422)
		return
	}

	if err = addItem(item, lookupID); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

func addItem(item *ConfigurationItem, lookupID string) error {
	var (
		check, insert_lookup, insert_item *sql.Stmt
		hostid                            int64
		err                               error
		look                              string
		jsonb                             []byte
	)

	if check, err = Eye.conn.Prepare(stmtCheckLookupExists); err != nil {
		return err
	}
	defer check.Close()

	if insert_lookup, err = Eye.conn.Prepare(stmtInsertLookupInformation); err != nil {
		return err
	}
	defer insert_lookup.Close()

	if insert_item, err = Eye.conn.Prepare(stmtInsertConfigurationItem); err != nil {
		return err
	}
	defer insert_item.Close()

	// string was generated from uint64, we need int now
	if hostid, err = strconv.ParseInt(item.HostId, 10, 64); err != nil {
		return err
	}
	if jsonb, err = json.Marshal(item); err != nil {
		return err
	}

	err = check.QueryRow(lookupID).Scan(&look)
	if err == sql.ErrNoRows {
		if _, err = insert_lookup.Exec(
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

	_, err = insert_item.Exec(
		item.ConfigurationItemId.String(),
		lookupID,
		jsonb,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
