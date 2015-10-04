package main

import (
  "log"
)

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

    sqlClusterTables01()
    log.Print("Installed: Cluster01")

    sqlGroupTables01()
    log.Print("Installed: Group01")

    sqlPermissionTables01()
    log.Print("Installed: Permission01")

    createTablesMetricsMonitoring( printOnly, verbose )

    done <- true
}


