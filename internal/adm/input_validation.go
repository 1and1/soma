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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
