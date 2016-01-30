package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetUserByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetUserIdByName(s)
	}
	return id.String()
}

func (u *SomaUtil) GetUserIdByName(user string) string {
	req := somaproto.ProtoRequestUser{}
	req.Filter = &somaproto.ProtoUserFilter{}
	req.Filter.UserName = user

	resp := u.PostRequestWithBody(req, "/filter/users/")
	userResult := u.DecodeProtoResultUserFromResponse(resp)

	if user != userResult.Users[0].UserName {
		u.Abort("Received result set for incorrect user")
	}
	return userResult.Users[0].Id
}

func (u *SomaUtil) DecodeProtoResultUserFromResponse(resp *resty.Response) *somaproto.ProtoResultUser {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
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
