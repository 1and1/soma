package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) TryGetTeamByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetTeamIdByName(s)
	}
	return id.String()
}

func (u SomaUtil) GetTeamIdByName(teamName string) string {
	req := somaproto.ProtoRequestTeam{}
	req.Filter = &somaproto.ProtoTeamFilter{}
	req.Filter.Name = teamName

	resp := u.PostRequestWithBody(req, "/filter/teams/")
	teamResult := u.DecodeProtoResultTeamFromResponse(resp)

	if teamName != teamResult.Teams[0].Name {
		u.Log.Fatal("Received result set for incorrect team")
	}
	return teamResult.Teams[0].Id
}

func (u SomaUtil) DecodeProtoResultTeamFromResponse(resp *resty.Response) *somaproto.ProtoResultTeam {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body()))
	var res somaproto.ProtoResultTeam
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
