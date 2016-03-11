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
		Eye.Database.Name,
		Eye.Database.User,
		Eye.Database.Pass,
		Eye.Database.Host,
		Eye.Database.Port,
		Eye.TlsMode,
		Eye.Timeout,
	)

	Eye.conn, err = sql.Open(driver, connect)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected to database")
	Eye.conn.Exec(`SET TIME ZONE 'UTC';`)
}

func pingDatabase() {
	ticker := time.NewTicker(time.Second).C

	for {
		<-ticker
		err := Eye.conn.Ping()
		if err != nil {
			log.Print(err)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
