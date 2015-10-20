package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
)

func cmdPermissionTypeAdd(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/permissions/types"

	validateCliArgumentCount(c, 1)
	permissionType := c.Args().First()
	Slog.Printf("Command: add permission type [%s]", permissionType)

	var req somaproto.ProtoRequestPermission
	req.PermissionType = permissionType

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
	Slog.Printf("Response: %s\n", resp.Status())
}

func cmdPermissionTypeDel(c *cli.Context) {
	url := getApiUrl()

	validateCliArgumentCount(c, 1)
	permissionType := c.Args().First()
	url.Path = fmt.Sprintf("/permissions/types/%s", permissionType)
	Slog.Printf("Command: delete permission type [%s]", permissionType)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
}

func cmdPermissionTypeRename(c *cli.Context) {
	url := getApiUrl()

	validateCliArgumentCount(c, 3)
	validateCliArgument(c, 2, "to") // starts args counting at 1
	permissionType := c.Args().Get(0)
	newPermissionType := c.Args().Get(2)
	url.Path = fmt.Sprintf("/permissions/types/%s", permissionType)

	var req somaproto.ProtoRequestPermission
	req.PermissionType = newPermissionType

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
}

func cmdPermissionTypeList(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/permissions/types"

	validateCliArgumentCount(c, 0)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
	// TODO display list result
}

func cmdPermissionTypeShow(c *cli.Context) {
	url := getApiUrl()

	validateCliArgumentCount(c, 1)
	permissionType := c.Args().Get(0)
	url.Path = fmt.Sprintf("/permissions/types/%s", permissionType)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
	// TODO display show result
}

func cmdPermissionAdd(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/permissions"

	validateCliArgumentCount(c, 3)
	validateCliArgument(c, 2, "type")
	permission := c.Args().Get(0)
	permissionType := c.Args().Get(2)

	var req somaproto.ProtoRequestPermission
	req.Permission = permission
	req.PermissionType = permissionType

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
}

func cmdPermissionDel(c *cli.Context) {
	url := getApiUrl()

	validateCliArgumentCount(c, 1)
	permission := c.Args().Get(0)
	url.Path = fmt.Sprintf("/permissions/%s", permission)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
}

func cmdPermissionList(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/permissions"

	validateCliArgumentCount(c, 0)

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
	// TODO list permissions
}

func cmdPermissionShowGeneric(c *cli.Context, objType string) {
	url := getApiUrl()
	var (
		objName string
		repo    string
		hasRepo bool
	)

	switch getCliArgumentCount(c) {
	case 1:
		hasRepo = false
	case 3:
		validateCliArgument(c, 2, "repository")
		hasRepo = true
	default:
		Slog.Fatal("Syntax error, unexpected argument count")
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
		Slog.Fatal(err)
	}
	// TODO list permissions
}

func cmdPermissionShowUser(c *cli.Context) {
	cmdPermissionShowGeneric(c, "user")
}

func cmdPermissionShowTeam(c *cli.Context) {
	cmdPermissionShowGeneric(c, "team")
}

func cmdPermissionShowTool(c *cli.Context) {
	cmdPermissionShowGeneric(c, "tool")
}

func cmdPermissionShowPermission(c *cli.Context) {
	url := getApiUrl()

	validateCliArgumentCount(c, 1)
	url.Path = fmt.Sprintf("/permissions/permission/%s", c.Args().Get(0))

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
	// TODO list users
}

func cmdPermissionAudit(c *cli.Context) {
	url := getApiUrl()

	validateCliArgumentCount(c, 1)
	url.Path = fmt.Sprintf("/permissions/repository/%s", c.Args().Get(0))

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
	// TODO list permissions
}

func cmdPermissionGrantEnable(c *cli.Context) {
	url := getApiUrl()

	validateCliArgumentCount(c, 1)
	url.Path = fmt.Sprintf("/permissions/user/%s", c.Args().Get(0))

	var req somaproto.ProtoRequestPermission
	req.GrantEnabled = true

	_, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Patch(url.String())
	if err != nil {
		Slog.Fatal(err)
	}
	//
}

func cmdPermissionGrantGlobal(c *cli.Context) {
	url := getApiUrl()
	var objType string

	validateCliArgumentCount(c, 3)
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
}

func cmdPermissionGrantLimited(c *cli.Context) {
	url := getApiUrl()
	keys := [...]string{"repository", "bucket", "group", "cluster"}
	keySlice := make([]string, 0)
	var objType string
	var req somaproto.ProtoRequestPermission
	req.Grant.GrantType = "limited"

	// the keys array is ordered towars increasing detail, ie. there can
	// not be a limited grant on a cluster without specifying the
	// repository
	switch getCliArgumentCount(c) {
	case 5:
		keySlice = keys[1:1]
	case 7:
		keySlice = keys[1:2]
	case 9:
		keySlice = keys[1:3]
	case 11:
		keySlice = keys[1:4]
	default:
		Slog.Fatal("Syntax error, unexpected argument count")
	}
	// Tail() skips the first argument, which is returned by First(),
	// thus contains arguments 2-n. The first 3 arguments are fixed in
	// order, the variable parts are arguments 4-n -- element 3-n in the
	// tail slice. Slice element numbering starts at 1...
	argSlice := c.Args().Tail()[3:]
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
}

func cmdPermissionGrantSystem(c *cli.Context) {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
