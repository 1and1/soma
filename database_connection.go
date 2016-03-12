package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func connectToDatabase() {
	var err error
	driver := "postgres"

	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		Eye.Database.Name,
		Eye.Database.User,
		Eye.Database.Pass,
		Eye.Database.Host,
		Eye.Database.Port,
		Eye.TlsMode,
		Eye.Timeout,
	)

	Eye.run.conn, err = sql.Open(driver, connect)
	if err != nil {
		log.Fatal(err)
	}
	if err = Eye.run.conn.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Print("Connected to database")
	if _, err = Eye.run.conn.Exec(`SET TIME ZONE 'UTC';`); err != nil {
		log.Fatal(err)
	}
}

func pingDatabase() {
	ticker := time.NewTicker(time.Second).C

	for {
		<-ticker
		err := Eye.run.conn.Ping()
		if err != nil {
			log.Print(err)
		}
	}
}

func prepareStatements() {
	var err error

	Eye.run.check_item, err = Eye.run.conn.Prepare(stmtCheckItemExists)
	abortOnError(err)

	Eye.run.check_lookup, err = Eye.run.conn.Prepare(stmtCheckLookupExists)
	abortOnError(err)

	Eye.run.delete_item, err = Eye.run.conn.Prepare(stmtDeleteConfigurationItem)
	abortOnError(err)

	Eye.run.delete_lookup, err = Eye.run.conn.Prepare(stmtDeleteLookupId)
	abortOnError(err)

	Eye.run.get_config, err = Eye.run.conn.Prepare(stmtGetSingleConfiguration)
	abortOnError(err)

	Eye.run.get_items, err = Eye.run.conn.Prepare(stmtGetConfigurationItemIds)
	abortOnError(err)

	Eye.run.get_lookup, err = Eye.run.conn.Prepare(stmtGetLookupIdForItem)
	abortOnError(err)

	Eye.run.insert_item, err = Eye.run.conn.Prepare(stmtInsertConfigurationItem)
	abortOnError(err)

	Eye.run.insert_lookup, err = Eye.run.conn.Prepare(stmtInsertLookupInformation)
	abortOnError(err)

	Eye.run.item_count, err = Eye.run.conn.Prepare(stmtGetItemCountForLookupId)
	abortOnError(err)

	Eye.run.retrieve, err = Eye.run.conn.Prepare(stmtRetrieveConfigurationsByLookup)
	abortOnError(err)

	Eye.run.update_item, err = Eye.run.conn.Prepare(stmtUpdateConfigurationItem)
	abortOnError(err)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
