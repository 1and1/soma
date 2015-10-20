package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func registerCommands(app cli.App) *cli.App {

	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "initialize local client files",
			Before: runtimePreCmd,
			Action: cmdClientInit,
		}, // end init
		{
			Name:   "servers",
			Usage:  "subcommands for servers",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create a new physical server",
					Action: func(c *cli.Context) {
						fmt.Printf("servers/create\n")
					},
				},
			},
		},
		{
			Name:   "buckets",
			Usage:  "subcommands for buckets",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create a new bucket",
					Action: func(c *cli.Context) {
						fmt.Printf("buckets/create\n")
					},
				},
				{
					Name:  "property",
					Usage: "subcommands for bucket properties",
					Subcommands: []cli.Command{
						{
							Name:        "add",
							Usage:       "add a property",
							Description: descBucketsPropertyAdd,
							Action: func(c *cli.Context) {
								fmt.Printf("buckets/property/add\n")
							},
						},
					},
				},
			},
		}, // end buckets
		{
			Name:   "views",
			Usage:  "subcommands for views",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "",
					Action: cmdViewsAdd,
				},
				{
					Name:   "remove",
					Usage:  "",
					Action: cmdViewsRemove,
				},
				{
					Name:   "rename",
					Usage:  "",
					Action: cmdViewsRename,
				},
				{
					Name:   "list",
					Usage:  "",
					Action: cmdViewsList,
				},
				{
					Name:   "show",
					Usage:  "",
					Action: cmdViewsShow,
				},
			},
		}, // end views
		{
			Name:   "environments",
			Usage:  "subcommands for environments",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "",
					Action: cmdEnvironmentsAdd,
				},
				{
					Name:   "remove",
					Usage:  "",
					Action: cmdEnvironmentsRemove,
				},
				{
					Name:   "rename",
					Usage:  "",
					Action: cmdEnvironmentsRename,
				},
				{
					Name:   "list",
					Usage:  "",
					Action: cmdEnvironmentsList,
				},
				{
					Name:   "show",
					Usage:  "",
					Action: cmdEnvironmentsShow,
				},
			},
		}, // end environments
		{
			Name:   "types",
			Usage:  "subcommands for object types",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "",
					Action: cmdObjectTypesAdd,
				},
				{
					Name:   "remove",
					Usage:  "",
					Action: cmdObjectTypesRemove,
				},
				{
					Name:   "rename",
					Usage:  "",
					Action: cmdObjectTypesRename,
				},
				{
					Name:   "list",
					Usage:  "",
					Action: cmdObjectTypesList,
				},
				{
					Name:   "show",
					Usage:  "",
					Action: cmdObjectTypesShow,
				},
			},
		}, // end types
		{
			Name:   "states",
			Usage:  "subcommands for states",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "",
					Action: cmdObjectStatesAdd,
				},
				{
					Name:   "remove",
					Usage:  "",
					Action: cmdObjectStatesRemove,
				},
				{
					Name:   "rename",
					Usage:  "",
					Action: cmdObjectStatesRename,
				},
				{
					Name:   "list",
					Usage:  "",
					Action: cmdObjectStatesList,
				},
				{
					Name:   "show",
					Usage:  "",
					Action: cmdObjectStatesShow,
				},
			},
		}, // end states
		{
			Name:   "datacenters",
			Usage:  "subcommands for datacenters",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "",
					Action: cmdDatacentersAdd,
				},
				{
					Name:   "remove",
					Usage:  "",
					Action: cmdDatacentersRemove,
				},
				{
					Name:   "rename",
					Usage:  "",
					Action: cmdDatacentersRename,
				},
				{
					Name:   "list",
					Usage:  "",
					Action: cmdDatacentersList,
				},
				{
					Name:   "show",
					Usage:  "",
					Action: cmdDatacentersShow,
				},
				{
					Name:   "groupadd",
					Usage:  "",
					Action: cmdDatacentersAddToGroup,
				},
				{
					Name:   "groupdel",
					Usage:  "",
					Action: cmdDatacentersRemoveFromGroup,
				},
				{
					Name:   "grouplist",
					Usage:  "",
					Action: cmdDatacentersListGroups,
				},
				{
					Name:   "groupshow",
					Usage:  "",
					Action: cmdDatacentersShowGroup,
				},
			},
		}, // end datacenters
		{
			Name:   "servers",
			Usage:  "subcommands for servers",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "",
					Action: cmdServerCreate,
				},
				{
					Name:   "delete",
					Usage:  "",
					Action: cmdServerMarkAsDeleted,
				},
				{
					Name:   "purge",
					Usage:  "",
					Action: cmdServerPurgeDeleted,
				},
				{
					Name:   "update",
					Usage:  "",
					Action: cmdServerUpdate,
				},
				{
					Name:   "rename",
					Usage:  "",
					Action: cmdServerRename,
				},
				{
					Name:   "online",
					Usage:  "",
					Action: cmdServerOnline,
				},
				{
					Name:   "offline",
					Usage:  "",
					Action: cmdServerOffline,
				},
				{
					Name:   "move",
					Usage:  "",
					Action: cmdServerMove,
				},
				{
					Name:   "list",
					Usage:  "",
					Action: cmdServerList,
				},
			},
		}, // end servers
		{
			Name:   "permissions",
			Usage:  "subcommands for permissions",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:  "type",
					Usage: "subcommands for permission types",
					Subcommands: []cli.Command{
						{
							Name:   "add",
							Action: cmdPermissionTypeAdd,
						},
						{
							Name:   "remove",
							Action: cmdPermissionTypeDel,
						},
						{
							Name:   "rename",
							Action: cmdPermissionTypeRename,
						},
						{
							Name:   "list",
							Action: cmdPermissionTypeList,
						},
						{
							Name:   "show",
							Action: cmdPermissionTypeShow,
						},
					}, // end permissions type
				},
				{
					Name:   "add",
					Action: cmdPermissionAdd,
				},
				{
					Name:   "remove",
					Action: cmdPermissionDel,
				},
				{
					Name:   "list",
					Action: cmdPermissionList,
				},
				{
					Name:  "show",
					Usage: "subcommands for permission show",
					Subcommands: []cli.Command{
						{
							Name:   "user",
							Action: cmdPermissionShowUser,
						},
						{
							Name:   "team",
							Action: cmdPermissionShowTeam,
						},
						{
							Name:   "tool",
							Action: cmdPermissionShowTool,
						},
						{
							Name:   "permission",
							Action: cmdPermissionShowPermission,
						},
					},
				}, // end permissions show
				{
					Name:   "audit",
					Action: cmdPermissionAudit,
				},
				{
					Name:  "grant",
					Usage: "subcommands for permission grant",
					Subcommands: []cli.Command{
						{
							Name:   "enable",
							Action: cmdPermissionGrantEnable,
						},
						{
							Name:   "global",
							Action: cmdPermissionGrantGlobal,
						},
						{
							Name:   "limited",
							Action: cmdPermissionGrantLimited,
						},
						{
							Name:   "system",
							Action: cmdPermissionGrantSystem,
						},
					},
				}, // end permissions grant
			},
		},
	}
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
