package main

import (
	"database/sql"
	"fmt"

)

func CheckUpdateOrInsert(details *somaproto.DeploymentDetails) error {
	var (
		stmt             *sql.Stmt
		err              error
		itemID, lookupID string
		item             *ConfigurationItem
	)
	if stmt, err = Eye.conn.Prepare(stmtCheckItemExists); err != nil {
		return err
	}

	if lookupID, item, err = Itemize(details); err != nil {
		return err
	}

	fmt.Println(lookupID)
	fmt.Println(item)

	err = stmt.QueryRow(item.ConfigurationItemId).Scan(&itemID)
	if err == sql.ErrNoRows {
		return addItem(item, lookupID)
	} else if err != nil {
		return err
	}
	// UPDATE
	if item.ConfigurationItemId.String() != itemID {
		panic(`Database corrupted`)
	}

	return updateItem(item, lookupID)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
