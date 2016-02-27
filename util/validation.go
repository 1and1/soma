package util

import (
	"fmt"
	"net/mail"
	"strconv"
	"unicode/utf8"
)

func (u *SomaUtil) ValidateStringAsNodeAssetId(s string) {
	_, err := strconv.ParseUint(s, 10, 64)
	u.AbortOnError(err)
}

func (u *SomaUtil) ValidateStringAsBool(s string) {
	_, err := strconv.ParseBool(s)
	u.AbortOnError(err)
}

func (u *SomaUtil) GetValidatedBool(s string) bool {
	b, err := strconv.ParseBool(s)
	u.AbortOnError(err)
	return b
}

func (u *SomaUtil) ValidateRuneCount(s string, l int) string {
	if cnt := utf8.RuneCountInString(s); cnt > l {
		u.Abort(fmt.Sprintf("Error, string above character limit %d: %s",
			l, s))
	}
	return s
}

func (u *SomaUtil) GetValidatedUint64(s string, min uint64) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	u.AbortOnError(err)
	if i < min {
		u.Abort(fmt.Sprintf("Error, value %s is less than the minimun %d", i, min))
	}
	return i
}

func (u *SomaUtil) ValidateStringAsEmployeeNumber(s string) {
	employeeNumber, err := strconv.Atoi(s)
	u.AbortOnError(err, "Syntax error, employeenr argument not a number")
	if employeeNumber < 0 {
		u.Abort("Negative employee number is not allowed")
	}
}

func (u *SomaUtil) ValidateStringAsMailAddress(s string) {
	_, err := mail.ParseAddress(s)
	u.AbortOnError(err, "Syntax error, mailaddr does not parse as RFC 5322 address")
}

func (u *SomaUtil) ValidateStringInSlice(s string, sl []string) {
	if !u.SliceContainsString(s, sl) {
		if len(sl) == 0 {
			u.Abort("Error, cannot compare '%s' to empty keyword list")
		}

		// default []string.String() method prints [elem0 elem1 ...]
		// without quoting whitespace: []string{"a", "a b",} -> "[a a b]"
		errStr := fmt.Sprintf("Error, '%s' not one of: '%s'", s, sl[0])
		for _, v := range sl[1:] {
			errStr = fmt.Sprintf("%s, '%s'", errStr, v)
		}
		u.Abort(errStr)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
