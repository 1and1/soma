package main

func commandInitialize(done chan<- bool, printOnly bool, verbose bool, version string) {
	if !printOnly {
		dbOpen()
	}

	createSqlSchema(printOnly, verbose)

	createTablesMetaData(printOnly, verbose)

	createTablesInventoryAssets(printOnly, verbose)

	createTablesDatacenterMetaData(printOnly, verbose)

	createTablesInventoryAccounts(printOnly, verbose)

	createTablesAuthentication(printOnly, verbose)

	createTablesRoot(printOnly, verbose)

	createTablesProperties(printOnly, verbose)

	createTableRepositories(printOnly, verbose)

	createTablesBuckets(printOnly, verbose)

	createTablesPropertyInstances(printOnly, verbose)

	createTableCustomProperties(printOnly, verbose)

	createTablesRepositoryProperties(printOnly, verbose)

	createTablesBucketsProperties(printOnly, verbose)

	createTablesNodes(printOnly, verbose)

	createTablesClusters(printOnly, verbose)

	createTablesGroups(printOnly, verbose)

	createTablesPermissions(printOnly, verbose)

	createTablesMetricsMonitoring(printOnly, verbose)

	createTablesChecks(printOnly, verbose)

	createTablesTemplates(printOnly, verbose)

	createTablesInstances(printOnly, verbose)

	createTablesJobs(printOnly, verbose)

	createTablesSchemaVersion(printOnly, verbose)

	schemaInserts(printOnly, verbose)

	schemaVersionInserts(printOnly, verbose, version)

	grantPermissions(printOnly, verbose)

	insertRootToken(printOnly, verbose)

	done <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
