package main

import (
	"fmt"

	"github.com/codegangsta/cli"
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
						Name:  "category",
						Usage: "SUBCOMMANDS for permission categories",
						Subcommands: []cli.Command{
							{
								Name:   "add",
								Usage:  "Register a new permission category",
								Action: runtime(cmdPermissionCategoryAdd),
							},
							{
								Name:   "remove",
								Usage:  "Remove an existing permission category",
								Action: runtime(cmdPermissionCategoryDel),
							},
							{
								Name:   "list",
								Usage:  "List all permission categories",
								Action: runtime(cmdPermissionCategoryList),
							},
							{
								Name:   "show",
								Usage:  "Show details for a permission category",
								Action: runtime(cmdPermissionCategoryShow),
							},
						}, // end permissions type
					},
					{
						Name:         "add",
						Usage:        "Register a new permission",
						Action:       runtime(cmdPermissionAdd),
						BashComplete: cmpl.PermissionAdd,
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
						Name:   "show",
						Usage:  "Show details for a permission",
						Action: runtime(cmdPermissionShow),
					},
					/*
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
					*/
				},
			}, // end permissions
		}...,
	)
	return &app
}

func cmdPermissionCategoryAdd(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.NewCategoryRequest()
	req.Category.Name = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, `/category/`)
	fmt.Println(resp)
	return nil
}

func cmdPermissionCategoryDel(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/category/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPermissionCategoryList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest(Client, `/category/`)
	fmt.Println(resp)
	return nil
}

func cmdPermissionCategoryShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/category/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPermissionAdd(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{`category`, `grants`}
	required := []string{`category`}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	req := proto.NewPermissionRequest()
	req.Permission.Name = c.Args().First()
	req.Permission.Category = opts[`category`][0]
	if sl, ok := opts[`grants`]; ok && len(sl) > 0 {
		req.Permission.Grants = opts[`grants`][0]
	}

	resp := utl.PostRequestWithBody(Client, req, `/permission/`)
	fmt.Println(resp)
	return nil
}

func cmdPermissionDel(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/permission/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPermissionList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest(Client, `/permission/`)
	fmt.Println(resp)
	return nil
}

func cmdPermissionShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	path := fmt.Sprintf("/permission/%s", c.Args().First())

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

/*
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
	userId := utl.TryGetUserByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/permissions/user/%s", userId)

	resp := utl.PatchRequestWithBody(Client, proto.Request{
		Grant: &proto.Grant{
			RecipientType: "user",
			RecipientId:   userId,
			PermissionId:  "global_grant_limited",
		},
	}, path,
	)
	fmt.Println(resp)
	return nil
}

func cmdPermissionGrantGlobal(c *cli.Context) error {
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
	return nil
}

func cmdPermissionGrantLimited(c *cli.Context) error {
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
	return nil
}

func cmdPermissionGrantSystem(c *cli.Context) error {
	return nil
}
*/

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
