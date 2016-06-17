package cmpl

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func LevelCreate(c *cli.Context) {
	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range []string{`shortname`, `numeric`} {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	hasSHORT := false
	hasNUMERIC := false

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		switch t {
		case `shortname`:
			skip = 1
			hasSHORT = true
			continue
		case `numeric`:
			skip = 1
			hasNUMERIC = true
			continue
		}
	}
	if skip > 0 {
		return
	}
	for _, t := range []string{`shortname`, `numeric`} {
		switch t {
		case `shortname`:
			if !hasSHORT {
				fmt.Println(t)
			}
		case `numeric`:
			if !hasNUMERIC {
				fmt.Println(t)
			}
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
