package main

import (
	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
	//"net/url"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func cmdServerCreate(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/servers"

	// required gymnastics to get a []string
	a := c.Args()
	args := make([]string, 1)
	args[0] = a.First()
	tail := a.Tail()
	args = append(args, tail...)

	var req somaproto.ProtoRequestServer
	var err error

	// golang on its own can't iterate over a slice two items at a time
	skipNext := false
	onlineArgGiven := false
	for pos, val := range args {
		if skipNext {
			skipNext = false
			continue
		}
		switch val {
		case "id":
			checkServerKeyword(args[pos+1])
			req.Server.AssetId, err = strconv.ParseUint(args[pos+1],
				10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"Cannot parse id argument to uint64\n")
				Slog.Fatal(err)
			}
			skipNext = true
		case "datacenter":
			checkServerKeyword(args[pos+1])
			req.Server.Datacenter = args[pos+1]
			skipNext = true
		case "location":
			checkServerKeyword(args[pos+1])
			req.Server.Location = args[pos+1]
			skipNext = true
		case "name":
			checkServerKeyword(args[pos+1])
			req.Server.Name = args[pos+1]
			skipNext = true
		case "online":
			checkServerKeyword(args[pos+1])
			req.Server.Online, err = strconv.ParseBool(args[pos+1])
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"parameter online must be true or false\n")
				Slog.Fatal(err)
			}
			skipNext = true
			onlineArgGiven = true
		}
	}

	// optional argument handling, false is the zero value of booleans,
	// so we need onlineArgGiven to detect if req.Server.Online was set
	// to false or is in its default field state.
	// Our servers are by default online
	if onlineArgGiven == false {
		req.Server.Online = true
	}

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	Slog.Printf("HTTP Response: %s\n", resp.Status())

	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var serverResult somaproto.ProtoResultServer
	err = decoder.Decode(&serverResult)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding server response body\n")
		Slog.Printf("Error decoding server response body\n")
		Slog.Fatal(err)
	}
	if serverResult.Code > 299 {
		fmt.Fprintf(os.Stderr, "Request failed: %d - %s\n",
			serverResult.Code, serverResult.Status)
		for _, e := range serverResult.Text {
			fmt.Fprintf(os.Stderr, "%s\n", e)
			Slog.Printf("%s\n", e)
		}
		os.Exit(1)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
