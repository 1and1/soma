package main

import (
  "log"
  "fmt"
)

func createSqlSchema(printOnly bool) {
  var err error
  idx := 0
  // map for storing the SQL statements by name
  queryMap := make( map[string]string )
  // slice storing the required statement order so foreign keys can
  // resolve successfully
  queries := make( []string, 5 )


  queryMap["createSchemaSoma"] = `
create schema if not exists soma;`
  queries[idx] = "createSchemaSoma"; idx++


  queryMap["createSchemaInventory"] = `create schema if not exists inventory;`
  queries[idx] = "createSchemaInventory"; idx++


  queryMap["createSchemaAuth"] = `create schema if not exists auth;`
  queries[idx] = "createSchemaAuth"; idx++


  queryMap["alterDatabaseDefaultSearchPath"] = fmt.Sprintf("alter database %s set search_path TO soma,inventory,auth;", Cfg.Database.Name)
  queries[idx] = "alterDatabaseDefaultSearchPath"; idx++


  queryMap["alterDatabaseSearchPath"] = `set search_path to soma,inventory,auth;`
  queries[idx] = "alterDatabaseSearchPath"; idx++


  for _, query := range queries {
    // ignore over-allocated slice
    if query == "" {
      continue
    }

    if printOnly {
      log.Print( queryMap[query] )
      continue
    }

    _, err = db.Exec( queryMap[query] ); if err != nil {
      log.Fatal( "Error executing query '", query, "': ", err )
    }
    log.Print( "Executed query: ", query )
  }
}
