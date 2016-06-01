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

func (u SomaUtil) TryGetServerByUUIDOrName(cache *db.DB, c *resty.Client, s string) string {
	id, err := uuid.FromString(s)
	if err != nil {
		if aid, err := strconv.ParseUint(s, 10, 64); err != nil {
			if m, err := cache.ServerByName(s); err != nil {
				// aborts on failure
				return u.GetServerIdByName(cache, c, s)
			} else {
				return m[`id`]
			}
		} else {
			if m, err := cache.ServerByAsset(s); err != nil {
				return u.GetServerIdByAssetId(cache, c, aid)
			} else {
				return m[`id`]
			}
		}
	}
	return id.String()
}

func (u SomaUtil) GetServerIdByAssetId(cache *db.DB, c *resty.Client, aid uint64) string {
	req := proto.NewServerFilter()
	req.Filter.Server.AssetId = aid

	resp := u.PostRequestWithBody(c, req, "/filter/servers/")
	res := u.DecodeResultFromResponse(resp)

	if aid != (*res.Servers)[0].AssetId {
		u.Abort("Received result set for incorrect server")
	}
	cache.Server(
		(*res.Servers)[0].Name,
		(*res.Servers)[0].Id,
		strconv.Itoa(int((*res.Servers)[0].AssetId)),
	)
	return (*res.Servers)[0].Id
}

func (u SomaUtil) GetServerIdByName(cache *db.DB, c *resty.Client, server string) string {
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
	cache.Server(
		(*serverResult.Servers)[0].Name,
		(*serverResult.Servers)[0].Id,
		strconv.Itoa(int((*serverResult.Servers)[0].AssetId)),
	)
	return (*serverResult.Servers)[0].Id
}

func (u SomaUtil) GetServerAssetIdByName(c *resty.Client, serverName string) uint64 {
	req := proto.Request{
		Filter: &proto.Filter{
			Server: &proto.ServerFilter{
				Name:     serverName,
				IsOnline: true,
			},
		},
	}

	resp := u.PostRequestWithBody(c, req, "/filter/servers/")
	serverResult := u.DecodeProtoResultServerFromResponse(resp)

	if serverName != (*serverResult.Servers)[0].Name {
		u.Log.Fatal("Received result set for incorrect server")
	}
	return (*serverResult.Servers)[0].AssetId
}

func (u SomaUtil) DecodeProtoResultServerFromResponse(resp *resty.Response) *proto.Result {
	return u.DecodeResultFromResponse(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
