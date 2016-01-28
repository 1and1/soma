package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) TryGetUserByUUIDOrName(s string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		id = u.GetUserIdByName(s)
	}
	return id
}

func (u *SomaUtil) UserIdByUuidOrName(c *cli.Context) uuid.UUID {
	var (
		id  uuid.UUID
		err error
	)

	switch u.GetCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
		u.AbortOnError(err, "Syntax error, argument not a uuid")
	case 2:
		u.ValidateCliArgument(c, 1, "by-name")
		id = u.GetUserIdByName(c.Args().Get(1))
	default:
		u.Abort("Syntax error, unexpected argument count")
	}
	return id
}

func (u *SomaUtil) GetUserIdByName(user string) uuid.UUID {
	var req somaproto.ProtoRequestUser
	req.Filter.UserName = user

	resp := u.GetRequestWithBody(req, "/users")
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
