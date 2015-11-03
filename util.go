package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
)

var Slog *log.Logger

func initLogFile() {
	f, err := os.OpenFile(path.Join(Cfg.Run.PathLogs, "somaadm.log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logging: %s\n", err.Error())
		os.Exit(1)
	}

	utl.SetLog(log.New(f, "", log.Ldate|log.Ltime|log.LUTC))
	// XXX COMPAT
	Cfg.Run.Logger = utl.Log
	Slog = utl.Log
}

func getApiUrl() *url.URL {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		utl.Log.Printf("Error parsing API address from config file")
		utl.Log.Fatal(err)
	}
	return url
}

func parseLimitedGrantArguments(keys []string, args []string) *map[string]string {
	result := make(map[string]string)
	argumentCheck := make(map[string]bool)
	for _, key := range keys {
		argumentCheck[key] = false
	}

	skipNext := false

	for pos, val := range args {
		if skipNext {
			skipNext = false
			continue
		}
		switch val {
		case "repository":
			utl.CheckStringNotAKeyword(args[pos+1], keys)
			result["repository"] = args[pos+1]
			skipNext = true
			argumentCheck["repository"] = true
		case "bucket":
			utl.CheckStringNotAKeyword(args[pos+1], keys)
			result["bucket"] = args[pos+1]
			skipNext = true
			argumentCheck["bucket"] = true
		case "group":
			utl.CheckStringNotAKeyword(args[pos+1], keys)
			result["group"] = args[pos+1]
			skipNext = true
			argumentCheck["group"] = true
		case "cluster":
			utl.CheckStringNotAKeyword(args[pos+1], keys)
			result["group"] = args[pos+1]
			skipNext = true
			argumentCheck["cluster"] = true
		}
	}

	// check we managed to collect all required keywords
	for k, v := range argumentCheck {
		if !v {
			utl.Log.Fatal("Syntax error, missing keyword for argument count: ", k)
		}
	}

	return &result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
