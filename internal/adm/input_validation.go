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
	"net/mail"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/1and1/soma/lib/proto"
	resty "gopkg.in/resty.v0"
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

// IsUUID validates if a string is one very narrow formatting of a UUID,
// namely the one used by the server. Other valid formats with braces etc
// are not accepted
func IsUUID(s string) bool {
	const reUUID string = `^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[1-5][[:xdigit:]]{3}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`
	const reUNIL string = `^0{8}-0{4}-0{4}-0{4}-0{12}$`
	re := regexp.MustCompile(fmt.Sprintf("%s|%s", reUUID, reUNIL))

	return re.MatchString(s)
}

// isUint64 validates if a string is a uint64 number
func isUint64(s string) (bool, uint64) {
	if u, err := strconv.ParseUint(s, 10, 64); err != nil {
		return false, 0
	} else {
		return true, u
	}
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

// ValidateEmployeeNumber tests if a string is a non-negative number
func ValidateEmployeeNumber(s string) error {
	employeeNumber, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("Error: employeenr argument not a number")
	}
	if employeeNumber < 0 {
		return fmt.Errorf("Negative employee number is not allowed")
	}
	return nil
}

// ValidateMailAddress tests if a string is a mail address
func ValidateMailAddress(s string) error {
	if _, err := mail.ParseAddress(s); err != nil {
		return fmt.Errorf("Error: mailaddr does not parse as RFC 5322 address")
	}
	return nil
}

// ValidateBool tests that string s represents and bool and sets
// its value in b.
// Returns error if s could not be parsed as bool.
func ValidateBool(s string, b *bool) error {
	var err error
	*b, err = strconv.ParseBool(s)
	return err
}

// ValidateLBoundUint64 tests that string s is an unsigned number
// with a specific minimum value and sets it in i.
// Returns error if the number conversion failed or the minimum
// value was not reached.
func ValidateLBoundUint64(s string, i *uint64, min uint64) error {
	var err error
	if *i, err = strconv.ParseUint(s, 10, 64); err != nil {
		return err
	}
	if *i < min {
		return fmt.Errorf("Error, value %d is less than the minimun %d", i, min)
	}
	return nil
}

// ValidateUnit tests against the server if string s is a valid
// unit.
func ValidateUnit(s string) error {
	res, err := fetchObjList(`/units/`)
	if err != nil {
		return err
	}

	if res.Units != nil {
		for _, unit := range *res.Units {
			if unit.Unit == s {
				return nil
			}
		}
	}
	return fmt.Errorf("Value %s is not a valid unit", s)
}

// ValidateProvider tests aginst the server if string s is a valid
// provider.
func ValidateProvider(s string) error {
	res, err := fetchObjList(`/providers/`)
	if err != nil {
		return err
	}

	if res.Providers != nil {
		for _, prov := range *res.Providers {
			if prov.Name == s {
				return nil
			}
		}
	}
	return fmt.Errorf("Value %s is not a valid provider", s)
}

// fetchObjList is a helper for ValidateUnit and ValidateProvider
func fetchObjList(path string) (*proto.Result, error) {
	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = GetReq(path); err != nil {
		return nil, err
	}
	if res, err = decodeResponse(resp); err != nil {
		return nil, err
	}
	if err = checkApplicationError(res); err != nil {
		return nil, err
	}
	return res, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
