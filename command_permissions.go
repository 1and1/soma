package main

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
)

func registerPermissions(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// permissions
			{
				Name:  "permissions",
				Usage: "SUBCOMMANDS for permissions",
				Subcommands: []cli.Command{
					{
						Name:  "type",
						Usage: "SUBCOMMANDS for permission types",
						Subcommands: []cli.Command{
							{
								Name:   "add",
								Usage:  "Register a new permission type",
								Action: runtime(cmdPermissionTypeAdd),
							},
							{
								Name:   "remove",
								Usage:  "Remove an existing permission type",
								Action: runtime(cmdPermissionTypeDel),
							},
							{
								Name:   "rename",
								Usage:  "Rename an existing permission type",
								Action: runtime(cmdPermissionTypeRename),
							},
							{
								Name:   "list",
								Usage:  "List all permission types",
								Action: runtime(cmdPermissionTypeList),
							},
							{
								Name:   "show",
								Usage:  "Show details for a permission type",
								Action: runtime(cmdPermissionTypeShow),
							},
						}, // end permissions type
					},
					{
						Name:   "add",
						Usage:  "Register a new permission",
						Action: runtime(cmdPermissionAdd),
					},
					{
						Name:   "remove",
						Usage:  "Remove a permission",
						Action: runtime(cmdPermissionDel),
					},
					{
						Name:   "list",
						Usage:  "List all permissions",
						Action: runtime(cmdPermissionList),
					},
					{
						Name:  "show",
						Usage: "SUBCOMMANDS for permission show",
						Subcommands: []cli.Command{
							{
								Name:   "user",
								Usage:  "Show permissions of a user",
								Action: runtime(cmdPermissionShowUser),
							},
							{
								Name:   "team",
								Usage:  "Show permissions of a team",
								Action: runtime(cmdPermissionShowTeam),
							},
							{
								Name:   "tool",
								Usage:  "Show permissions of a tool account",
								Action: runtime(cmdPermissionShowTool),
							},
							{
								Name:   "permission",
								Usage:  "Show details about a permission",
								Action: runtime(cmdPermissionShowPermission),
							},
						},
					}, // end permissions show
					{
						Name:   "audit",
						Usage:  "Show all limited permissions associated with a repository",
						Action: runtime(cmdPermissionAudit),
					},
					{
						Name:  "grant",
						Usage: "SUBCOMMANDS for permission grant",
						Subcommands: []cli.Command{
							{
								Name:   "enable",
								Usage:  "Enable a useraccount to receive GRANT permissions",
								Action: runtime(cmdPermissionGrantEnable),
							},
							{
								Name:   "global",
								Usage:  "Grant a global permission",
								Action: runtime(cmdPermissionGrantGlobal),
							},
							{
								Name:   "limited",
								Usage:  "Grant a limited permission",
								Action: runtime(cmdPermissionGrantLimited),
							},
							{
								Name:   "system",
								Usage:  "Grant a system permission",
								Action: runtime(cmdPermissionGrantSystem),
							},
						},
					}, // end permissions grant
				},
			}, // end permissions
		}...,
	)
	return &app
}

func cmdPermissionTypeAdd(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	permissionType := c.Args().First()

	var req proto.Request
	req.Permission = &proto.Permission{}
	req.Permission.Category = permissionType

	resp := utl.PostRequestWithBody(req, "/permission/types/")
	fmt.Println(resp)
	return nil
}

