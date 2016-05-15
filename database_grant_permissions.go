package main

func grantPermissions(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 100)

	queryMap["grantServiceUserSchemaSoma"] = `grant select, insert, update, delete on all tables in schema soma to soma_svc;`
	queries[idx] = "grantServiceUserSchemaSoma"
	idx++

	queryMap["grantServiceUserSequencesSoma"] = `grant usage, select on all sequences in schema soma to soma_svc;`
	queries[idx] = "grantServiceUserSequencesSoma"
	idx++

	queryMap["grantServiceUserSchemaRoot"] = `grant select on all tables in schema root ro soma_svc;`
	queries[idx] = "grantServiceUserSequencesSoma"
	idx++

	queryMap["grantServiceUserSchemaInventory"] = `grant select, insert, update, delete on all tables in schema inventory to soma_svc;`
	queries[idx] = "grantServiceUserSchemaInventory"
	idx++

	queryMap["grantServiceUserSequencesInventory"] = `grant usage, select on all sequences in schema inventory to soma_svc;`
	queries[idx] = "grantServiceUserSequencesInventory"
	idx++

	queryMap["grantServiceUserSchemaAuth"] = `grant select, insert, update, delete on all tables in schema auth to soma_svc;`
	queries[idx] = "grantServiceUserSchemaAuth"
	idx++

	queryMap["grantServiceUserSequencesAuth"] = `grant usage, select on all sequences in schema auth to soma_svc;`
	queries[idx] = "grantServiceUserSequencesAuth"
	idx++

	queryMap["grantInventoryUserSchemaInventory"] = `grant select, insert, update, delete on all tables in schema inventory to soma_inv;`
	queries[idx] = "grantInventoryUserSchemaInventory"
	idx++

	queryMap["grantServiceUserSchemaPublic"] = `grant select on all tables in schema public to soma_svc;`
	queries[idx] = "grantServiceUserSchemaPublic"
	idx++

	queryMap["grantServiceUserSequencesPublic"] = `grant usage, select on all sequences in schema public to soma_svc;`
	queries[idx] = "grantServiceUserSequencesPublic"
	idx++

	performDatabaseTask(printOnly, verbose, queries[:idx], queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
