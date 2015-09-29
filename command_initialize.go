package main

import (
)

func commandInitialize(done chan<- bool) {
    dbOpen()
    sqlSchema()
    sqlMetaTables01()
    sqlInventoryTables01()
    sqlMetaTables02()
    sqlInventoryTables02()
    sqlAuthTables01()
    // root_token table
    // user_keys
    // user_certificates
    // tool_keys
    // tool_certificate
    done <- true
}


