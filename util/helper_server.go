package util

import (
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

func (u SomaUtil) TryGetServerByUUIDOrName(s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		// aborts on failure
		return u.GetServerIdByName(s)
	}
	return id.String()
}

func (u SomaUtil) GetServerIdByName(server string) string {
	req := somaproto.ProtoRequestServer{}
	req.Filter = &somaproto.ProtoServerFilter{}
	req.Filter.Name = server

	resp := u.PostRequestWithBody(req, "/filter/servers/")
	serverResult := u.DecodeProtoResultServerFromResponse(resp)

	if server != serverResult.Servers[0].Name {
		u.Abort("Received result set for incorrect oncall duty")
	}
	return serverResult.Servers[0].Id
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

func (u SomaUtil) DecodeProtoResultServerFromResponse(resp *resty.Response) *somaproto.Result {
	return DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
