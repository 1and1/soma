package util

import (
	//	"fmt"
	"net/mail"
	"strconv"
)

func (u *SomaUtil) ValidateStringAsNodeAssetId(s string) {
	_, err := strconv.ParseUint(s, 10, 64)
	u.AbortOnError(err)
}

func (u *SomaUtil) ValidateStringAsBool(s string) {
	_, err := strconv.ParseBool(s)
	u.AbortOnError(err)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
