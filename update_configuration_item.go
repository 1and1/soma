package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func UpdateConfigurationItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

	if lookupID != params.ByName("item") {
		http.Error(w, "Mismatching ConfigurationItemID", http.StatusBadRequest)
		return
	}

	if err = updateItem(item, lookupID); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

func updateItem(item *ConfigurationItem, lookupID string) error {
	var (
		check, update_item *sql.Stmt
		itemID             string
		err                error
		jsonb              []byte
	)

	if check, err = Eye.conn.Prepare(stmtCheckItemExists); err != nil {
		return err
	}
	if update_item, err = Eye.conn.Prepare(stmtUpdateConfigurationItem); err != nil {
		return err
	}

	if jsonb, err = json.Marshal(item); err != nil {
		return err
	}

	// since this was an explicit update request, non-existance is a
	// hard error
	if err = check.QueryRow(item.ConfigurationItemId.String()).Scan(&itemID); err != nil {
		return err
	}
	if itemID != item.ConfigurationItemId.String() {
		panic(`Database corrupted`)
	}

	_, err = update_item.Exec(
		item.ConfigurationItemId.String(),
		lookupID,
		jsonb,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
