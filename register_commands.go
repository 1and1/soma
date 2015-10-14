package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
)

func registerCommands(app cli.App) *cli.App {
	log.Print("Registering cli commands")

	app.Commands = []cli.Command{
		{
			Name:   "servers",
			Usage:  "subcommands for servers",
			Before: configSetup,
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
			Before: configSetup,
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
			Before: configSetup,
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
			Before: configSetup,
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
			Before: configSetup,
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
			Before: configSetup,
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
			Before: configSetup,
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
		}, // end datacenters
	}
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
