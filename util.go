package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

var Slog *log.Logger

func initLogFile() {
	f, err := os.OpenFile(".somaadm.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	Slog = log.New(f, "somaadm: ", log.Ldate|log.Ltime|log.LUTC)

	Slog.Print("Initialized logfile")
}

func getApiUrl() *url.URL {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		Slog.Printf("Error parsing API address from config file")
		Slog.Fatal(err)
	}
	return url
}

func checkServerKeyword(s string) {
	keywords := []string{"id", "datacenter", "location", "name", "online"}
	for _, k := range keywords {
		if s == k {
			fmt.Fprintf(os.Stderr, "Syntax error: back-to-back keywords")
			os.Exit(1)
		}
	}
}
