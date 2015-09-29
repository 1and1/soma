package main

import (
  "fmt"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
)

func dbOpen() {
  var err error
  driver := "postgres"
  connect := fmt.Sprintf( "dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
    Cfg.Database.Name,
    Cfg.Database.User,
    Cfg.Database.Pass,
    Cfg.Database.Host,
    Cfg.Database.Port,
    Cfg.TlsMode,
    Cfg.Timeout,
  )

  db, err = sql.Open( driver, connect )
  if err != nil {
    log.Fatal( err )
  }
  if err = db.Ping(); err != nil {
    log.Fatal( err )
  }
}
