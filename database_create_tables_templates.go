package main

func createTablesTemplates(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableTemplates"] = `
create table if not exists soma.templates (
  template_id                 uuid            PRIMARY KEY,
  repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ),
  template_name               varchar(128)    NOT NULL,
  UNIQUE( repository_id, template_name )
);`
	queries[idx] = "createTableTemplates"
	idx++

	queryMap["createTableTemplateAssignments"] = `
create table if not exists soma.template_assignments (
  template_id                 uuid            NOT NULL REFERENCES soma.templates ( template_id ),
  configuration_object        uuid            NOT NULL,
  configuration_object_type   varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ),
  inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
  children_on                 boolean         NOT NULL DEFAULT 'no',
  UNIQUE ( template_id, configuration_object ),
  CHECK ( configuration_object_type != 'server' ),
  CHECK ( configuration_object_type != 'template' )
);`
	queries[idx] = "createTableTemplateAssignments"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
