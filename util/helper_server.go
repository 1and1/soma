package util

import (
	"fmt"
	"os"
	"strconv"

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

func (u SomaUtil) TryGetServerByUUIDOrName(c *resty.Client, s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		if aid, err := strconv.ParseUint(s, 10, 64); err != nil {
			// aborts on failure
			return u.GetServerIdByName(c, s)
		} else {
			return u.GetServerIdByAssetId(c, aid)
		}
	}
	return id.String()
}

func (u SomaUtil) GetServerIdByAssetId(c *resty.Client, aid uint64) string {
	req := proto.NewServerFilter()
	req.Filter.Server.AssetId = aid

	resp := u.PostRequestWithBody(c, req, "/filter/servers/")
	res := u.DecodeResultFromResponse(resp)

	if aid != (*res.Servers)[0].AssetId {
		u.Abort("Received result set for incorrect server")
	}
	return (*res.Servers)[0].Id
}

func (u SomaUtil) GetServerIdByName(c *resty.Client, server string) string {
	req := proto.Request{
		Filter: &proto.Filter{
			Server: &proto.ServerFilter{
				Name: server,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/servers/")
	serverResult := u.DecodeProtoResultServerFromResponse(resp)

	if server != (*serverResult.Servers)[0].Name {
		u.Abort("Received result set for incorrect server")
	}
	return (*serverResult.Servers)[0].Id
}

func (u SomaUtil) GetServerAssetIdByName(serverName string) uint64 {
	url := u.ApiUrl
	url.Path = "/servers"

	var err error
	req := proto.Request{
		Filter: &proto.Filter{
			Server: &proto.ServerFilter{
				Name:     serverName,
				IsOnline: true,
			},
		},
	}

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
	if len(*serverResult.Servers) != 1 {
		u.Log.Fatal("Unexpected result set length - expected one server result")
	}
	if serverName != (*serverResult.Servers)[0].Name {
		u.Log.Fatal("Received result set for incorrect server")
	}
	return (*serverResult.Servers)[0].AssetId
}

func (u SomaUtil) DecodeProtoResultServerFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
