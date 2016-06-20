package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/resty.v0"
)

type SomaError struct {
	code     int
	somaCode uint16
	text     string
}

func (e SomaError) Error() string {
	return e.text
}

func (e SomaError) RequestError() bool {
	return e.code > 299
}

func (e SomaError) Code() uint16 {
	return e.somaCode
}

func (u SomaUtil) DecodeResultFromResponse(resp *resty.Response) *proto.Result {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	res := proto.Result{}
	err := decoder.Decode(&res)
	u.AbortOnError(err, "Error decoding server response body")
	if res.StatusCode > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.StatusCode, res.StatusText)
		msgs := []string{s}
		if res.Errors != nil { // pointer to slice
			msgs = append(msgs, *res.Errors...)
		}
		u.Abort(msgs...)
	}
	return &res
}

func (u SomaUtil) UnfilteredResultFromResponse(resp *resty.Response) *proto.Result {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	res := proto.Result{}
	err := decoder.Decode(&res)
	u.AbortOnError(err, "Error decoding server response body")
	return &res
}

func (u SomaUtil) VerifyEnvironment(c *resty.Client, env string) {
	resp := u.GetRequest(c, "/environments/")
	res := u.DecodeResultFromResponse(resp)
	for _, e := range *res.Environments {
		if e.Name == env {
			return
		}
	}
	u.Abort(fmt.Sprintf("Invalid environment specified: %s", env))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
