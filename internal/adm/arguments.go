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
	"log"
)

// ParseVariadicArguments parses whitespace separated argument lists
// of keyword/value pairs were keywords can be specified multiple
// times, some keywords are required and some only allowed once.
// Sequence of multiple keywords are detected and lead to abort
//
// multKeys => [ "port", "transport" ]
// uniqKeys => [ "team" ]
// reqKeys  => [ "team" ]
// args     => [ "port", "53", "transport", "tcp", "transport",
//               "udp", "team", "GenericOps" ]
//
// result => result["team"] = [ "GenericOps" ]
//           result["port"] = [ "53" ]
//           result["transport"] = [ "tcp", "udp" ]
func ParseVariadicArguments(
	multKeys []string, // keys that may appear multiple times
	uniqKeys []string, // keys that are allowed at most once
	reqKeys []string, // keys that are required at least one
	args []string, // arguments to parse
) map[string][]string {
	// returns a map of slices of string
	result := make(map[string][]string)

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

		if SliceContainsString(val, keys) {
			// there must be at least one arguments left
			if len(args[pos+1:]) < 1 {
				Abort("Syntax error, incomplete key/value specification (too few items left to parse)")
			}
			// check for back-to-back keyswords
			CheckStringNotAKeyword(args[pos+1], keys)

			// append value of current keyword into result map
			result[val] = append(result[val], args[pos+1])
			skip = true
			continue
		}
		// keywords trigger continue before this
		// values after keywords are skip'ed
		// reaching this is an error
		Abort(fmt.Sprintf("Syntax error, erroneus argument: %s", val))
	}

	// check if we managed to collect all required keywords
	for _, key := range reqKeys {
		// ok is false if slice is nil
		if _, ok := result[key]; !ok {
			Abort(fmt.Sprintf("Syntax error, missing keyword: %s", key))
		}
	}

	// check if unique keywords were only specified once
	for _, key := range uniqKeys {
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			Abort(fmt.Sprintf("Syntax error, keyword must only be provided once: %s", key))
		}
	}

	return result
}

func SliceContainsString(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

func CheckStringNotAKeyword(s string, keys []string) {
	if SliceContainsString(s, keys) {
		log.Fatal(`Syntax error, back-to-back keywords`) // XXX
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
