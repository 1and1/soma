package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
)

func registerServers(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// servers
			{
				Name:   "servers",
				Usage:  "SUBCOMMANDS for servers",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:        "create",
						Usage:       "Create a new physical server",
						Description: help.CmdServerCreate,
						Action:      cmdServerCreate,
					},
					{
						Name:   "delete",
						Usage:  "Mark an existing physical server as deleted",
						Action: cmdServerMarkAsDeleted,
					},
					{
						Name:   "purge",
						Usage:  "Remove all unreferenced servers marked as deleted",
						Action: cmdServerPurgeDeleted,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "all, a",
								Usage: "Purge all deleted servers",
							},
						},
					},
					{
						Name:   "update",
						Usage:  "Full update of server attributes (replace, not merge)",
						Action: cmdServerUpdate,
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing server",
						Action: cmdServerRename,
					},
					{
						Name:   "online",
						Usage:  "Set an existing server to online",
						Action: cmdServerOnline,
					},
					{
						Name:   "offline",
						Usage:  "Set an existing server to offline",
						Action: cmdServerOffline,
					},
					{
						Name:   "move",
						Usage:  "Change a server's registered location",
						Action: cmdServerMove,
					},
					{
						Name:   "list",
						Usage:  "List all servers, see full description for possible filters",
						Action: cmdServerList,
					},
					{
						Name:   "show",
						Usage:  "Show details about a specific server",
						Action: cmdServerShow,
					},
					{
						Name:   "sync",
						Usage:  "Request a data sync for a server",
						Action: cmdServerSyncRequest,
					},
					{
						Name:   "null",
						Usage:  "Bootstrap the null server",
						Action: cmdServerNull,
					},
				},
			}, // end servers
		}...,
	)
	return &app
}

func cmdServerCreate(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/servers/"

	// required gymnastics to get a []string
	a := c.Args()
	args := make([]string, 1)
	args[0] = a.First()
	tail := a.Tail()
	args = append(args, tail...)

	req := somaproto.ProtoRequestServer{}
	req.Server = &somaproto.ProtoServer{}
	var err error

	// golang on its own can't iterate over a slice two items at a time
	skipNext := false
	argumentCheck := map[string]bool{
		"id":         false,
		"datacenter": false,
		"location":   false,
		"name":       false,
		"online":     false,
	}
	for pos, val := range args {
		if skipNext {
			skipNext = false
			continue
		}
		switch val {
		case "id":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.AssetId, err = strconv.ParseUint(args[pos+1],
				10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"Cannot parse id argument to uint64\n")
				Slog.Fatal(err)
			}
			skipNext = true
			argumentCheck["id"] = true
		case "datacenter":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.Datacenter = args[pos+1]
			skipNext = true
			argumentCheck["datacenter"] = true
		case "location":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.Location = args[pos+1]
			skipNext = true
			argumentCheck["location"] = true
		case "name":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.Name = args[pos+1]
			skipNext = true
			argumentCheck["name"] = true
		case "online":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.IsOnline, err = strconv.ParseBool(args[pos+1])
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"parameter online must be true or false\n")
				Slog.Fatal(err)
			}
			skipNext = true
			argumentCheck["online"] = true
		}
	}

	// online argument is optional and defaults to true
	if !argumentCheck["online"] {
		argumentCheck["online"] = true
		req.Server.IsOnline = true
	}
	missingArgument := false
	for k, v := range argumentCheck {
		if !v {
			fmt.Fprintf(os.Stderr, "Missing argument: %s\n", k)
			missingArgument = true
		}
	}
	if missingArgument {
		os.Exit(1)
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
	utl.CheckRestyResponse(resp)
	// checks the embedded status code
	_ = utl.DecodeProtoResultServerFromResponse(resp)
	fmt.Println(resp)
}

