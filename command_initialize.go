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

    sqlRepositoryTables01()
    log.Print("Installed: Repository01")

    createTableCustomProperties( printOnly, verbose )

    sqlRepositoryTables02()
    log.Print("Installed: Repository02")

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


