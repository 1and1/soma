package main

func createTablesPropertyInstances(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTablePropertyInstances"] = `
create table if not exists soma.property_instances (
  instance_id                 uuid            PRIMARY KEY,
  repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
  bucket_id                   uuid            NOT NULL REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
  source_instance_id          uuid            NOT NULL,
  source_object_type          varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
  source_object_id            uuid            NOT NULL,
  UNIQUE ( instance_id, repository_id ),
  FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE,
  FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTablePropertyInstances"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
