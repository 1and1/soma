package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

var Slog *log.Logger

func initLogFile() {
	f, err := os.OpenFile(path.Join(Cfg.Run.PathLogs, "somaadm.log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logging: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	Cfg.Run.Logger = log.New(f, "", log.Ldate|log.Ltime|log.LUTC)
	Slog = Cfg.Run.Logger

	Slog.Print("Initialized logger")
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

	// XXX really needed?
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

func validateCliArgumentCount(c *cli.Context, i uint8) {
	a := c.Args()
	if i == 0 {
		if a.Present() {
			Slog.Fatal("Syntax error, command takes no arguments")
		}
	} else {
		if !a.Present() || len(a.Tail()) != (int(i)-1) {
			Slog.Fatal("Syntax error")
		}
	}
}

func validateCliArgument(c *cli.Context, pos uint8, s string) {
	a := c.Args()
	if a.Get(int(pos)-1) != s {
		Slog.Fatal("Syntax error, missing keyword: ", s)
	}
}

func getCliArgumentCount(c *cli.Context) int {
	a := c.Args()
	if !a.Present() {
		return 0
	}
	return len(a.Tail()) + 1
}

func parseLimitedGrantArguments(keys []string, args []string) *map[string]string {
	result := make(map[string]string)
	argumentCheck := make(map[string]bool)
	for _, key := range keys {
		argumentCheck[key] = false
	}

	skipNext := false

	for pos, val := range args {
		if skipNext {
			skipNext = false
			continue
		}
		switch val {
		case "repository":
			checkStringNotAKeyword(args[pos+1], keys)
			result["repository"] = args[pos+1]
			skipNext = true
			argumentCheck["repository"] = true
		case "bucket":
			checkStringNotAKeyword(args[pos+1], keys)
			result["bucket"] = args[pos+1]
			skipNext = true
			argumentCheck["bucket"] = true
		case "group":
			checkStringNotAKeyword(args[pos+1], keys)
			result["group"] = args[pos+1]
			skipNext = true
			argumentCheck["group"] = true
		case "cluster":
			checkStringNotAKeyword(args[pos+1], keys)
			result["group"] = args[pos+1]
			skipNext = true
			argumentCheck["cluster"] = true
		}
	}

	// check we managed to collect all required keywords
	for k, v := range argumentCheck {
		if !v {
			Slog.Fatal("Syntax error, missing keyword for argument count: ", k)
		}
	}

	return &result
}

func checkStringNotAKeyword(s string, keys []string) {
	for _, val := range keys {
		if val == s {
			Slog.Fatal("Syntax error, back-to-back keywords")
		}
	}
}

func getTeamIdByName(teamName string) uuid.UUID {
	url := getApiUrl()
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
		Slog.Fatal(err)
	}

	checkRestyResponse(resp)
	teamResult := decodeProtoResultTeamFromResponse(resp)

	if teamName != teamResult.Teams[0].TeamName {
		Slog.Fatal("Received result set for incorrect team")
	}
	return teamResult.Teams[0].TeamId
}

func getOncallIdByName(oncall string) uuid.UUID {
	url := getApiUrl()
	url.Path = "/oncall"

	var req somaproto.ProtoRequestOncall
	var err error
	req.Filter.Name = oncall

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Get(url.String())
	abortOnError(err)
	checkRestyResponse(resp)
	oncallResult := decodeProtoResultOncallFromResponse(resp)

	if oncall != oncallResult.Oncalls[0].Name {
		abort("Received result set for incorrect team")
	}
	return oncallResult.Oncalls[0].Id
}

func decodeProtoResultTeamFromResponse(resp *resty.Response) *somaproto.ProtoResultTeam {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultTeam
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

func decodeProtoResultOncallFromResponse(resp *resty.Response) *somaproto.ProtoResultOncall {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultOncall
	err := decoder.Decode(&res)
	abortOnError(err, "Error decoding server response body")
	if res.Code > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.Code, res.Status)
		msgs := []string{s}
		msgs = append(msgs, res.Text...)
		abort(msgs...)
	}
	return &res
}

func parseVariableArguments(keys []string, args []string) *map[string]string {
	result := make(map[string]string)
	argumentCheck := make(map[string]bool)
	for _, key := range keys {
		argumentCheck[key] = false
	}
	skipNext := false

	for pos, val := range args {
		if skipNext {
			skipNext = false
			continue
		}

		if stringIsKeyword(val, keys) {
			checkStringNotAKeyword(args[pos+1], keys)
			result[val] = args[pos+1]
			argumentCheck[val] = true
			skipNext = true
			continue
		}
		// keywords trigger continue, arguments are skipped over.
		// reaching this is an error
		Slog.Fatal("Syntax error, erroneus argument: ", val)
	}

	// check we managed to collect all required keywords
	for k, v := range argumentCheck {
		if !v {
			Slog.Fatal("Syntax error, missing keyword: ", k)
		}
	}

	return &result
}

func stringIsKeyword(s string, keys []string) bool {
	for _, key := range keys {
		if key == s {
			return true
		}
	}
	return false
}

func abortOnError(err error, txt ...string) {
	if err != nil {
		for _, s := range txt {
			fmt.Fprintf(os.Stderr, "%s\n", s)
			Slog.Print(s)
		}
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
}

func abort(txt ...string) {
	for _, s := range txt {
		fmt.Fprintf(os.Stderr, "%s\n", s)
		Slog.Print(s)
	}

	// ensure there is _something_
	if len(txt) == 0 {
		e := `abort called without error message. Sorry!`
		fmt.Fprintf(os.Stderr, "%s\n", e)
		Slog.Print(e)
	}
	os.Exit(1)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
