package util

import (
	"fmt"

	"github.com/codegangsta/cli"
)

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

func (u *SomaUtil) ValidateCliMinArgumentCount(c *cli.Context, i uint8) {
	ct := u.GetCliArgumentCount(c)
	if ct < int(i) {
		u.Abort("Syntax error, incorrect argument count")
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
			u.Abort("Syntax error, incorrect argument count")
		}
	}
}

func (u *SomaUtil) GetFullArgumentSlice(c *cli.Context) []string {
	sl := []string{c.Args().First()}
	sl = append(sl, c.Args().Tail()...)
	return sl
}

func (u *SomaUtil) ParseVariableArguments(keys []string, rKeys []string, args []string) (map[string]string, []string) {
	// return map of the parse result
	result := make(map[string]string)
	// map to test which required keys were found
	argumentCheck := make(map[string]bool)
	// return slice which optional keys were found
	optionalKeys := make([]string, 0)
	// no required keys is valid
	if len(rKeys) > 0 {
		for _, key := range rKeys {
			argumentCheck[key] = false
		}
	}
	skipNext := false

	for pos, val := range args {
		// skip current argument if last argument was a keyword
		if skipNext {
			skipNext = false
			continue
		}

		if u.SliceContainsString(val, keys) {
			// check back-to-back keywords
			u.CheckStringNotAKeyword(args[pos+1], keys)
			result[val] = args[pos+1]
			argumentCheck[val] = true
			skipNext = true
			if !u.SliceContainsString(val, rKeys) {
				optionalKeys = append(optionalKeys, val)
			}
			continue
		}
		// keywords trigger continue, arguments are skipped over.
		// reaching this is an error
		u.Abort(fmt.Sprintf("Syntax error, erroneus argument: %s", val))
	}

	// check we managed to collect all required keywords
	for _, v := range argumentCheck {
		if !v {
			u.Abort(fmt.Sprintf("Syntax error, missing keyword: %s", v))
		}
	}

	return result, optionalKeys
}

/*
 * This function parses whitespace separated argument lists of
 * keyword/value pairs were keywords can be specified multiple
 * times, some keywords are required and some only allowed once.
 * Sequence of multiple keywords are detected and lead to abort
 *
 * multKeys => [ "port", "transport" ]
 * uniqKeys => [ "team" ]
 * reqKeys  => [ "team" ]
 * args     => [ "port", "53", "transport", "tcp", "transport",
 *               "udp", "team", "ITOMI" ]
 *
 * result => result["team"] = [ "ITOMI" ]
 *           result["port"] = [ "53" ]
 *           result["transport"] = [ "tcp", "udp" ]
 */
func (u *SomaUtil) ParseVariadicArguments(
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

		if u.SliceContainsString(val, keys) {
			// check for back-to-back keyswords
			u.CheckStringNotAKeyword(args[pos+1], keys)

			// append value of current keyword into result map
			result[val] = append(result[val], args[pos+1])
			skip = true
			continue
		}
		// keywords trigger continue before this
		// values after keywords are skip'ed
		// reaching this is an error
		u.Abort(fmt.Sprintf("Syntax error, erroneus argument: %s", val))
	}

	// check if we managed to collect all required keywords
	for _, key := range reqKeys {
		_, ok := result[key] // ok is false if slice is nil
		if !ok {
			u.Abort(fmt.Sprintf("Syntax error, missing keyword: %s", key))
		}
	}

	// check if unique keywords were only specified once
	for _, key := range uniqKeys {
		sl, ok := result[key]
		if ok && (len(sl) > 1) {
			u.Abort(fmt.Sprintf("Syntax error, keyword must only be provided once: %s", key))
		}
	}

	return result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
