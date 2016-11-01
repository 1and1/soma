package util

import (
	"fmt"
	"strconv"

	"gopkg.in/resty.v0"

	proto "github.com/1and1/soma/lib/proto"

	"github.com/codegangsta/cli"
)

func (u *SomaUtil) getCliArgumentCount(c *cli.Context) int {
	a := c.Args()
	if !a.Present() {
		return 0
	}
	return len(a.Tail()) + 1
}

func (u *SomaUtil) validateCliArgument(c *cli.Context, pos uint8, s string) {
	a := c.Args()
	if a.Get(int(pos)-1) != s {
		u.abort(fmt.Sprintf("Syntax error, missing keyword: %s", s))
	}
}

func (u *SomaUtil) validateCliMinArgumentCount(c *cli.Context, i uint8) {
	ct := u.getCliArgumentCount(c)
	if ct < int(i) {
		u.abort(fmt.Sprintf(
			"Syntax error, incorrect argument count (%d < %d+ expected)",
			ct,
			i,
		))
	}
}

func (u *SomaUtil) validateCliArgumentCount(c *cli.Context, i uint8) {
	a := c.Args()
	if i == 0 {
		if a.Present() {
			u.abort("Syntax error, command takes no arguments")
		}
	} else {
		if !a.Present() || len(a.Tail()) != (int(i)-1) {
			u.abort(fmt.Sprintf(
				"Syntax error, incorrect argument count (expected: %d, received %d)",
				i,
				len(a.Tail()),
			))
		}
	}
}

func (u *SomaUtil) GetFullArgumentSlice(c *cli.Context) []string {
	sl := []string{c.Args().First()}
	sl = append(sl, c.Args().Tail()...)
	return sl
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
func (u *SomaUtil) parseVariadicArguments(
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
			// there must be at least one arguments left
			if len(args[pos+1:]) < 1 {
				u.abort("Syntax error, incomplete key/value specification (too few items left to parse)")
			}
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
		u.abort(fmt.Sprintf("Syntax error, erroneus argument: %s", val))
	}

	// check if we managed to collect all required keywords
	for _, key := range reqKeys {
		// ok is false if slice is nil
		if _, ok := result[key]; !ok {
			u.abort(fmt.Sprintf("Syntax error, missing keyword: %s", key))
		}
	}

	// check if unique keywords were only specified once
	for _, key := range uniqKeys {
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			u.abort(fmt.Sprintf("Syntax error, keyword must only be provided once: %s", key))
		}
	}

	return result
}

func (u *SomaUtil) ParseVariadicCheckArguments(args []string) (
	map[string][]string,
	[]proto.CheckConfigConstraint,
	[]proto.CheckConfigThreshold,
) {
	// create return objects
	result := make(map[string][]string)
	constraints := []proto.CheckConfigConstraint{}
	thresholds := []proto.CheckConfigThreshold{}
	var err error

	multiple := []string{"threshold", "constraint"}
	unique := []string{"in", "on", "with", "interval", "inheritance", "childrenonly", "extern"}
	required := []string{"in", "on", "with", "interval"}

	constraintTypes := []string{"service", "attribute", "system", "custom", "oncall", "native"}
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

		if u.SliceContainsString(val, keys) {
			// there must be at least one arguments left
			if len(args[pos+1:]) < 1 {
				u.abort("Syntax error, incomplete key/value specification (too few items left to parse)")
			}
			// check for back-to-back keyswords
			u.CheckStringNotAKeyword(args[pos+1], keys)

			switch val {
			case "threshold":
				if len(args[pos+1:]) < 6 {
					u.abort("Syntax error, incomplete threshold specification")
				}
				t := u.parseVariadicArguments(
					[]string{},
					[]string{"predicate", "level", "value"},
					[]string{"predicate", "level", "value"},
					args[pos+1:pos+7])
				thr := proto.CheckConfigThreshold{}
				thr.Predicate.Symbol = t["predicate"][0]
				thr.Level.Name = t["level"][0]
				if thr.Value, err = strconv.ParseInt(t["value"][0], 10, 64); err != nil {
					u.abort(fmt.Sprintf("Syntax error, value argument not numeric: %s",
						t["value"][0]))
				}
				thresholds = append(thresholds, thr)
				skip = true
				skipcount = 6
				continue argloop
			case "constraint":
				// argument is the start of a constraint specification.
				// check we have enough arguments left
				if len(args[pos+1:]) < 3 {
					u.abort("Syntax error, incomplete constraint specification")
				}
				// check constraint type specification
				if !u.SliceContainsString(args[pos+1], constraintTypes) {
					u.abort(fmt.Sprintf("Syntax error, unknown contraint type: %s",
						args[pos+1]))
				}
				constr := proto.CheckConfigConstraint{}
				constr.ConstraintType = args[pos+1]
				switch constr.ConstraintType {
				case "service":
					constr.Service = &proto.PropertyService{}
					switch args[pos+2] {
					case "name":
						constr.Service.Name = args[pos+3]
					default:
						u.abort(fmt.Sprintf("Syntax error, can not constraint service to %s",
							args[pos+2]))
					}
				case "attribute":
					constr.Attribute = &proto.ServiceAttribute{
						Name:  args[pos+2],
						Value: args[pos+3],
					}
				case "system":
					constr.System = &proto.PropertySystem{
						Name:  args[pos+2],
						Value: args[pos+3],
					}
				case "custom":
					constr.Custom = &proto.PropertyCustom{
						Name:  args[pos+2],
						Value: args[pos+3],
					}
				case "oncall":
					constr.Oncall = &proto.PropertyOncall{}
					switch args[pos+2] {
					case "id":
						constr.Oncall.Id = args[pos+3]
					case "name":
						constr.Oncall.Name = args[pos+3]
					default:
						u.abort(fmt.Sprintf("Syntax error, can not constraint oncall to %s",
							args[pos+2]))
					}
				case "native":
					constr.Native = &proto.PropertyNative{
						Name:  args[pos+2],
						Value: args[pos+3],
					}
				}
				constraints = append(constraints, constr)
				skip = true
				skipcount = 3
				continue argloop
			case "on":
				result["on/type"] = append(result["on/type"], args[pos+1])
				result["on/object"] = append(result["on/object"], args[pos+2])
				result[val] = append(result[val], fmt.Sprintf("%s::%s", args[pos+1], args[pos+2]))
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
		u.abort(fmt.Sprintf("Syntax error, erroneus argument: %s", val))
	}

	// check if all required keywords were collected
	for _, key := range required {
		if _, ok := result[key]; !ok {
			u.abort(fmt.Sprintf("Syntax error, missing keyword: %s", key))
		}
	}

	// check if unique keywords were only specuified once
	for _, key := range unique {
		// check ok since unique may still be optional
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			u.abort(fmt.Sprintf("Syntax error, keyword must only be provided once: %s", key))
		}
	}

	return result, constraints, thresholds
}

