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
	"regexp"
	"strconv"
	"unicode/utf8"
)

// ValidateRuneCount tests if a string's number of unicode runes is
// below an upper limit
func ValidateRuneCount(s string, l int) error {
	if cnt := utf8.RuneCountInString(s); cnt > l {
		return fmt.Errorf("Validation error, string '%s' above character limit %d.",
			s, l)
	}
	return nil
}

// ValidateRuneCountRange tests if a string's number of unicode runes is
// between an upper and lower bound
func ValidateRuneCountRange(s string, lower, higher int) error {
	if utf8.RuneCountInString(s) < lower || utf8.RuneCountInString(s) > higher {
		return fmt.Errorf("Validation error, string '%s' outside permitted length."+
			"Required: %d < len(%s) < %d.", s, lower, s, higher)
	}
	return nil
}

// isUUID validates if a string is one very narrow formatting of a UUID,
// namely the one used by the server. Other valid formats with braces etc
// are not accepted
func IsUUID(s string) bool {
	const reUUID string = `^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[1-5][[:xdigit:]]{3}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`
	const reUNIL string = `^0{8}-0{4}-0{4}-0{4}-0{12}$`
	re := regexp.MustCompile(fmt.Sprintf("%s|%s", reUUID, reUNIL))

	return re.MatchString(s)
}

// ValidateOncallNumber tests if a string is a 4 digit number
func ValidateOncallNumber(n string) error {
	num, err := strconv.Atoi(n)
	if err != nil {
		return fmt.Errorf("Syntax error, argument is not a number: %s", err.Error())
	}

	if num <= 0 || num > 9999 {
		return fmt.Errorf("Phone number must be 4-digit extension")
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
