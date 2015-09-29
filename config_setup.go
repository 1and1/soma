package main

import (
  "log"
  "fmt"
  "strconv"
  "github.com/codegangsta/cli"
  "github.com/howeyc/gopass"
)

func configSetup(c *cli.Context) error {
  log.Print( "Starting runtime config initialization" )

  // try loading a configuration file
  err := Cfg.populateFromFile(c.GlobalString("config"))
  if err != nil {
    if c.GlobalIsSet("config") {
      // missing cli argument file is fatal error
      log.Fatal( err )
    }
    log.Print( "No configuration file found" )
  }

  // finish setting up runtime configuration
  params := []string{"host", "port", "database", "user", "timeout", "tls"}

  for p := range params {
    // update configuration with cli argument overrides
    if c.GlobalIsSet(params[p]) {
      log.Printf("Setting cli value override for %s", params[p])
      switch params[p] {
      case "host":
        Cfg.Database.Host = c.GlobalString(params[p])
      case "port":
        Cfg.Database.Port = strconv.Itoa(c.GlobalInt(params[p]))
      case "database":
        Cfg.Database.Name = c.GlobalString(params[p])
      case "user":
        Cfg.Database.User = c.GlobalString(params[p])
      case "timeout":
        Cfg.Timeout = strconv.Itoa(c.GlobalInt(params[p]))
      case "tls":
        Cfg.TlsMode = c.GlobalString(params[p])
      }
      continue
    }
    // set default values for unset configuration parameters
    switch params[p] {
    case "host":
      if Cfg.Database.Host == "" {
        log.Printf("Setting default value for %s", params[p])
        Cfg.Database.Host = "localhost"
      }
    case "port":
      if Cfg.Database.Port == "" {
        log.Printf("Setting default value for %s", params[p])
        Cfg.Database.Port = strconv.Itoa(5432)
      }
    case "database":
      if Cfg.Database.Name == "" {
        log.Printf("Setting default value for %s", params[p])
        Cfg.Database.Name = "soma"
      }
    case "user":
      if Cfg.Database.User == "" {
        log.Printf("Setting default value for %s", params[p])
        Cfg.Database.User = "soma_dba"
      }
    case "timeout":
      if Cfg.Timeout == "" {
        log.Printf("Setting default value for %s", params[p])
        Cfg.Timeout = strconv.Itoa(3)
      }
    case "tls":
      if Cfg.TlsMode == "" {
        log.Printf("Setting default value for %s", params[p])
        Cfg.TlsMode = "verify-full"
      }
    }
  }

  // prompt for password if the cli flag was set
  if c.GlobalBool("password") {
    fmt.Printf("Enter password: ")
    pass := gopass.GetPasswd()
    Cfg.Database.Pass = string(pass)
  }
  // abort if we have no connection password at this point
  if Cfg.Database.Pass == "" {
    log.Fatal("Can not continue without database connection password")
  }
  return nil
}
