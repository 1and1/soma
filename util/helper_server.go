package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func (u *SomaUtil) CheckServerKeyword(s string) {
	keywords := []string{"id", "datacenter", "location", "name", "online"}
	for _, k := range keywords {
		if s == k {
			fmt.Fprintf(os.Stderr, "Syntax error: back-to-back keywords")
			os.Exit(1)
		}
	}
}

func (u SomaUtil) TryGetServerByUUIDOrName(s string) uuid.UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		id = u.GetNodeIdByName(s)
	}
	return id
}

func (u SomaUtil) GetServerAssetIdByName(serverName string) uint64 {
	url := u.ApiUrl
	url.Path = "/servers"

	var req somaproto.ProtoRequestServer
	var err error
	req.Filter.Name = serverName
	req.Filter.Online = true

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
	serverResult := u.DecodeProtoResultServerFromResponse(resp)

	// XXX really needed?
	if len(serverResult.Servers) != 1 {
		u.Log.Fatal("Unexpected result set length - expected one server result")
	}
	if serverName != serverResult.Servers[0].Name {
		u.Log.Fatal("Received result set for incorrect server")
	}
	return serverResult.Servers[0].AssetId
}

func (u SomaUtil) DecodeProtoResultServerFromResponse(resp *resty.Response) *somaproto.ProtoResultServer {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultServer
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
