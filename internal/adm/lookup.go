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

	"github.com/1and1/soma/lib/proto"
	resty "gopkg.in/resty.v0"
)

// LookupOncallId looks up the UUID for an oncall duty on the
// server with name s. Error is set if no such oncall duty was
// found or an error occured.
// If s is already a UUID, then s is immediately returned.
func LookupOncallId(s string) (string, error) {
	if isUUID(s) {
		return s, nil
	}
	return oncallIdByName(s)
}

// oncallIdByName implements the actual serverside lookup of the
// oncall duty UUID
func oncallIdByName(oncall string) (string, error) {
	req := proto.Request{}
	req.Filter = &proto.Filter{}
	req.Filter.Oncall = &proto.OncallFilter{Name: oncall}

	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	if resp, err = PostReqBody(req, `/filter/oncall/`); err != nil {
		// transport errors
		goto abort
	}

	if res, err = decodeResponse(resp); err != nil {
		// http code errors
		goto abort
	}

	if res.StatusCode >= 300 {
		// application errors
		s := fmt.Sprintf("Request failed: %d - %s",
			res.StatusCode, res.StatusText)
		m := []string{s}

		if res.Errors != nil {
			m = append(m, *res.Errors...)
		}

		err = fmt.Errorf(combineStrings(m...))
		goto abort
	}

	// check the received record against the input
	if oncall != (*res.Oncalls)[0].Name {
		err = fmt.Errorf("Name mismatch: %s vs %s",
			oncall, (*res.Oncalls)[0].Name)
		goto abort
	}
	return (*res.Oncalls)[0].Id, nil

abort:
	return ``, fmt.Errorf("OncallId lookup failed: %s", err.Error())
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
