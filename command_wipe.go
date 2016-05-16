package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func commandWipe(done chan<- bool, forced, printOnly bool) {
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
	stmts := []string{
		`DROP SCHEMA IF EXISTS auth CASCADE;`,
		`DROP SCHEMA IF EXISTS inventory CASCADE;`,
		`DROP SCHEMA IF EXISTS soma CASCADE;`,
		`DROP SCHEMA IF EXISTS root CASCADE;`,
		`DROP TABLE IF EXISTS public.schema_versions;`,
	}

	if !printOnly {
		dbOpen()

		for _, stmt := range stmts {
			db.Exec(stmt)
		}
	} else {
		for _, stmt := range stmts {
			fmt.Println(stmt)
		}
	}

	done <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
