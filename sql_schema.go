package main

import (
  "log"
  "fmt"
)

func sqlSchema() {
  var err error;
  var query string;

  _, err = db.Exec("CREATE SCHEMA IF NOT EXISTS soma;"); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec("CREATE SCHEMA IF NOT EXISTS inventory;"); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec("CREATE SCHEMA IF NOT EXISTS auth;"); if err != nil {
    log.Fatal( err )
  }

  query = fmt.Sprintf("ALTER DATABASE %s SET search_path TO soma,inventory,auth;", Cfg.Database.Name)
  _, err = db.Exec(query); if err != nil {
    log.Fatal( err )
  }

  _, err = db.Exec("SET search_path TO soma,inventory,auth;"); if err != nil {
    log.Fatal( err )
  }
}
