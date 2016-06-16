package main

func createTablesJobs(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableJobStatus"] = `
create table if not exists soma.job_status (
    job_status                  varchar(32)     PRIMARY KEY
);`
	queries[idx] = "createTableJobStatus"
	idx++

	queryMap["createTableJobResults"] = `
create table if not exists soma.job_results (
    job_result                  varchar(32)     PRIMARY KEY
);`
	queries[idx] = "createTableJobResults"
	idx++

	queryMap["createTableJobTypes"] = `
create table if not exists soma.job_types (
    job_type                    varchar(128)    PRIMARY KEY
);`
	queries[idx] = "createTableJobTypes"
	idx++

	queryMap["createTableJobs"] = `
create table if not exists soma.jobs (
    job_id                      uuid            PRIMARY KEY,
    job_status                  varchar(32)     NOT NULL REFERENCES soma.job_status ( job_status ) DEFERRABLE,
    job_result                  varchar(32)     NOT NULL REFERENCES soma.job_results ( job_result ) DEFERRABLE,
    job_type                    varchar(128)    NOT NULL REFERENCES soma.job_types ( job_type ) DEFERRABLE,
    job_serial                  bigserial       NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    user_id                     uuid            NOT NULL REFERENCES inventory.users ( user_id ) DEFERRABLE,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    job_error                   text            NOT NULL DEFAULT '',
    job_queued                  timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    job_started                 timestamptz(3),
    job_finished                timestamptz(3),
    job                         jsonb           NOT NULL
);`
	queries[idx] = "createTableJobs"
	idx++

	queryMap["createIndexJobStatus"] = `
create index _not_processed_jobs
    on soma.jobs ( organizational_team_id, user_id, job_status, job_id )
    where job_status != 'processed'
;`
	queries[idx] = "createIndexJobStatus"
	idx++

	queryMap["createIndexRepoJobs"] = `
create index _jobs_by_repo
    on soma.jobs ( repository_id, job_serial, job_id, job_status )
;`
	queries[idx] = "createIndexRepoJobs"
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
