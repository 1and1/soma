package main

import (
  "log"
)

func commandInitialize(done chan<- bool) {
    dbOpen()
    sqlSchema()
    log.Print("Installed: Schema")
    sqlMetaTables01()
    log.Print("Installed: Meta01")
    sqlInventoryTables01()
    log.Print("Installed: Inventory01")
    sqlMetaTables02()
    log.Print("Installed: Meta02")
    sqlInventoryTables02()
    log.Print("Installed: Inventory02")
    sqlAuthTables01()
    log.Print("Installed: Auth01")
    // root_token table
    // user_keys
    // user_certificates
    // tool_keys
    // tool_certificate
    sqlPropertyTables01()
    log.Print("Installed: Property01")
    sqlRepositoryTables01()
    log.Print("Installed: Repository01")
    sqlPropertyTables02()
    log.Print("Installed: Property02")
    sqlRepositoryTables02()
    log.Print("Installed: Repository02")
    sqlBucketsTables01()
    log.Print("Installed: Buckets01")
    sqlNodeTables01()
    log.Print("Installed: Node01")
    sqlClusterTables01()
    log.Print("Installed: Cluster01")
    sqlGroupTables01()
    log.Print("Installed: Group01")
    sqlPermissionTables01()
    log.Print("Installed: Permission01")
    createTablesMetricsMonitoring( false )
    log.Print("Installed: MetricsMonitoring")
    done <- true
}


