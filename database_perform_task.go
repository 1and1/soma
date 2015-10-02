package main

import (
  "fmt"
  "log"
)

func performDatabaseTask(printOnly bool, verbose bool, queries []string, queryMap map[string]string) {
  var err error

  for _, query := range queries {
    // ignore over-allocated slice
    if query == "" {
      continue
    }

    if printOnly || verbose {
      fmt.Printf( "%s\n", queryMap[query] )
      if printOnly {
        continue
      }
    }

    _, err = db.Exec( queryMap[query] ); if err != nil {
      log.Fatal( "Error executing query '", query, "': ", err )
    }
    log.Print( "Executed query: ", query )
  }
}