func cmdPermissionTypeDel(c *cli.Context) error {
	url := Cfg.Run.SomaAPI

	utl.ValidateCliArgumentCount(c, 1)
	permissionType := c.Args().First()
	url.Path = fmt.Sprintf("/permissions/types/%s", permissionType)
	log.Printf("Command: delete permission type [%s]", permissionType)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func cmdPermissionTypeRename(c *cli.Context) error {
	url := Cfg.Run.SomaAPI

	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "to") // starts args counting at 1
	permissionType := c.Args().Get(0)
	newPermissionType := c.Args().Get(2)
	url.Path = fmt.Sprintf("/permissions/types/%s", permissionType)

	var req proto.Request
	req.Permission = &proto.Permission{}
	req.Permission.Category = newPermissionType

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func cmdPermissionTypeList(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	url.Path = "/permissions/types"

	utl.ValidateCliArgumentCount(c, 0)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	// TODO display list result
	return nil
}

func cmdPermissionTypeShow(c *cli.Context) error {
	url := Cfg.Run.SomaAPI

	utl.ValidateCliArgumentCount(c, 1)
	permissionType := c.Args().Get(0)
	url.Path = fmt.Sprintf("/permissions/types/%s", permissionType)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	// TODO display show result
	return nil
}

func cmdPermissionAdd(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	url.Path = "/permissions"

	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "type")
	permission := c.Args().Get(0)
	permissionType := c.Args().Get(2)

	var req proto.Request
	req.Permission = &proto.Permission{}
	req.Permission.Name = permission
	req.Permission.Category = permissionType

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func cmdPermissionDel(c *cli.Context) error {
	url := Cfg.Run.SomaAPI

	utl.ValidateCliArgumentCount(c, 1)
	permission := c.Args().Get(0)
	url.Path = fmt.Sprintf("/permissions/%s", permission)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func cmdPermissionList(c *cli.Context) error {
	url := Cfg.Run.SomaAPI
	url.Path = "/permissions"

	utl.ValidateCliArgumentCount(c, 0)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	// TODO list permissions
	return nil
}

func cmdPermissionShowGeneric(c *cli.Context, objType string) {
	url := Cfg.Run.SomaAPI
	var (
		objName string
		repo    string
		hasRepo bool
	)

	switch utl.GetCliArgumentCount(c) {
	case 1:
		hasRepo = false
	case 3:
		utl.ValidateCliArgument(c, 2, "repository")
		hasRepo = true
	default:
		log.Fatal("Syntax error, unexpected argument count")
	}

	objName = c.Args().Get(0)
	if hasRepo {
		repo = c.Args().Get(2)
		url.Path = fmt.Sprintf("/permissions/%s/%s/repository/%s",
			objType, objName, repo)
	} else {
		url.Path = fmt.Sprintf("/permissions/%s/%s", objType, objName)
	}

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	// TODO list permissions
}

func cmdPermissionShowUser(c *cli.Context) error {
	cmdPermissionShowGeneric(c, "user")
	return nil
}

func cmdPermissionShowTeam(c *cli.Context) error {
	cmdPermissionShowGeneric(c, "team")
	return nil
}

func cmdPermissionShowTool(c *cli.Context) error {
	cmdPermissionShowGeneric(c, "tool")
	return nil
}

func cmdPermissionShowPermission(c *cli.Context) error {
	url := Cfg.Run.SomaAPI

	utl.ValidateCliArgumentCount(c, 1)
	url.Path = fmt.Sprintf("/permissions/permission/%s", c.Args().Get(0))

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	// TODO list users
	return nil
}

func cmdPermissionAudit(c *cli.Context) error {
	url := Cfg.Run.SomaAPI

	utl.ValidateCliArgumentCount(c, 1)
	url.Path = fmt.Sprintf("/permissions/repository/%s", c.Args().Get(0))

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	// TODO list permissions
	return nil
}

func cmdPermissionGrantEnable(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	userId := utl.TryGetUserByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/permissions/user/%s", userId)

	resp := utl.PatchRequestWithBody(proto.Request{
		Grant: &proto.Grant{
			RecipientType: "user",
			RecipientId:   userId,
			Permission:    "global_grant_limited",
		},
	}, path,
	)
	fmt.Println(resp)
	return nil
}

func cmdPermissionGrantGlobal(c *cli.Context) error {
	/*
		url := getApiUrl()
		var objType string

		utl.ValidateCliArgumentCount(c, 3)
		switch c.Args().Get(1) {
		case "user":
			objType = "user"
		case "team":
			objType = "team"
		case "tool":
			objType = "tool"
		default:
			Slog.Fatal("Syntax error")
		}
		url.Path = fmt.Sprintf("/permissions/%s/%s", objType, c.Args().Get(2))

		var req somaproto.ProtoRequestPermission
		req.Permission = c.Args().Get(0)

		_, err := resty.New().
			SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
			R().
			SetBody(req).
			Patch(url.String())
		if err != nil {
			Slog.Fatal(err)
		}
	*/
	return nil
}

func cmdPermissionGrantLimited(c *cli.Context) error {
	/*
		url := getApiUrl()
		keys := [...]string{"repository", "bucket", "group", "cluster"}
		keySlice := make([]string, 0)
		var objType string
		var req somaproto.ProtoRequestPermission
		req.Grant.GrantType = "limited"

		// the keys array is ordered towars increasing detail, ie. there can
		// not be a limited grant on a cluster without specifying the
		// repository
		// Also, slicing uses halfopen interval [a:b), ie `a` is part of the
		// resulting slice, but `b` is not
		switch utl.GetCliArgumentCount(c) {
		case 5:
			keySlice = keys[0:1]
		case 7:
			keySlice = keys[0:2]
		case 9:
			keySlice = keys[0:3]
		case 11:
			keySlice = keys[:]
		default:
			Slog.Fatal("Syntax error, unexpected argument count")
		}
		// Tail() skips the first argument 0, which is returned by First(),
		// thus contains arguments 1-n. The first 3 arguments 0-2 are fixed in
		// order, the variable parts are arguments 4+ (argv 3-10) -- element
		// 2-9 in the tail slice.
		argSlice := c.Args().Tail()[2:]
		options := *parseLimitedGrantArguments(keySlice, argSlice)

		switch c.Args().Get(1) {
		case "user":
			objType = "user"
		case "team":
			objType = "team"
		case "tool":
			objType = "tool"
		default:
			Slog.Fatal("Syntax error")
		}
		req.Permission = c.Args().Get(0)
		url.Path = fmt.Sprintf("/permissions/%s/%s", objType, c.Args().Get(2))

		for k, v := range options {
			switch k {
			case "repository":
				req.Grant.Repository = v
			case "bucket":
				req.Grant.Bucket = v
			case "group":
				req.Grant.Group = v
			case "cluster":
				req.Grant.Cluster = v
			}
		}

		_, err := resty.New().
			SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
			R().
			SetBody(req).
			Patch(url.String())
		if err != nil {
			Slog.Fatal(err)
		}
	*/
	return nil
}

func cmdPermissionGrantSystem(c *cli.Context) error {
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
