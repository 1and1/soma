package main

import (
	"fmt"
	"log"
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
				Name:  "servers",
				Usage: "SUBCOMMANDS for servers",
				Subcommands: []cli.Command{
					{
						Name:        "create",
						Usage:       "Create a new physical server",
						Description: help.CmdServerCreate,
						Action:      runtime(cmdServerCreate),
					},
					{
						Name:   "delete",
						Usage:  "Mark an existing physical server as deleted",
						Action: runtime(cmdServerMarkAsDeleted),
					},
					{
						Name:   "purge",
						Usage:  "Remove all unreferenced servers marked as deleted",
						Action: runtime(cmdServerPurgeDeleted),
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
						Action: runtime(cmdServerUpdate),
					},
					{
						Name:   "rename",
						Usage:  "Rename an existing server",
						Action: runtime(cmdServerRename),
					},
					{
						Name:   "online",
						Usage:  "Set an existing server to online",
						Action: runtime(cmdServerOnline),
					},
					{
						Name:   "offline",
						Usage:  "Set an existing server to offline",
						Action: runtime(cmdServerOffline),
					},
					{
						Name:   "move",
						Usage:  "Change a server's registered location",
						Action: runtime(cmdServerMove),
					},
					{
						Name:   "list",
						Usage:  "List all servers, see full description for possible filters",
						Action: runtime(cmdServerList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a specific server",
						Action: runtime(cmdServerShow),
					},
					{
						Name:   "sync",
						Usage:  "Request a data sync for a server",
						Action: runtime(cmdServerSyncRequest),
					},
					{
						Name:   "null",
						Usage:  "Bootstrap the null server",
						Action: runtime(cmdServerNull),
					},
				},
			}, // end servers
		}...,
	)
	return &app
}

func cmdServerCreate(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	url.Path = "/servers/"

	// required gymnastics to get a []string
	a := c.Args()
	args := make([]string, 1)
	args[0] = a.First()
	tail := a.Tail()
	args = append(args, tail...)

	req := proto.Request{}
	req.Server = &proto.Server{}
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
				log.Fatal(err)
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
				log.Fatal(err)
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
		log.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// checks the embedded status code
	_ = utl.DecodeProtoResultServerFromResponse(resp)
	fmt.Println(resp)
	return nil
}

func cmdServerMarkAsDeleted(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	var (
		assetId uint64
		err     error
	)

	a := c.Args()
	if !a.Present() {
		log.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		server := a.Get(1)
		if server == "" {
			log.Fatal("Syntax error")
		}
		assetId = utl.GetServerAssetIdByName(server)
	} else {
		assetId, err = strconv.ParseUint(a.First(), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			log.Fatal(err)
		}
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		log.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
	return nil
}

func cmdServerPurgeDeleted(c *cli.Context) error {
	url := Cfg.Run.SomaAPI

	if c.Bool("all") {
		url.Path = fmt.Sprintf("/servers")
	} else {
		a := c.Args()
		if !a.Present() || len(a.Tail()) != 0 {
			log.Fatal("Syntax error")
		}
		assetId, err := strconv.ParseUint(a.First(), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			log.Fatal(err)
		}
		url.Path = fmt.Sprintf("/servers/%d", assetId)
	}

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Delete(url.String())
	if err != nil {
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
	return nil
}

func cmdServerUpdate(c *cli.Context) error {
	url := Cfg.Run.SomaAPI

	// required gymnastics to get a []string
	a := c.Args()
	args := make([]string, 1)
	args[0] = a.First()
	tail := a.Tail()
	args = append(args, tail...)

	var req proto.Request
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
				log.Fatal(err)
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
				log.Fatal(err)
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
		log.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// checks the embedded status code
	_ = utl.DecodeProtoResultServerFromResponse(resp)
	return nil
}

func cmdServerRename(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	var (
		assetId uint64
		err     error
		newName string
	)

	a := c.Args()
	if !a.Present() {
		log.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		server := a.Get(1)
		if server == "" || a.Get(2) != "to" || a.Get(3) == "" {
			log.Fatal("Syntax error")
		}
		assetId = utl.GetServerAssetIdByName(server)
		newName = a.Get(3)
	} else {
		assetId, err = strconv.ParseUint(a.First(), 10, 64)
		if err != nil || a.Get(1) != "to" || a.Get(2) == "" {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			log.Fatal(err)
		}
		newName = a.Get(2)
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	var req proto.Request
	req.Server.Name = newName

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		log.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
	return nil
}

func cmdServerOnline(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	var (
		assetId uint64
		err     error
	)

	a := c.Args()
	if !a.Present() {
		log.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		server := a.Get(1)
		if server == "" {
			log.Fatal("Syntax error")
		}
		assetId = utl.GetServerAssetIdByName(server)
	} else {
		idString := a.First()
		if idString == "" {
			fmt.Fprintf(os.Stderr, "Could not read assetId\n")
			log.Fatal(err)
		}
		assetId, err = strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			log.Fatal(err)
		}
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	var req proto.Request
	req.Server.IsOnline = true

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		log.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
	return nil
}

func cmdServerOffline(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	var (
		assetId uint64
		err     error
	)

	a := c.Args()
	if !a.Present() {
		log.Fatal("Syntax error")
	}
	if a.First() == "by-name" {
		server := a.Get(1)
		if server == "" {
			log.Fatal("Syntax error")
		}
		assetId = utl.GetServerAssetIdByName(server)
	} else {
		idString := a.First()
		if idString == "" {
			fmt.Fprintf(os.Stderr, "Could not read assetId\n")
			log.Fatal(err)
		}
		assetId, err = strconv.ParseUint(idString, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse assetId\n")
			log.Fatal(err)
		}
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	var req proto.Request
	req.Server.IsOnline = false

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		log.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// TODO check delete action success
	return nil
}

func cmdServerMove(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	var (
		assetId uint64
		err     error
	)

	a := c.Args()
	args := make([]string, 1)
	if !a.Present() {
		log.Fatal("Syntax error")
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
			log.Fatal(err)
		}
		tail := a.Tail()
		args = append(args, tail...)
	}
	url.Path = fmt.Sprintf("/servers/%d", assetId)

	var req proto.Request

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
		log.Fatal(err)
	}
	utl.CheckRestyResponse(resp)
	// checks the embedded status code
	_ = utl.DecodeProtoResultServerFromResponse(resp)
	return nil
}

func cmdServerList(c *cli.Context) error {
	resp := utl.GetRequest("/servers/")
	fmt.Println(resp)
	return nil
}

func cmdServerShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	serverId := utl.TryGetServerByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/servers/%s", serverId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	return nil
}

func cmdServerSyncRequest(c *cli.Context) error {
	/*
		url := getApiUrl()
		url.Path = "/jobs"

		a := c.Args()
		// arguments must be present, and the arguments after the first must
		// be zero => 1 argument given
		if !a.Present() || len(a.Tail()) != 0 {
			log.Fatal("Syntax error")
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
	*/
	return nil
}

func cmdServerNull(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 2)
	key := []string{"datacenter"}

	opts := utl.ParseVariadicArguments(key, key, key, c.Args())

	req := proto.Request{}
	req.Server = &proto.Server{}
	req.Server.Id = "00000000-0000-0000-0000-000000000000"
	req.Server.Datacenter = opts["datacenter"][0]

	resp := utl.PutRequestWithBody(req, "/servers/null")
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
