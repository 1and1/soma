package main

func commandInitialize(done chan<- bool, printOnly bool, verbose bool) {
    dbOpen()

    createSqlSchema( printOnly, verbose )

    createTablesMetaData( printOnly, verbose )

    createTablesInventoryAssets( printOnly, verbose )

    createTablesDatacenterMetaData( printOnly, verbose )

    createTablesInventoryAccounts( printOnly, verbose )

    createTablesAuthentication( printOnly, verbose )
    // root_token table

    createTablesProperties( printOnly, verbose )

    createTableRepositories( printOnly, verbose )

    createTableCustomProperties( printOnly, verbose )

    createTablesRepositoryProperties( printOnly, verbose )

    createTablesBuckets( printOnly, verbose )

    createTablesNodes( printOnly, verbose )

    createTablesClusters( printOnly, verbose )

    createTablesGroups( printOnly, verbose )

    createTablesPermissions( printOnly, verbose )

    createTablesMetricsMonitoring( printOnly, verbose )

    done <- true
}


