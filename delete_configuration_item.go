package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

func DeleteConfigurationItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		itemID string
		err    error
	)
	itemID = params.ByName("item")
	if _, err = uuid.FromString(itemID); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = deleteItem(itemID); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}

func deleteItem(itemID string) error {
	var (
		delete_item, get_lookup, item_count, delete_lookup *sql.Stmt
		lookupID                                           string
		count                                              int
		err                                                error
	)

	if delete_item, err = Eye.conn.Prepare(stmtDeleteConfigurationItem); err != nil {
		return err
	}
	defer delete_item.Close()

	if delete_lookup, err = Eye.conn.Prepare(stmtDeleteLookupId); err != nil {
		return err
	}
	defer delete_lookup.Close()

	if get_lookup, err = Eye.conn.Prepare(stmtGetLookupIdForItem); err != nil {
		return err
	}
	defer get_lookup.Close()

	if item_count, err = Eye.conn.Prepare(stmtGetItemCountForLookupId); err != nil {
		return err
	}
	defer item_count.Close()

	if err = get_lookup.QueryRow(itemID).Scan(&lookupID); err != nil {
		// either a real error, or what is to be deleted does not exist
		return err
	}

	if _, err = delete_item.Exec(itemID); err != nil {
		return err
	}

	if err = item_count.QueryRow(lookupID).Scan(&count); err != nil {
		return err
	}

	if count != 0 {
		return nil
	}
	_, err = delete_lookup.Exec(lookupID)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
