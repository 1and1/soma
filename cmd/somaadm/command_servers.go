package main

import (
	"fmt"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/internal/help"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
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
						Name:         "create",
						Usage:        "Create a new physical server",
						Description:  help.Text(`ServersCreate`),
						Action:       runtime(cmdServerCreate),
						BashComplete: cmpl.ServerCreate,
					},
					{
						Name:   "delete",
						Usage:  "Mark an existing physical server as deleted",
						Action: runtime(cmdServerMarkAsDeleted),
					},
					/*
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
					*/
					{
						Name:         "update",
						Usage:        "Full update of server attributes (replace, not merge)",
						Description:  help.Text(`ServersUpdate`),
						Action:       runtime(cmdServerUpdate),
						BashComplete: cmpl.ServerUpdate,
					},
					/*
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
					*/
					{
						Name:   "list",
						Usage:  "List all servers, see full description for possible filters",
						Action: runtime(cmdServerList),
					},
					{
						Name:   "synclist",
						Usage:  "Export a list of all servers suitable for sync",
						Action: runtime(cmdServerSync),
					},
					{
						Name:   "show",
						Usage:  "Show details about a specific server",
						Action: runtime(cmdServerShow),
					},
					/*
						{
							Name:   "sync",
							Usage:  "Request a data sync for a server",
							Action: runtime(cmdServerSyncRequest),
						},
					*/
					{
						Name:         "null",
						Usage:        "Bootstrap the null server",
						Action:       runtime(cmdServerNull),
						BashComplete: cmpl.Datacenter,
					},
				},
			}, // end servers
		}...,
	)
	return &app
}

func cmdServerCreate(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 7)

	multiple := []string{}
	unique := []string{`assetid`, `datacenter`, `location`, `online`}
	required := []string{`assetid`, `datacenter`, `location`}
	opts := adm.ParseVariadicArguments(multiple, unique, required,
		c.Args().Tail())

	req := proto.NewServerRequest()
	req.Server.Name = c.Args().First()
	req.Server.AssetId = utl.GetValidatedUint64(opts[`assetid`][0], 1)
	req.Server.Datacenter = opts[`datacenter`][0]
	req.Server.Location = opts[`location`][0]

	// optional argument: online
	if ov, ok := opts[`online`]; ok {
		req.Server.IsOnline = utl.GetValidatedBool(ov[0])
	} else {
		// online defaults to true
		req.Server.IsOnline = true
	}

	resp := utl.PostRequestWithBody(Client, req, `/servers/`)
	fmt.Println(resp)
	return nil
}

func cmdServerMarkAsDeleted(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	sid := utl.TryGetServerByUUIDOrName(&store, Client, c.Args().First())
	path := fmt.Sprintf("/servers/%s", sid)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdServerPurgeDeleted(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	// TODO this will currently never return a deleted server
	sid := utl.TryGetServerByUUIDOrName(&store, Client, c.Args().First())
	path := fmt.Sprintf("/servers/%s", sid)
	req := proto.NewServerRequest()
	req.Flags.Purge = true

	resp := utl.DeleteRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdServerUpdate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 13)

	if !utl.IsUUID(c.Args().First()) {
		adm.Abort(
			fmt.Sprintf("Server to update not referenced by UUID: %s",
				c.Args().First()))
	}

	multiple := []string{}
	unique := []string{`name`, `assetid`, `datacenter`, `location`, `online`, `deleted`}
	required := []string{`name`, `assetid`, `datacenter`, `location`, `online`, `deleted`}
	opts := adm.ParseVariadicArguments(multiple, unique, required,
		c.Args().Tail())

	req := proto.NewServerRequest()
	req.Server.Id = c.Args().First()
	req.Server.Name = opts[`name`][0]
	req.Server.AssetId = utl.GetValidatedUint64(opts[`assetid`][0], 1)
	req.Server.Datacenter = opts[`datacenter`][0]
	req.Server.Location = opts[`location`][0]
	req.Server.IsOnline = utl.GetValidatedBool(opts[`online`][0])
	req.Server.IsDeleted = utl.GetValidatedBool(opts[`deleted`][0])

	path := fmt.Sprintf("/servers/%s", c.Args().First())
	resp := utl.PutRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdServerRename(c *cli.Context) error {
	return nil
}

func cmdServerOnline(c *cli.Context) error {
	return nil
}

func cmdServerOffline(c *cli.Context) error {
	return nil
}

func cmdServerMove(c *cli.Context) error {
	return nil
}

func cmdServerList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp, err := adm.GetReq(`/servers/`)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func cmdServerSync(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp, err := adm.GetReq(`/sync/servers/`)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func cmdServerShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	serverId := utl.TryGetServerByUUIDOrName(&store, Client, c.Args().First())
	path := fmt.Sprintf("/servers/%s", serverId)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdServerSyncRequest(c *cli.Context) error {
	return nil
}

func cmdServerNull(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 2)
	key := []string{"datacenter"}

	opts := adm.ParseVariadicArguments(key, key, key, c.Args())

	req := proto.Request{}
	req.Server = &proto.Server{}
	req.Server.Id = "00000000-0000-0000-0000-000000000000"
	req.Server.Datacenter = opts["datacenter"][0]

	resp := utl.PostRequestWithBody(Client, req, "/servers/null")
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
