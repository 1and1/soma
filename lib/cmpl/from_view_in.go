package cmpl

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func FromViewIn(c *cli.Context) {
	switch {
	case c.NArg() == 0:
		return
	case c.NArg() == 1:
		for _, t := range []string{`from`, `view`, `in`} {
			fmt.Println(t)
		}
		return
	}

	skip := 0
	hasFROM := false
	hasVIEW := false
	hasIN := false

	for _, t := range c.Args().Tail() {
		if skip > 0 {
			skip--
			continue
		}
		switch t {
		case `from`:
			skip = 1
			hasFROM = true
			continue
		case `view`:
			skip = 1
			hasVIEW = true
			continue
		case `in`:
			skip = 1
			hasIN = true
			continue
		}
	}
	if skip > 0 {
		return
	}
	for _, t := range []string{`from`, `view`, `in`} {
		switch t {
		case `from`:
			if !hasFROM {
				fmt.Println(t)
			}
		case `view`:
			if !hasVIEW {
				fmt.Println(t)
			}
		case `in`:
			if !hasIN {
				fmt.Println(t)
			}
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
