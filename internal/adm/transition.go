/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2015-2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"fmt"
	"os"
)

func Abort(txt ...string) {
	for _, s := range txt {
		fmt.Fprintf(os.Stderr, "%s\n", s)
	}

	// ensure there is _something_
	if len(txt) == 0 {
		e := `Abort() called without error message. Sorry!`
		fmt.Fprintf(os.Stderr, "%s\n", e)
	}
	os.Exit(1)
}

func AbortOnError(err error, txt ...string) {
	if err != nil {
		for _, s := range txt {
			fmt.Fprintf(os.Stderr, "%s\n", s)
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
