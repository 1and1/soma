/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package adm

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

// ParseVariadicArguments parses split up argument lists of
// keyword/value pairs were keywords can be specified multiple
// times, some keywords are required and some only allowed once.
// Sequence of multiple keywords are detected and lead to abort
//
//	multKeys => [ "port", "transport" ]
//	uniqKeys => [ "team" ]
//	reqKeys  => [ "team" ]
//	args     => [ "port", "53", "transport", "tcp", "transport",
//	              "udp", "team", "GenericOps" ]
//
//	result => result["team"] = [ "GenericOps" ]
//	          result["port"] = [ "53" ]
//	          result["transport"] = [ "tcp", "udp" ]
func ParseVariadicArguments(
	result map[string][]string, // provided result map
	multKeys []string, // keys that may appear multiple times
	uniqKeys []string, // keys that are allowed at most once
	reqKeys []string, // keys that are required at least one
	args []string, // arguments to parse
) error {
	// ensure the result is empty
	result = make(map[string][]string)
	// used to hold found errors, so if three keywords are missing they can
	// all be mentioned in one call
	errors := []string{}

	// merge key slices
	keys := append(multKeys, uniqKeys...)

	// helper to skip over next value in args slice
	skip := false

	for pos, val := range args {
		// skip current arg if last argument was a keyword
		if skip {
			skip = false
			continue
		}

		if sliceContainsString(val, keys) {
			// there must be at least one arguments left
			if len(args[pos+1:]) < 1 {
				errors = append(errors,
					`Syntax error, incomplete key/value specification (too few items left to parse)`,
				)
				goto abort
			}
			// check for back-to-back keyswords
			if err := checkStringNotAKeyword(args[pos+1], keys); err != nil {
				errors = append(errors, err.Error())
				goto abort
			}

			// append value of current keyword into result map
			result[val] = append(result[val], args[pos+1])
			skip = true
			continue
		}
		// keywords trigger continue before this
		// values after keywords are skip'ed
		// reaching this is an error
		errors = append(errors, fmt.Sprintf("Syntax error, erroneus argument: %s", val))
	}

	// check if we managed to collect all required keywords
	for _, key := range reqKeys {
		// ok is false if slice is nil
		if _, ok := result[key]; !ok {
			errors = append(errors, fmt.Sprintf("Syntax error, missing keyword: %s", key))
		}
	}

	// check if unique keywords were only specified once
	for _, key := range uniqKeys {
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			errors = append(errors, fmt.Sprintf("Syntax error, keyword must only be provided once: %s", key))
		}
	}

abort:
	if len(errors) > 0 {
		result = nil
		return fmt.Errorf(combineStrings(errors...))
	}

	return nil
}

// VerifySingleArgument takes a context and verifies there is only one
// commandline argument
func VerifySingleArgument(c *cli.Context) error {
	a := c.Args()
	if !a.Present() {
		return fmt.Errorf(`Syntax error, command requires argument`)
	}

	if len(a.Tail()) != 0 {
		return fmt.Errorf(
			"Syntax error, too many arguments (expected: 1, received %d)",
			len(a.Tail())+1,
		)
	}
	return nil
}

// VerifyNoArgument takes a context and verifies there is no
// commandline argument
func VerifyNoArgument(c *cli.Context) error {
	a := c.Args()
	if a.Present() {
		return fmt.Errorf(`Syntax error, command takes no arguments`)
	}

	return nil
}

func AllArguments(c *cli.Context) []string {
	sl := []string{c.Args().First()}
	sl = append(sl, c.Args().Tail()...)
	return sl
}

func sliceContainsString(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

func checkStringNotAKeyword(s string, keys []string) error {
	if sliceContainsString(s, keys) {
		return fmt.Errorf("Syntax error, back-to-back keyword: %s", s)
	}
	return nil
}

// combineStrings takes an arbitray number of strings and combines them
// into one, separated by `.\n`
func combineStrings(s ...string) string {
	var out string
	spacer := ``
	for _, in := range s {
		// ensure a single trailing .
		out = fmt.Sprintf("%s%s", out+spacer, strings.TrimRight(in, `.`)+`.`)
		spacer = "\n"
	}
	return out
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
