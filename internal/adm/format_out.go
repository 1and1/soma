package adm

import (
	"fmt"

	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
)

func FormatOut(c *cli.Context, resp *resty.Response, cmd string) error {
	if c.GlobalBool(`json`) {
		fmt.Println(resp)
		return nil
	}

	// hardwire JSON output for now
	fmt.Println(resp)
	return nil

	switch cmd {
	case `list`:
	case `show`:
	case `tree`:
	default:
	}

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
