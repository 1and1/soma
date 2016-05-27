package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/resty.v0"
)

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

func (u SomaUtil) VerifyEnvironment(env string) {
	resp := u.GetRequest("/environments/")
	res := u.DecodeResultFromResponse(resp)
	for _, e := range *res.Environments {
		if e.Name == env {
			return
		}
	}
	u.Abort(fmt.Sprintf("Invalid environment specified: %s", env))
}

func (u SomaUtil) AsyncWait(enabled bool, resp *resty.Response) {
	if !enabled {
		return
	}
	r := u.DecodeResultFromResponse(resp)
	if r.StatusCode == 202 && r.JobId != "" {
		fmt.Fprintf(os.Stderr, "Waiting for job: %s\n", r.JobId)
		u.PutRequest(fmt.Sprintf("/jobs/%s", r.JobId))
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
