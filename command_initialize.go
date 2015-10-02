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

    sqlPropertyTables01()
    log.Print("Installed: Property01")

    sqlRepositoryTables01()
    log.Print("Installed: Repository01")

    sqlPropertyTables02()
    log.Print("Installed: Property02")

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


