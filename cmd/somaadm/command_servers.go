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
					{
						Name:         "update",
						Usage:        "Full update of server attributes (replace, not merge)",
						Description:  help.Text(`ServersUpdate`),
						Action:       runtime(cmdServerUpdate),
						BashComplete: cmpl.ServerUpdate,
					},
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
	unique := []string{
		`assetid`,
		`datacenter`,
		`location`,
		`online`}
	required := []string{
		`assetid`,
		`datacenter`,
		`location`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewServerRequest()
	req.Server.Name = c.Args().First()
	if err := adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Server.AssetId, 1); err != nil {
		return err
	}
	req.Server.Datacenter = opts[`datacenter`][0]
	req.Server.Location = opts[`location`][0]

	// optional argument: online
	if ov, ok := opts[`online`]; ok {
		if err := adm.ValidateBool(ov[0],
			&req.Server.IsOnline); err != nil {
			return err
		}
	} else {
		// online defaults to true
		req.Server.IsOnline = true
	}

	return adm.Perform(`postbody`, `/servers/`, `command`, req, c)
}

func cmdServerMarkAsDeleted(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	sid, err := adm.LookupServerId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/servers/%s", sid)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdServerPurgeDeleted(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	// TODO this will currently never return a deleted server
	sid, err := adm.LookupServerId(c.Args().First())
	if err != nil {
		return err
	}
	req := proto.NewServerRequest()
	req.Flags.Purge = true

	path := fmt.Sprintf("/servers/%s", sid)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

func cmdServerUpdate(c *cli.Context) error {
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Server to update not referenced by"+
			" UUID: %s", c.Args().First())
	}

	unique := []string{
		`name`,
		`assetid`,
		`datacenter`,
		`location`,
		`online`,
		`deleted`}
	required := []string{
		`name`,
		`assetid`,
		`datacenter`,
		`location`,
		`online`,
		`deleted`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewServerRequest()
	req.Server.Id = c.Args().First()
	req.Server.Name = opts[`name`][0]
	req.Server.Datacenter = opts[`datacenter`][0]
	req.Server.Location = opts[`location`][0]
	if err := adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Server.AssetId, 1); err != nil {
		return err
	}
	if err := adm.ValidateBool(opts[`online`][0],
		&req.Server.IsOnline); err != nil {
		return err
	}
	if err := adm.ValidateBool(opts[`deleted`][0],
		&req.Server.IsDeleted); err != nil {
		return err
	}

	path := fmt.Sprintf("/servers/%s", c.Args().First())
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdServerRename(c *cli.Context) error {
	return fmt.Errorf(`Not implemented.`)
}

func cmdServerOnline(c *cli.Context) error {
	return fmt.Errorf(`Not implemented.`)
}

func cmdServerOffline(c *cli.Context) error {
	return fmt.Errorf(`Not implemented.`)
}

func cmdServerMove(c *cli.Context) error {
	return fmt.Errorf(`Not implemented.`)
}

func cmdServerList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/servers/`, `list`, nil, c)
}

func cmdServerSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/servers/`, `list`, nil, c)
}

func cmdServerShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	serverId, err := adm.LookupServerId(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/servers/%s", serverId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdServerSyncRequest(c *cli.Context) error {
	return fmt.Errorf(`Not implemented.`)
}

func cmdServerNull(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`datacenter`},
		[]string{`datacenter`},
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	req := proto.NewServerRequest()
	req.Server.Id = "00000000-0000-0000-0000-000000000000"
	req.Server.Datacenter = opts["datacenter"][0]

	return adm.Perform(`postbody`, `/servers/null`, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
