package adm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
)

// Exported functions

// WRAPPER
func Perform(rqType, path, tmpl string, body interface{}, c *cli.Context) error {
	var (
		err  error
		resp *resty.Response
	)

	if strings.HasSuffix(rqType, `body`) && body == nil {
		goto noattachment
	}

	switch rqType {
	case `get`:
		resp, err = GetReq(path)
	case `delete`:
		resp, err = DeleteReq(path)
	case `deletebody`:
		resp, err = DeleteReqBody(body, path)
	case `putbody`:
		resp, err = PutReqBody(body, path)
	case `postbody`:
		resp, err = PostReqBody(body, path)
	case `patchbody`:
		resp, err = PatchReqBody(body, path)
	}

	if err != nil {
		return err
	}
	return FormatOut(c, resp, tmpl)

noattachment:
	return fmt.Errorf(`Missing body to client request that requires it.`)
}

func DecodedResponse(resp *resty.Response, res *proto.Result) error {
	var err error
	res, err = decodeResponse(resp)
	if err != nil {
		return err
	}
	return checkApplicationError(res)
}

// DELETE
func DeleteReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Delete(p))
}

func DeleteReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Delete(p))
}

// GET
func GetReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Get(p))
}

// PATCH
func PatchReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Patch(p))
}

// POST
func PostReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Post(p))
}

// PUT
func PutReq(p string) (*resty.Response, error) {
	return handleRequestOptions(client.R().Put(p))
}

func PutReqBody(body interface{}, p string) (*resty.Response, error) {
	return handleRequestOptions(
		client.R().SetBody(body).SetContentLength(true).Put(p))
}

// Private functions

func handleRequestOptions(resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() >= 300 {
		return resp, fmt.Errorf("Request error: %s", resp.Status())
	}

	if !(async || jobSave) {
		return resp, nil
	}

	var result *proto.Result
	if result, err = decodeResponse(resp); err != nil {
		return nil, err
	}

	if jobSave {
		if result.StatusCode == 202 && result.JobId != "" {
			cache.SaveJob(result.JobId, result.JobType)
		}
	}

	if async {
		asyncWait(result)
	}
	return resp, nil
}

func asyncWait(result *proto.Result) {
	if !async {
		return
	}

	if result.StatusCode == 202 && result.JobId != "" {
		fmt.Fprintf(os.Stderr, "Waiting for job: %s\n", result.JobId)
		_, err := PutReq(fmt.Sprintf("/jobs/%s", result.JobId))
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Wait error: %s\n", err.Error())
		}
	}
}

func decodeResponse(resp *resty.Response) (*proto.Result, error) {
	result := proto.Result{}
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
