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
	"strconv"
	"strings"

	"github.com/1and1/soma/lib/proto"
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

// ParseVariadicCheckArguments is a version of ParseVariadicArguments
// that can handle the threshold and constraint keywords for
// checks, which consist of multiple key/value pairs. Keywords do not
// have to be passed in.
func ParseVariadicCheckArguments(
	result map[string][]string,
	constraints []proto.CheckConfigConstraint,
	thresholds []proto.CheckConfigThreshold,
	args []string,
) error {
	// used to hold found errors, so if three keywords are missing they can
	// all be mentioned in one call
	errors := []string{}

	multiple := []string{
		`threshold`,
		`constraint`}
	unique := []string{
		`in`,
		`on`,
		`with`,
		`interval`,
		`inheritance`,
		`childrenonly`,
		`extern`}
	required := []string{
		`in`,
		`on`,
		`with`,
		`interval`}

	// merge key slices
	keys := append(multiple, unique...)

	// iteration helper
	skip := false
	skipcount := 0

argloop:
	for pos, val := range args {
		// skip current arg if it was already consumed
		if skip {
			skipcount--
			if skipcount == 0 {
				skip = false
			}
			continue
		}

		if sliceContainsString(val, keys) {
			// there must be at least one arguments left
			if len(args[pos+1:]) < 1 {
				errors = append(errors, `Syntax error, incomplete`+
					` key/value specification (too few items left`+
					` to parse)`,
				)
				goto abort
			}
			// check for back-to-back keyswords
			if err := checkStringNotAKeyword(
				args[pos+1], keys,
			); err != nil {
				errors = append(errors, err.Error())
				goto abort
			}

			switch val {
			case `threshold`:
				if len(args[pos+1:]) < 6 {
					errors = append(errors, `Syntax error, incomplete`+
						`threshold specification`)
					goto abort
				}
				thr := proto.CheckConfigThreshold{}
				if err := parseThresholdChain(
					thr,
					args[pos+1:pos+7],
				); err != nil {
					errors = append(errors, err.Error())
				} else {
					thresholds = append(thresholds, thr)
				}
				skip = true
				skipcount = 6
				continue argloop

			case `constraint`:
				// argument is the start of a constraint specification.
				// check we have enough arguments left
				if len(args[pos+1:]) < 3 {
					errors = append(errors, `Syntax error, incomplete`+
						` constraint specification`)
				}
				constr := proto.CheckConfigConstraint{}
				if err := parseConstraintChain(
					constr,
					args[pos+1:pos+3],
				); err != nil {
					errors = append(errors, err.Error())
				} else {
					constraints = append(constraints, constr)
				}
				skip = true
				skipcount = 3
				continue argloop

			case `on`:
				result[`on/type`] = append(result[`on/type`],
					args[pos+1])
				result[`on/object`] = append(result[`on/object`],
					args[pos+2])
				// set for required+unique checks
				result[val] = append(result[val], fmt.Sprintf(
					"%s::%s", args[pos+1], args[pos+2]))
				skip = true
				skipcount = 2
				continue argloop

			default:
				// regular key/value keyword
				result[val] = append(result[val], args[pos+1])
				skip = true
				skipcount = 1
				continue argloop
			}
		}
		// error is reached if argument was not skipped and not a
		// recognized keyword
		errors = append(errors, fmt.Sprintf("Syntax error, erroneus"+
			" argument: %s", val))
	}

	// check if all required keywords were collected
	for _, key := range required {
		if _, ok := result[key]; !ok {
			errors = append(errors, fmt.Sprintf("Syntax error,"+
				" missing keyword: %s", key))
		}
	}

	// check if unique keywords were only specuified once
	for _, key := range unique {
		// check ok since unique may still be optional
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			errors = append(errors, fmt.Sprintf("Syntax error,"+
				" keyword must only be provided once: %s", key))
		}
	}

abort:
	if len(errors) > 0 {
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

// AllArguments returns all arguments from the given cli.Context
func AllArguments(c *cli.Context) []string {
	sl := []string{c.Args().First()}
	sl = append(sl, c.Args().Tail()...)
	return sl
}

// sliceContainsString checks whether string s is in slice sl
func sliceContainsString(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

// checkStringNotAKeyword checks whether string s in not in slice keys
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

// parseThresholdChain parses a single threshold specification given
// to ParseVariadicCheckArguments
func parseThresholdChain(result proto.CheckConfigThreshold,
	args []string) error {
	tParse := make(map[string][]string)
	if err := ParseVariadicArguments(
		tParse,
		[]string{},
		[]string{`predicate`, `level`, `value`},
		[]string{`predicate`, `level`, `value`},
		args,
	); err != nil {
		return err
	}
	var err error
	result.Predicate.Symbol = tParse[`predicate`][0]
	result.Level.Name = tParse[`level`][0]
	if result.Value, err = strconv.ParseInt(
		tParse[`value`][0], 10, 64,
	); err != nil {
		return fmt.Errorf("Syntax error, value argument not"+
			" numeric: %s", tParse[`value`][0])
	}
	return nil
}

// parseConstraintChain parses a single constraint specification
// given to ParseVariadicCheckArguments
func parseConstraintChain(result proto.CheckConfigConstraint,
	args []string) error {
	result.ConstraintType = args[0]
	switch result.ConstraintType {
	case `service`:
		result.Service = &proto.PropertyService{}
		switch args[1] {
		case `name`:
			result.Service.Name = args[2]
		default:
			return fmt.Errorf("Syntax error, can not constraint"+
				" service to %s", args[1])
		}

	case `oncall`:
		result.Oncall = &proto.PropertyOncall{}
		switch args[1] {
		case `id`:
			result.Oncall.Id = args[2]
		case `name`:
			result.Oncall.Name = args[2]
		default:
			return fmt.Errorf("Syntax error, can not constraint"+
				" oncall to %s", args[1])
		}

	case `attribute`:
		result.Attribute = &proto.ServiceAttribute{
			Name:  args[1],
			Value: args[2],
		}

	case `system`:
		result.System = &proto.PropertySystem{
			Name:  args[1],
			Value: args[2],
		}

	case `custom`:
		result.Custom = &proto.PropertyCustom{
			Name:  args[1],
			Value: args[2],
		}

	case `native`:
		result.Native = &proto.PropertyNative{
			Name:  args[1],
			Value: args[2],
		}

	default:
		return fmt.Errorf("Syntax error, unknown constraint type: %s",
			args[0])
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
