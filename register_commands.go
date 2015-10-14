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
	}
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
