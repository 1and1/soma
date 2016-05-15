package main

func createTablesRoot(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 25)

	queryMap["createTableRootToken"] = `
create table if not exists root.token (
    token                       varchar(256)    NOT NULL
);`
	queries[idx] = "createTableRootToken"
	idx++

	queryMap["createTableRootFlags"] = `
create table if not exists root.flags (
    flag                        varchar(256)    NOT NULL,
    status                      boolean         NOT NULL DEFAULT 'no'
);`
	queries[idx] = "createTableRootFlags"
	idx++

	performDatabaseTask(printOnly, verbose, queries[:idx], queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
