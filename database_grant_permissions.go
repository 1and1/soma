package main

func grantPermissions(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 4)

	queryMap["grantServiceUserSchemaSoma"] = `grant select, insert, update, delete on all tables in schema soma to soma_svc;`
	queries[idx] = "grantServiceUserSchemaSoma"
	idx++

	queryMap["grantServiceUserSchemaInventory"] = `grant select, insert, update, delete on all tables in schema inventory to soma_svc;`
	queries[idx] = "grantServiceUserSchemaInventory"
	idx++

	queryMap["grantServiceUserSchemaAuth"] = `grant select, insert, update, delete on all tables in schema auth to soma_svc;`
	queries[idx] = "grantServiceUserSchemaAuth"
	idx++

	queryMap["grantInventoryUserSchemaInventory"] = `grant select, insert, update, delete on all tables in schema inventory to soma_inv;`
	queries[idx] = "grantInventoryUserSchemaInventory"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
