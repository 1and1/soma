package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/1and1/soma/lib/proto"
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
	u.abortOnError(err, "Error decoding server response body")
	if res.StatusCode > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.StatusCode, res.StatusText)
		msgs := []string{s}
		if res.Errors != nil { // pointer to slice
			msgs = append(msgs, *res.Errors...)
		}
		u.abort(msgs...)
	}
	return &res
}

func (u SomaUtil) unfilteredResultFromResponse(resp *resty.Response) *proto.Result {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	res := proto.Result{}
	err := decoder.Decode(&res)
	u.abortOnError(err, "Error decoding server response body")
	return &res
}

func (u SomaUtil) resultFromResponse(resp *resty.Response) (*proto.Result, error) {
	if resp.StatusCode() > 299 {
		return nil, SomaError{
			code: resp.StatusCode(),
			text: fmt.Sprintf("Received HTTP statuscode %d", resp.StatusCode()),
		}
	}
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	res := proto.Result{}
	err := decoder.Decode(&res)
	if err != nil {
		return nil, SomaError{
			code:     resp.StatusCode(),
			somaCode: 500,
			text:     fmt.Sprintf("somaadm: %s", err.Error()),
		}
	}
	if res.StatusCode > 299 {
		txt := fmt.Sprintf("SOMA returned internal code %d: %s",
			res.StatusCode,
			res.StatusText,
		)
		if res.Errors != nil && len(*res.Errors) > 0 {
			a := []string{txt}
			a = append(a, *res.Errors...)
			txt = strings.Join(a, `,`)
		}
		return nil, SomaError{
			code:     resp.StatusCode(),
			somaCode: res.StatusCode,
			text:     txt,
		}
	}
	return &res, nil
}

func (u SomaUtil) VerifyEnvironment(c *resty.Client, env string) {
	resp := u.GetRequest(c, "/environments/")
	res := u.DecodeResultFromResponse(resp)
	for _, e := range *res.Environments {
		if e.Name == env {
			return
		}
	}
	u.abort(fmt.Sprintf("Invalid environment specified: %s", env))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
