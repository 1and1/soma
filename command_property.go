package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

/*
 * CREATE
 */
func cmdPropertyCustomCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	var req somaproto.ProtoRequestProperty
	req.Custom.Property = c.Args().First()
	req.Custom.Repository = c.Args().Get(2)
	path := fmt.Sprintf("/property/custom/%s/", repoId.String())

	_ = utl.PostRequestWithBody(req, path)
}

func cmdPropertySystemCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	var req somaproto.ProtoRequestProperty
	req.System.Property = c.Args().First()

	_ = utl.PostRequestWithBody(req, "/property/system/")
}

func cmdPropertyServiceCreate(c *cli.Context) {
	multKeys := []string{"transport", "application", "port", "process",
		"file", "directory", "socket", "uid"}
	uniqKeys := []string{"tls", "provider", "team"}
	reqKeys := []string{"team"}

	argCount := utl.GetCliArgumentCount(c)
	switch {
	// first argument is the name of template, then attributes and
	// values are added in pairs of two -> valid are 3,5,7,... args
	case argCount == 0 || argCount == 1:
		utl.Abort("Syntax error, unexpected argument count")
	case (argCount % 2) == 0:
		break
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	argSlice := c.Args().Tail()

	opts := utl.ParseVariadicArguments(multKeys, uniqKeys,
		reqKeys, argSlice)

	teamId := utl.TryGetTeamByUUIDOrName(opts["team"][0])

	var req somaproto.ProtoRequestProperty
	req.Service.Property = c.Args().First()
	req.Service.Team = opts["team"][0] // required argument
	for optKey, optVal := range opts {
		switch optKey {
		case "transport":
			req.Service.Attributes.Transport = optVal
		case "application":
			req.Service.Attributes.Application = optVal
		case "port":
			req.Service.Attributes.Port = optVal
		case "process":
			req.Service.Attributes.Process = optVal
		case "file":
			req.Service.Attributes.File = optVal
		case "directory":
			req.Service.Attributes.Directory = optVal
		case "socket":
			req.Service.Attributes.Socket = optVal
		case "uid":
			req.Service.Attributes.Uid = optVal
		case "provider":
			if optVal != nil {
				req.Service.Attributes.Tls = optVal[0]
			}
		case "tls":
			if optVal != nil {
				req.Service.Attributes.Provider = optVal[0]
			}
		default:
			utl.Abort("Error assigning service attributes")
		}
	}
	path := fmt.Sprintf("/property/service/team/%s/", teamId)

	_ = utl.PostRequestWithBody(req, path)
}

func cmdPropertyTemplateCreate(c *cli.Context) {
	multKeys := []string{"transport", "application", "port", "process",
		"file", "directory", "socket", "uid"}
	uniqKeys := []string{"tls", "provider"}
	reqKeys := []string{}

	argCount := utl.GetCliArgumentCount(c)
	switch {
	// first argument is the name of template, then attributes and
	// values are added in pairs of two -> valid are 1,3,5,7,... args
	case argCount == 0:
		utl.Abort("Syntax error, unexpected argument count")
	case (argCount % 2) == 0:
		break
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	argSlice := c.Args().Tail()

	opts := utl.ParseVariadicArguments(multKeys, uniqKeys,
		reqKeys, argSlice)

	var req somaproto.ProtoRequestProperty
	req.Service.Property = c.Args().First()
	for optKey, optVal := range opts {
		switch optKey {
		case "transport":
			req.Service.Attributes.Transport = optVal
		case "application":
			req.Service.Attributes.Application = optVal
		case "port":
			req.Service.Attributes.Port = optVal
		case "process":
			req.Service.Attributes.Process = optVal
		case "file":
			req.Service.Attributes.File = optVal
		case "directory":
			req.Service.Attributes.Directory = optVal
		case "socket":
			req.Service.Attributes.Socket = optVal
		case "uid":
			req.Service.Attributes.Uid = optVal
		case "provider":
			if optVal != nil {
				req.Service.Attributes.Tls = optVal[0]
			}
		case "tls":
			if optVal != nil {
				req.Service.Attributes.Provider = optVal[0]
			}
		default:
			utl.Abort("Error assigning service attributes")
		}
	}

	_ = utl.PostRequestWithBody(req, "/property/service/global/")
}

/*
 * DELETE
 */
func cmdPropertyCustomDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetCustomPropertyByUUIDOrName(c.Args().Get(0), repoId.String())
	path := fmt.Sprintf("/property/custom/%s/%s", repoId.String(), propId.String())

	_ = utl.DeleteRequest(path)
}

func cmdPropertySystemDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetSystemPropertyByUUIDOrName(c.Args().Get(0))
	path := fmt.Sprintf("/property/system/%s", propId.String())

	_ = utl.DeleteRequest(path)
}

func cmdPropertyServiceDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")
	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetServicePropertyByUUIDOrName(c.Args().Get(0), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId.String())

	_ = utl.DeleteRequest(path)
}

func cmdPropertyTemplateDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetTemplatePropertyByUUIDOrName(c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId.String())

	_ = utl.DeleteRequest(path)
}

/*
 * EDIT
 */
func cmdPropertyServiceEdit(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateEdit(c *cli.Context) {
	utl.NotImplemented()
}

/*
 * RENAME
 */
func cmdPropertyCustomRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertySystemRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyServiceRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateRename(c *cli.Context) {
	utl.NotImplemented()
}

/*
 * SHOW
 */
func cmdPropertyCustomShow(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertySystemShow(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyServiceShow(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateShow(c *cli.Context) {
	utl.NotImplemented()
}

/*
 * LIST
 */
func cmdPropertyCustomList(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertySystemList(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyServiceList(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateList(c *cli.Context) {
	utl.NotImplemented()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
