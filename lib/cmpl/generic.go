package cmpl

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func Generic(c *cli.Context, keywords []string) {
	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range keywords {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	match := make(map[string]bool)

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		skip = 1
		match[t] = true
		continue
	}
	// do not complete in porisitons where arguments are expected
	if skip > 0 {
		return
	}
	for _, t := range keywords {
		if !match[t] {
			fmt.Println(t)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
