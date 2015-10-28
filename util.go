package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

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
			utl.CheckStringNotAKeyword(args[pos+1], keys)
			result["repository"] = args[pos+1]
			skipNext = true
			argumentCheck["repository"] = true
		case "bucket":
			utl.CheckStringNotAKeyword(args[pos+1], keys)
			result["bucket"] = args[pos+1]
			skipNext = true
			argumentCheck["bucket"] = true
		case "group":
			utl.CheckStringNotAKeyword(args[pos+1], keys)
			result["group"] = args[pos+1]
			skipNext = true
			argumentCheck["group"] = true
		case "cluster":
			utl.CheckStringNotAKeyword(args[pos+1], keys)
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

	utl.CheckRestyResponse(resp)
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
	utl.AbortOnError(err)
	utl.CheckRestyResponse(resp)
	oncallResult := decodeProtoResultOncallFromResponse(resp)

	if oncall != oncallResult.Oncalls[0].Name {
		utl.Abort("Received result set for incorrect oncall duty")
	}
	return oncallResult.Oncalls[0].Id
}

func getUserIdByName(user string) uuid.UUID {
	url := getApiUrl()
	url.Path = "/users"

	var req somaproto.ProtoRequestUser
	var err error
	req.Filter.UserName = user

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Get(url.String())
	utl.AbortOnError(err)
	utl.CheckRestyResponse(resp)
	userResult := decodeProtoResultUserFromResponse(resp)

	if user != userResult.Users[0].UserName {
		utl.Abort("Received result set for incorrect user")
	}
	return userResult.Users[0].Id
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
	utl.AbortOnError(err, "Error decoding server response body")
	if res.Code > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.Code, res.Status)
		msgs := []string{s}
		msgs = append(msgs, res.Text...)
		utl.Abort(msgs...)
	}
	return &res
}

func decodeProtoResultUserFromResponse(resp *resty.Response) *somaproto.ProtoResultUser {
	decoder := json.NewDecoder(bytes.NewReader(resp.Body))
	var res somaproto.ProtoResultUser
	err := decoder.Decode(&res)
	utl.AbortOnError(err, "Error decoding server response body")
	if res.Code > 299 {
		s := fmt.Sprintf("Request failed: %d - %s", res.Code, res.Status)
		msgs := []string{s}
		msgs = append(msgs, res.Text...)
		utl.Abort(msgs...)
	}
	return &res
}

func parseVariableArguments(keys []string, rKeys []string, args []string) (map[string]string, []string) {
	result := make(map[string]string)
	argumentCheck := make(map[string]bool)
	optionalKeys := make([]string, 0)
	for _, key := range rKeys {
		argumentCheck[key] = false
	}
	skipNext := false

	for pos, val := range args {
		if skipNext {
			skipNext = false
			continue
		}

		if utl.StringIsKeyword(val, keys) {
			utl.CheckStringNotAKeyword(args[pos+1], keys)
			result[val] = args[pos+1]
			argumentCheck[val] = true
			skipNext = true
			if !utl.StringIsKeyword(val, rKeys) {
				optionalKeys = append(optionalKeys, val)
			}
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

	return result, optionalKeys
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
