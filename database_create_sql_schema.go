package main

import (
  "fmt"
)

func createSqlSchema(printOnly bool, verbose bool) {
  idx := 0
  // map for storing the SQL statements by name
  queryMap := make( map[string]string )
  // slice storing the required statement order so foreign keys can
  // resolve successfully
  queries := make( []string, 10 )


  queryMap["createSchemaSoma"] = `create schema if not exists soma;`
  queries[idx] = "createSchemaSoma"; idx++


  queryMap["createSchemaInventory"] = `create schema if not exists inventory;`
  queries[idx] = "createSchemaInventory"; idx++


  queryMap["createSchemaAuth"] = `create schema if not exists auth;`
  queries[idx] = "createSchemaAuth"; idx++


  queryMap["alterDatabaseDefaultSearchPath"] = fmt.Sprintf("alter database %s set search_path TO soma,inventory,auth;", Cfg.Database.Name)
  queries[idx] = "alterDatabaseDefaultSearchPath"; idx++


  queryMap["alterDatabaseSearchPath"] = `set search_path to soma,inventory,auth;`
  queries[idx] = "alterDatabaseSearchPath"; idx++


  queryMap["grantServiceUserSchemaSoma"] = `grant select, insert, update, delete on all tables in schema soma to soma_svc;`
  queries[idx] = "grantServiceUserSchemaSoma"; idx++


  queryMap["grantServiceUserSchemaInventory"] = `grant select, insert, update, delete on all tables in schema inventory to soma_svc;`
  queries[idx] = "grantServiceUserSchemaInventory"; idx++


  queryMap["grantServiceUserSchemaAuth"] = `grant select, insert, update, delete on all tables in schema auth to soma_svc;`
  queries[idx] = "grantServiceUserSchemaAuth"; idx++


  performDatabaseTask( printOnly, verbose, queries, queryMap )
}
