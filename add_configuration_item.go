package main

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddConfigurationItem(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func addItem(item ConfigurationItem, lookupID string) error {
	var (
		stmt *sql.Stmt
		err  error
	)
	if stmt, err = Eye.conn.Prepare(stmtCheckLookupExists); err != nil {
		return err
	}
	err = stmt.QueryRow(lookupID).Scan()
	if err == sql.ErrNoRows {
		return err
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
