package main

func createTablesSchemaVersion(printOnly bool, verbose bool) {

	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 3)

	queryMap["createTableSchemaVersionSOMA"] = `
create table if not exists public.schema_versions (
    serial                      bigserial       PRIMARY KEY,
    schema                      varchar(16)     NOT NULL,
    version                     numeric(16,0)   NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    description                 text            NOT NULL
);`
	queries[idx] = "createTableSchemaVersionSOMA"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
