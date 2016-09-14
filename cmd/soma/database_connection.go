package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

func connectToDatabase() {
	var err error
	var rows *sql.Rows
	var schema string
	var schemaVer int64

	driver := "postgres"

	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		SomaCfg.Database.Name,
		SomaCfg.Database.User,
		SomaCfg.Database.Pass,
		SomaCfg.Database.Host,
		SomaCfg.Database.Port,
		SomaCfg.Database.TlsMode,
		SomaCfg.Database.Timeout,
	)

	// enable handling of infinity timestamps
	pq.EnableInfinityTs(NegTimeInf, PosTimeInf)

	conn, err = sql.Open(driver, connect)
	if err != nil {
		log.Fatal(err)
	}
	if err = conn.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Print("Connected main pool to database")
	if _, err = conn.Exec(`SET TIME ZONE 'UTC';`); err != nil {
		log.Fatal(err)
	}
	if _, err = conn.Exec(`SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE;`); err != nil {
		log.Fatal(err)
	}

	// size the connection pool
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(15)
	conn.SetConnMaxLifetime(12 * time.Hour)

	// required schema versions
	required := map[string]int64{
		"inventory": 201605060001,
		"root":      201605160001,
		"auth":      201605190001,
		"soma":      201609120001,
	}

	if rows, err = conn.Query(`
SELECT schema,
       MAX(version) AS version
FROM   public.schema_versions
GROUP  BY schema;`); err != nil {
		log.Fatal("Query db schema versions: ", err)
	}

	for rows.Next() {
		if err = rows.Scan(
			&schema,
			&schemaVer,
		); err != nil {
			log.Fatal("Schema check: ", err)
		}
		if rsv, ok := required[schema]; ok {
			if rsv != schemaVer {
				log.Fatalf("Incompatible schema %s: %d != %d", schema, rsv, schemaVer)
			} else {
				log.Printf("DB Schema %s, version: %d", schema, schemaVer)
				delete(required, schema)
			}
		} else {
			log.Fatal("Unknown schema: ", schema)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal("Schema check: ", err)
	}
	if len(required) != 0 {
		for s, _ := range required {
			log.Printf("Missing schema: %s", s)
		}
		log.Fatal("FATAL - database incomplete")
	}
}

func newDatabaseConnection() (*sql.DB, error) {
	driver := "postgres"

	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		SomaCfg.Database.Name,
		SomaCfg.Database.User,
		SomaCfg.Database.Pass,
		SomaCfg.Database.Host,
		SomaCfg.Database.Port,
		SomaCfg.Database.TlsMode,
		SomaCfg.Database.Timeout,
	)

	dbcon, err := sql.Open(driver, connect)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		return nil, err
	}
	if _, err = conn.Exec(`SET TIME ZONE 'UTC';`); err != nil {
		return nil, err
	}
	if _, err = conn.Exec(`SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE;`); err != nil {
		log.Fatal(err)
	}
	dbcon.SetMaxIdleConns(1)
	dbcon.SetMaxOpenConns(5)
	dbcon.SetConnMaxLifetime(12 * time.Hour)
	log.Print("Connected new secondary pool to database")
	return dbcon, nil
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
