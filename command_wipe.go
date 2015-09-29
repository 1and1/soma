package main

import (
  "log"
  "bufio"
  "os"
  "fmt"
  "strings"
)

func commandWipe(done chan<- bool, forced bool) {
  if !forced {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("Are you sure (yes/no)? ")
    text, _ := reader.ReadString('\n')
    text = strings.TrimSpace(text)
    if text != "yes" {
      os.Exit(0)
    }
  }
  log.Printf("Wiping database %s", Cfg.Database.Name)

  dbOpen()

  db.Exec(`DROP SCHEMA auth CASCADE;`);
  db.Exec(`DROP SCHEMA inventory CASCADE;`)
  db.Exec(`DROP SCHEMA soma CASCADE;`);

  done <- true
}
