package main

import (
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
		dispatchBadRequest(&w, err.Error())
		return
	}
	if err = deleteItem(itemID); err != nil {
		dispatchInternalServerError(&w, err.Error())
		return
	}
	dispatchNoContent(&w)
}

func deleteItem(itemID string) error {
	var (
		lookupID string
		count    int
		err      error
	)

	if err = Eye.run.get_lookup.QueryRow(itemID).Scan(&lookupID); err != nil {
		// either a real error, or what is to be deleted does not exist
		return err
	}

	if _, err = Eye.run.delete_item.Exec(itemID); err != nil {
		return err
	}

	if err = Eye.run.item_count.QueryRow(lookupID).Scan(&count); err != nil {
		return err
	}

	if count != 0 {
		return nil
	}
	_, err = Eye.run.delete_lookup.Exec(lookupID)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
