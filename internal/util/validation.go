package util

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"unicode/utf8"

	"gopkg.in/resty.v0"
)

func (u *SomaUtil) ValidateStringAsNodeAssetId(s string) {
	_, err := strconv.ParseUint(s, 10, 64)
	u.abortOnError(err)
}

func (u *SomaUtil) ValidateStringAsBool(s string) {
	_, err := strconv.ParseBool(s)
	u.abortOnError(err)
}

func (u *SomaUtil) GetValidatedBool(s string) bool {
	b, err := strconv.ParseBool(s)
	u.abortOnError(err)
	return b
}

func (u *SomaUtil) validateRuneCount(s string, l int) string {
	if cnt := utf8.RuneCountInString(s); cnt > l {
		u.abort(fmt.Sprintf("Error, string above character limit %d: %s",
			l, s))
	}
	return s
}

func (u *SomaUtil) validateRuneCountRange(s string, lower, higher int) {
	if utf8.RuneCountInString(s) < lower || utf8.RuneCountInString(s) > higher {
		u.abort(fmt.Sprintf("Error, invalid string length (%d < length < %d)",
			lower, higher))
	}
}

func (u *SomaUtil) GetValidatedUint64(s string, min uint64) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	u.abortOnError(err)
	if i < min {
		u.abort(fmt.Sprintf("Error, value %d is less than the minimun %d", i, min))
	}
	return i
}

func (u *SomaUtil) validateStringAsEmployeeNumber(s string) {
	employeeNumber, err := strconv.Atoi(s)
	u.abortOnError(err, "Syntax error, employeenr argument not a number")
	if employeeNumber < 0 {
		u.abort("Negative employee number is not allowed")
	}
}

func (u *SomaUtil) validateStringAsMailAddress(s string) {
	_, err := mail.ParseAddress(s)
	u.abortOnError(err, "Syntax error, mailaddr does not parse as RFC 5322 address")
}

func (u *SomaUtil) ValidateStringInSlice(s string, sl []string) {
	if !u.SliceContainsString(s, sl) {
		if len(sl) == 0 {
			u.abort("Error, cannot compare '%s' to empty keyword list")
		}

		// default []string.String() method prints [elem0 elem1 ...]
		// without quoting whitespace: []string{"a", "a b",} -> "[a a b]"
		errStr := fmt.Sprintf("Error, '%s' not one of: '%s'", s, sl[0])
		for _, v := range sl[1:] {
			errStr = fmt.Sprintf("%s, '%s'", errStr, v)
		}
		u.abort(errStr)
	}
}

func (u *SomaUtil) ValidateProviderExists(c *resty.Client, s string) {
	resp := u.GetRequest(c, "/providers/")
	res := u.DecodeResultFromResponse(resp)

	if res.Providers != nil {
		for _, prov := range *res.Providers {
			if prov.Name == s {
				return
			}
		}
	}
	u.abort(fmt.Sprintf("Referenced provider %s is not registered with SOMA, see `somaadm providers help create`", s))
}

func (u *SomaUtil) ValidateUnitExists(c *resty.Client, s string) {
	resp := u.GetRequest(c, "/units/")
	res := u.DecodeResultFromResponse(resp)

	if res.Units != nil {
		for _, unit := range *res.Units {
			if unit.Unit == s {
				return
			}
		}
	}
	u.abort(fmt.Sprintf("Referenced unit %s is not registered with SOMA, see `somaadm units help create`", s))
}

func (u *SomaUtil) IsUUID(s string) bool {
	const reUUID string = `^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[1-5][[:xdigit:]]{3}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`
	const reUNIL string = `^0{8}-0{4}-0{4}-0{4}-0{12}$`
	re := regexp.MustCompile(fmt.Sprintf("%s|%s", reUUID, reUNIL))

	return re.MatchString(s)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
