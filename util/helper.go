package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) GetUserIdByName(user string) uuid.UUID {
	url := *u.ApiUrl
	url.Path = "/users"

	var req somaproto.ProtoRequestUser
	var err error
	req.Filter.UserName = user

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Get(url.String())
	u.AbortOnError(err)
	u.CheckRestyResponse(resp)
	userResult := u.DecodeProtoResultUserFromResponse(resp)

	if user != userResult.Users[0].UserName {
		u.Abort("Received result set for incorrect user")
	}
	return userResult.Users[0].Id
}

func (u *SomaUtil) CheckRestyResponse(resp *resty.Response) {
	if resp.StatusCode() >= 400 {
		u.Abort(fmt.Sprintf("Request error: %s\n", resp.Status()))
	}
}

func (u *SomaUtil) DecodeProtoResultUserFromResponse(resp *resty.Response) *somaproto.ProtoResultUser {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultUser
	err := decoder.Decode(&res)
	u.AbortOnError(err, "Error decoding server response body")
	if res.Code > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.Code, res.Status)
		msgs := []string{s}
		msgs = append(msgs, res.Text...)
		u.Abort(msgs...)
	}
	return &res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
