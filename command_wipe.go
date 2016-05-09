package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	db.Exec(`DROP SCHEMA auth CASCADE;`)
	db.Exec(`DROP SCHEMA inventory CASCADE;`)
	db.Exec(`DROP SCHEMA soma CASCADE;`)
	db.Exec(`DROP TABLE public.schema_versions;`)

	done <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
