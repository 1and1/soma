package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u SomaUtil) GetTeamIdByName(teamName string) uuid.UUID {
	url := u.ApiUrl
	url.Path = "/teams"

	var req somaproto.ProtoRequestTeam
	var err error
	req.Filter.TeamName = teamName

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Get(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		u.Log.Fatal(err)
	}

	u.CheckRestyResponse(resp)
	teamResult := u.DecodeProtoResultTeamFromResponse(resp)

	if teamName != teamResult.Teams[0].TeamName {
		u.Log.Fatal("Received result set for incorrect team")
	}
	return teamResult.Teams[0].TeamId
}

func (u SomaUtil) DecodeProtoResultTeamFromResponse(resp *resty.Response) *somaproto.ProtoResultTeam {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultTeam
	err := decoder.Decode(&res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding server response body\n")
		u.Log.Printf("Error decoding server response body\n")
		u.Log.Fatal(err)
	}
	if res.Code > 299 {
		fmt.Fprintf(os.Stderr, "Request failed: %d - %s\n",
			res.Code, res.Status)
		for _, e := range res.Text {
			fmt.Fprintf(os.Stderr, "%s\n", e)
			u.Log.Printf("%s\n", e)
		}
		os.Exit(1)
	}
	return &res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
