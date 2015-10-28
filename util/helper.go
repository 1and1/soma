package util

import (
	"fmt"

	"gopkg.in/resty.v0"
)

func (u *SomaUtil) CheckRestyResponse(resp *resty.Response) {
	if resp.StatusCode() >= 400 {
		u.Abort(fmt.Sprintf("Request error: %s\n", resp.Status()))
	}
}

func (u *SomaUtil) SliceContainsString(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

// XXX DEPRECATED FOR SliceContainsString
func (u *SomaUtil) StringIsKeyword(s string, keys []string) bool {
	for _, key := range keys {
		if key == s {
			return true
		}
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
