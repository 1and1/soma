package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func dbOpen() {
	var err error
	driver := "postgres"
	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		Cfg.Database.Name,
		Cfg.Database.User,
		Cfg.Database.Pass,
		Cfg.Database.Host,
		Cfg.Database.Port,
		Cfg.TlsMode,
		Cfg.Timeout,
	)

	db, err = sql.Open(driver, connect)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