func cmdServerMarkAsDeleted(c *cli.Context) {
	url := getApiUrl()
	var (
		assetId uint64
		err     error
	)

	a := c.Args()
	if !a.Present() {
		Slog.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		server := a.Get(1)
		if server == "" {
			Slog.Fatal("Syntax error")
		}
		assetId = utl.GetServerAssetIdByName(server)
	} else {
		assetId, err = strconv.ParseUint(a.First(), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			Slog.Fatal(err)
		}
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
}

func cmdServerPurgeDeleted(c *cli.Context) {
	url := getApiUrl()

	if c.Bool("all") {
		url.Path = fmt.Sprintf("/servers")
	} else {
		a := c.Args()
		if !a.Present() || len(a.Tail()) != 0 {
			Slog.Fatal("Syntax error")
		}
		assetId, err := strconv.ParseUint(a.First(), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			Slog.Fatal(err)
		}
		url.Path = fmt.Sprintf("/servers/%d", assetId)
	}

	var req somaproto.ProtoRequestServer
	req.Purge = true

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Delete(url.String())
	if err != nil {
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
}

func cmdServerUpdate(c *cli.Context) {
	url := getApiUrl()

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
	argumentCheck := map[string]bool{
		"id":         false,
		"datacenter": false,
		"location":   false,
		"name":       false,
		"online":     false,
	}
	for pos, val := range args {
		if skipNext {
			skipNext = false
			continue
		}
		switch val {
		case "id":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.AssetId, err = strconv.ParseUint(args[pos+1],
				10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"Cannot parse id argument to uint64\n")
				Slog.Fatal(err)
			}
			skipNext = true
			argumentCheck["id"] = true
		case "datacenter":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.Datacenter = args[pos+1]
			skipNext = true
			argumentCheck["datacenter"] = true
		case "location":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.Location = args[pos+1]
			skipNext = true
			argumentCheck["location"] = true
		case "name":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.Name = args[pos+1]
			skipNext = true
			argumentCheck["name"] = true
		case "online":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.IsOnline, err = strconv.ParseBool(args[pos+1])
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"parameter online must be true or false\n")
				Slog.Fatal(err)
			}
			skipNext = true
			argumentCheck["online"] = true
		}
	}

	// online argument is optional and defaults to true
	if !argumentCheck["online"] {
		argumentCheck["online"] = true
		req.Server.IsOnline = true
	}
	missingArgument := false
	for k, v := range argumentCheck {
		if !v {
			fmt.Fprintf(os.Stderr, "Missing argument: %s\n", k)
			missingArgument = true
		}
	}
	if missingArgument {
		os.Exit(1)
	}
	url.Path = fmt.Sprintf("/servers/%d", req.Server.AssetId)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// checks the embedded status code
	_ = utl.DecodeProtoResultServerFromResponse(resp)
}