func (u *SomaUtil) ParseVariadicCapabilityArguments(
	c *resty.Client,
	multKeys []string, // keys that may appear multiple times
	uniqKeys []string, // keys that are allowed at most once
	reqKeys []string, // keys that are required at least one
	args []string, // arguments to parse
) (map[string][]string, []proto.CapabilityConstraint) {
	// returns a map of slices of string
	result := make(map[string][]string)
	constr := make([]proto.CapabilityConstraint, 0)

	// merge key slices
	multKeys = append(multKeys, []string{"constraint", "demux"}...)
	keys := append(multKeys, uniqKeys...)

	// helper to skip over next value in args slice
	skip := false
	skipcount := 0

	for pos, val := range args {
		// skip current arg if last argument was a keyword
		if skip {
			skipcount--
			if skipcount == 0 {
				skip = false
			}
			continue
		}

		if u.SliceContainsString(val, keys) {
			// there must be at least one arguments left
			if len(args[pos+1:]) < 1 {
				u.abort("Syntax error, incomplete key/value specification (too few items left to parse)")
			}
			// check for back-to-back keyswords
			u.CheckStringNotAKeyword(args[pos+1], keys)

			switch val {
			case "constraint":
				// must be at least 3 items left
				if len(args[pos+1:]) < 3 {
					u.abort("Syntax error, incomplete constraint specification")
				}
				// constraint must be type `system` or `attribute`
				switch args[pos+1] {
				case "system":
					u.CheckStringIsSystemProperty(c, args[pos+2])
				case "attribute":
					u.CheckStringIsServiceAttribute(c, args[pos+2])
				default:
					u.abort(fmt.Sprintf("Syntax error, invalid constraint type: %s", args[pos+1]))
				}
				constr = append(constr, proto.CapabilityConstraint{
					Type:  args[pos+1],
					Name:  args[pos+2],
					Value: args[pos+3],
				})
				skip = true
				skipcount = 3
				continue
			case "demux":
				// argument to demux must be a service attribute
				u.CheckStringIsServiceAttribute(c, args[pos+1])
				fallthrough
			default:
				// append value of current keyword into result map
				result[val] = append(result[val], args[pos+1])
				skip = true
				skipcount = 1
				continue
			}
		}
		// keywords trigger continue before this
		// values after keywords are skip'ed
		// reaching this is an error
		u.abort(fmt.Sprintf("Syntax error, erroneus argument: %s", val))
	}

	// check if we managed to collect all required keywords
	for _, key := range reqKeys {
		// ok is false if slice is nil
		if _, ok := result[key]; !ok {
			u.abort(fmt.Sprintf("Syntax error, missing required keyword: %s", key))
		}
	}

	// check if unique keywords were only specified once
	for _, key := range uniqKeys {
		if sl, ok := result[key]; ok && (len(sl) > 1) {
			u.abort(fmt.Sprintf("Syntax error, keyword must only be provided once: %s", key))
		}
	}

	return result, constr
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
