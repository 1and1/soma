package util

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
	//"log"
	//"os"
)

func (u *SomaUtil) UserIdByUuidOrName(c *cli.Context) uuid.UUID {
	var (
		id  uuid.UUID
		err error
	)

	switch u.GetCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
		u.AbortOnError(err, "Syntax error, argument not a uuid")
	case 2:
		u.ValidateCliArgument(c, 1, "by-name")
		id = u.GetUserIdByName(c.Args().Get(1))
	default:
		u.Abort("Syntax error, unexpected argument count")
	}
	return id
}

func (u *SomaUtil) GetCliArgumentCount(c *cli.Context) int {
	a := c.Args()
	if !a.Present() {
		return 0
	}
	return len(a.Tail()) + 1
}

func (u *SomaUtil) ValidateCliArgument(c *cli.Context, pos uint8, s string) {
	a := c.Args()
	if a.Get(int(pos)-1) != s {
		u.Abort(fmt.Sprintf("Syntax error, missing keyword: ", s))
	}
}

func (u *SomaUtil) ValidateCliArgumentCount(c *cli.Context, i uint8) {
	a := c.Args()
	if i == 0 {
		if a.Present() {
			u.Abort("Syntax error, command takes no arguments")
		}
	} else {
		if !a.Present() || len(a.Tail()) != (int(i)-1) {
			u.Abort("Syntax error")
		}
	}
}

func (u *SomaUtil) GetFullArgumentSlice(c *cli.Context) []string {
	sl := []string{c.Args().First()}
	sl = append(sl, c.Args().Tail()...)
	return sl
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
