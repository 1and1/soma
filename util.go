package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/resty.v0"
	"log"
	"net/url"
	"os"
)

var Slog *log.Logger

func initLogFile() {
	f, err := os.OpenFile(".somaadm.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	Slog = log.New(f, "somaadm: ", log.Ldate|log.Ltime|log.LUTC)

	Slog.Print("Initialized logfile")
}

func getApiUrl() *url.URL {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		Slog.Printf("Error parsing API address from config file")
		Slog.Fatal(err)
	}
	return url
}

func checkServerKeyword(s string) {
	keywords := []string{"id", "datacenter", "location", "name", "online"}
	for _, k := range keywords {
		if s == k {
			fmt.Fprintf(os.Stderr, "Syntax error: back-to-back keywords")
			os.Exit(1)
		}
	}
}

func getServerAssetIdByName(serverName string) uint64 {
	url := getApiUrl()
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
		Slog.Fatal(err)
	}

	checkRestyResponse(resp)
	serverResult := decodeProtoResultServerFromResponse(resp)

	if len(serverResult.Servers) != 1 {
		Slog.Fatal("Unexpected result set length - expected one server result")
	}
	if serverName != serverResult.Servers[0].Name {
		Slog.Fatal("Received result set for incorrect server")
	}
	return serverResult.Servers[0].AssetId
}

func checkRestyResponse(resp *resty.Response) {
	if resp.StatusCode() >= 400 {
		fmt.Fprintf(os.Stderr, "Request error: %s\n", resp.Status())
		Slog.Fatal(resp.Status())
	}
}

func decodeProtoResultServerFromResponse(resp *resty.Response) *somaproto.ProtoResultServer {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultServer
	err := decoder.Decode(&res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding server response body\n")
		Slog.Printf("Error decoding server response body\n")
		Slog.Fatal(err)
	}
	if res.Code > 299 {
		fmt.Fprintf(os.Stderr, "Request failed: %d - %s\n",
			res.Code, res.Status)
		for _, e := range res.Text {
			fmt.Fprintf(os.Stderr, "%s\n", e)
			Slog.Printf("%s\n", e)
		}
		os.Exit(1)
	}
	return &res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
