package main

import (
	"database/sql"
	"fmt"

)

func CheckUpdateOrInsertOrDelete(details *somaproto.DeploymentDetails) error {
	var (
		stmt             *sql.Stmt
		err              error
		itemID, lookupID string
		item             *ConfigurationItem
	)
	if stmt, err = Eye.conn.Prepare(stmtCheckItemExists); err != nil {
		return err
	}
	defer stmt.Close()

	if lookupID, item, err = Itemize(details); err != nil {
		return err
	}

	fmt.Println(lookupID)
	fmt.Println(item)

	err = stmt.QueryRow(item.ConfigurationItemId).Scan(&itemID)
	switch details.Task {
	case "rollout":
		if err == sql.ErrNoRows {
			return addItem(item, lookupID)
		} else if err != nil {
			return err
		}
	case "deprovision":
		if err != nil {
			return err
		}
	}

	if item.ConfigurationItemId.String() != itemID {
		panic(`Database corrupted`)
	}
	switch details.Task {
	case "rollout":
		return updateItem(item, lookupID)
	case "deprovision":
		return deleteItem(itemID)
	default:
		return fmt.Errorf(`Unknown Task requested`)
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
