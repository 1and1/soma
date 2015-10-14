package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func connectToDatabase() {
	var err error
	driver := "postgres"

	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		SomaCfg.Database.Name,
		SomaCfg.Database.User,
		SomaCfg.Database.Pass,
		SomaCfg.Database.Host,
		SomaCfg.Database.Port,
		SomaCfg.TlsMode,
		SomaCfg.Timeout,
	)

	conn, err = sql.Open(driver, connect)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected to database")
}

func pingDatabase() {
	ticker := time.NewTicker(time.Second).C

	for {
		<-ticker
		err := conn.Ping()
		if err != nil {
			log.Print(err)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