func cmdServerRename(c *cli.Context) {
	url := getApiUrl()
	var (
		assetId uint64
		err     error
		newName string
	)

	a := c.Args()
	if !a.Present() {
		Slog.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		server := a.Get(1)
		if server == "" || a.Get(2) != "to" || a.Get(3) == "" {
			Slog.Fatal("Syntax error")
		}
		assetId = utl.GetServerAssetIdByName(server)
		newName = a.Get(3)
	} else {
		assetId, err = strconv.ParseUint(a.First(), 10, 64)
		if err != nil || a.Get(1) != "to" || a.Get(2) == "" {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			Slog.Fatal(err)
		}
		newName = a.Get(2)
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	var req somaproto.ProtoRequestServer
	req.Server.Name = newName

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
}

func cmdServerOnline(c *cli.Context) {
	url := getApiUrl()
	var (
		assetId uint64
		err     error
	)

	a := c.Args()
	if !a.Present() {
		Slog.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		server := a.Get(1)
		if server == "" {
			Slog.Fatal("Syntax error")
		}
		assetId = utl.GetServerAssetIdByName(server)
	} else {
		idString := a.First()
		if idString == "" {
			fmt.Fprintf(os.Stderr, "Could not read assetId\n")
			Slog.Fatal(err)
		}
		assetId, err = strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			Slog.Fatal(err)
		}
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	var req somaproto.ProtoRequestServer
	req.Server.IsOnline = true

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
}

func cmdServerOffline(c *cli.Context) {
	url := getApiUrl()
	var (
		assetId uint64
		err     error
	)

	a := c.Args()
	if !a.Present() {
		Slog.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		server := a.Get(1)
		if server == "" {
			Slog.Fatal("Syntax error")
		}
		assetId = utl.GetServerAssetIdByName(server)
	} else {
		idString := a.First()
		if idString == "" {
			fmt.Fprintf(os.Stderr, "Could not read assetId\n")
			Slog.Fatal(err)
		}
		assetId, err = strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			Slog.Fatal(err)
		}
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	var req somaproto.ProtoRequestServer
	req.Server.IsOnline = false

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
}

func cmdServerMove(c *cli.Context) {
	url := getApiUrl()
	var (
		assetId uint64
		err     error
	)

	a := c.Args()
	args := make([]string, 1)
	if !a.Present() {
		Slog.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		assetId = utl.GetServerAssetIdByName(a.Get(1))
		tail := a.Tail()
		subTail := tail[1:]
		args = append(args, subTail...)
	} else {
		assetId, err = strconv.ParseUint(a.First(), 10, 64)
		if err != nil || a.Get(1) != "to" || a.Get(2) == "" {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			Slog.Fatal(err)
		}
		tail := a.Tail()
		args = append(args, tail...)
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	var req somaproto.ProtoRequestServer

	// golang on its own can't iterate over a slice two items at a time
	skipNext := false
	argumentCheck := map[string]bool{
		"datacenter": false,
		"location":   false,
	}
	for pos, val := range args {
		if skipNext {
			skipNext = false
			continue
		}
		switch val {
		case "datacenter":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.Datacenter = args[pos+1]
			skipNext = true
			argumentCheck["datacenter"] = true
		case "location":
			utl.CheckServerKeyword(args[pos+1])
			req.Server.Datacenter = args[pos+1]
			skipNext = true
			argumentCheck["location"] = true
		}
	}
	missingArgument := false
	for k, v := range argumentCheck {
		if !v {
			fmt.Fprintf(os.Stderr, "Missing argument: %s\n", k)
			missingArgument = true
		}
	}
	if missingArgument {
		os.Exit(1)
	}

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// checks the embedded status code
	_ = utl.DecodeProtoResultServerFromResponse(resp)
}

func cmdServerList(c *cli.Context) {
	resp := utl.GetRequest("/servers/")
	fmt.Println(resp)
}

func cmdServerShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	serverId := utl.TryGetServerByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/servers/%s", serverId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdServerSyncRequest(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/jobs"

	a := c.Args()
	// arguments must be present, and the arguments after the first must
	// be zero => 1 argument given
	if !a.Present() || len(a.Tail()) != 0 {
		Slog.Fatal("Syntax error")
	}
	assetId, err := strconv.ParseUint(a.First(), 10, 64)

	var req somaproto.ProtoRequestJob
	req.JobType = "server"
	req.Server.Action = "sync"
	req.Server.Server.AssetId = assetId

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		Slog.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO save jobid locally as outstanding
	/*
		TODO: decode resp.Body -> ProtoResultJob
		job := ProtoResultJob.JobId
		jobDbAddOutstandingJob(job)
	*/
}

func cmdServerNull(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 2)
	key := []string{"datacenter"}

	opts := utl.ParseVariadicArguments(key, key, key, c.Args())

	req := somaproto.ProtoRequestServer{}
	req.Server = &somaproto.ProtoServer{}
	req.Server.Id = "00000000-0000-0000-0000-000000000000"
	req.Server.Datacenter = opts["datacenter"][0]

	resp := utl.PutRequestWithBody(req, "/servers/null")
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
