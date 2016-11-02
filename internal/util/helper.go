package util

import (
	"fmt"
	"log"

	"gopkg.in/resty.v0"
)

func (u *SomaUtil) CheckRestyResponse(resp *resty.Response) {
	if resp.StatusCode() >= 400 {
		u.abort(fmt.Sprintf("Request error: %s\n", resp.Status()))
	}
}

func (u *SomaUtil) sliceContainsString(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

func (u *SomaUtil) CheckStringNotAKeyword(s string, keys []string) {
	if u.sliceContainsString(s, keys) {
		log.Fatal(`Syntax error, back-to-back keywords`)
	}
}

// XXX DEPRECATED FOR SliceContainsString
func (u *SomaUtil) stringIsKeyword(s string, keys []string) bool {
	for _, key := range keys {
		if key == s {
			return true
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
