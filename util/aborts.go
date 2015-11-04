package util

import (
	"fmt"
	"os"
)

func (u *SomaUtil) AbortOnError(err error, txt ...string) {
	if err != nil {
		for _, s := range txt {
			fmt.Fprintf(os.Stderr, "%s\n", s)
			u.Log.Print(s)
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		u.Log.Fatal(err)
	}
}

func (u *SomaUtil) Abort(txt ...string) {
	for _, s := range txt {
		fmt.Fprintf(os.Stderr, "%s\n", s)
		u.Log.Print(s)
	}

	// ensure there is _something_
	if len(txt) == 0 {
		e := `Abort() called without error message. Sorry!`
		fmt.Fprintf(os.Stderr, "%s\n", e)
		u.Log.Print(e)
	}
	os.Exit(1)
}

func (u *SomaUtil) NotImplemented(txt ...string) {
	for _, s := range txt {
		fmt.Fprintf(os.Stderr, "%s\n", s)
		u.Log.Print(s)
	}

	// ensure there is _something_ in the log
	if len(txt) == 0 {
		e := `This command is currently not implemented.`
		fmt.Fprintf(os.Stderr, "%s\n", e)
		u.Log.Print(e)
	}
	os.Exit(1)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
