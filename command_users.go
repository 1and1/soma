package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
)

func cmdUserAdd(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/users"

	keySlice := []string{"firstname", "lastname", "employeenr",
		"mailaddr", "team", "active", "system"}
	reqSlice := []string{"firstname", "lastname", "employeenr",
		"mailaddr", "team"}
	var err error

	switch utl.GetCliArgumentCount(c) {
	case 11, 13, 15:
		break
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}

	argSlice := c.Args().Tail()
	options, opts := utl.ParseVariableArguments(keySlice, reqSlice, argSlice)

	// validate
	utl.ValidateStringAsEmployeeNumber(options["employeenr"])
	utl.ValidateStringAsMailAddress(options["mailaddr"])

	var req somaproto.ProtoRequestUser
	req.User.UserName = c.Args().First()
	req.User.FirstName = options["firstname"]
	req.User.LastName = options["LastName"]
	req.User.Team = options["team"]
	req.User.MailAddress = options["mailaddr"]
	req.User.EmployeeNumber = options["employeenr"]
	req.User.IsDeleted = false

	// optional arguments
	if utl.SliceContainsString("active", opts) {
		req.User.IsActive, err = strconv.ParseBool(options["active"])
		utl.AbortOnError(err, "Syntax error, active argument not boolean")
	} else {
		req.User.IsActive = true
	}

	if utl.SliceContainsString("system", opts) {
		req.User.IsSystem, err = strconv.ParseBool(options["system"])
		utl.AbortOnError(err, "Syntax error, system argument not boolean")
	} else {
		req.User.IsSystem = false
	}

	_ = utl.PostRequestWithBody(req, url.String())
}

func cmdUserMarkDeleted(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	switch utl.GetCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Syntax error, argument not a uuid")
	case 2:
		utl.ValidateCliArgument(c, 1, "by-name")
		id = utl.GetUserIdByName(c.Args().Get(1))
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/users/%s", id.String())

	_ = utl.DeleteRequest(url.String())
}

func cmdUserPurgeDeleted(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		url.Path = fmt.Sprintf("/users")
	} else {
		switch utl.GetCliArgumentCount(c) {
		case 1:
			id, err = uuid.FromString(c.Args().First())
			utl.AbortOnError(err, "Syntax error, argument not a uuid")
		case 2:
			utl.ValidateCliArgument(c, 1, "by-name")
			id = utl.GetUserIdByName(c.Args().Get(1))
		default:
			utl.Abort("Syntax error, unexpected argument count")
		}
		url.Path = fmt.Sprintf("/users/%s", id.String)
	}

	var req somaproto.ProtoRequestUser
	req.Purge = true

	_ = utl.DeleteRequestWithBody(req, url.String())
}

func cmdUserRestoreDeleted(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	if c.Bool("all") {
		utl.ValidateCliArgumentCount(c, 0)
		url.Path = fmt.Sprintf("/users")
	} else {
		switch utl.GetCliArgumentCount(c) {
		case 1:
			id, err = uuid.FromString(c.Args().First())
			utl.AbortOnError(err, "Syntax error, argument not a uuid")
		case 2:
			utl.ValidateCliArgument(c, 1, "by-name")
			id = utl.GetUserIdByName(c.Args().Get(1))
		default:
			utl.Abort("Syntax error, unexpected argument count")
		}
		url.Path = fmt.Sprintf("/users/%s", id.String())
	}

	var req somaproto.ProtoRequestUser
	req.Restore = true

	_ = utl.PatchRequestWithBody(req, url.String())
}

func cmdUserUpdate(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	argSlice := make([]string, 0)
	keySlice := []string{"firstname", "lastname", "employeenr", "mailaddr", "team"}
	reqSlice := make([]string, 0)

	switch utl.GetCliArgumentCount(c) {
	case 1, 3, 5, 7, 9, 11:
		id, err = uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Syntax error, argument not a uuid")
		argSlice = c.Args().Tail()
	case 2, 4, 6, 8, 10, 12:
		utl.ValidateCliArgument(c, 1, "by-name")
		id = utl.GetUserIdByName(c.Args().Tail()[0])
		argSlice = c.Args().Tail()[1:]
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/users/%s", id.String())

	options, opts := utl.ParseVariableArguments(keySlice, reqSlice, argSlice)
	var req somaproto.ProtoRequestUser

	for _, v := range opts {
		switch v {
		case "firstname":
			req.User.FirstName = options["firstname"]
		case "lastname":
			req.User.LastName = options["lastname"]
		case "employeenr":
			utl.ValidateStringAsEmployeeNumber(options["employeenr"])
			req.User.EmployeeNumber = options["employeenr"]
		case "mailaddr":
			utl.ValidateStringAsMailAddress(options["mailaddr"])
			req.User.MailAddress = options["mailaddr"]
		case "team":
			req.User.Team = options["team"]
		}
	}

	_ = utl.PatchRequestWithBody(req, url.String())
}

func cmdUserRename(c *cli.Context) {
	url := getApiUrl()
	var (
		id      uuid.UUID
		err     error
		newName string
	)

	switch utl.GetCliArgumentCount(c) {
	case 3:
		utl.ValidateCliArgument(c, 2, "to")
		id, err = uuid.FromString(c.Args().First())
		utl.AbortOnError(err, "Syntax error, argument not a uuid")
		newName = c.Args().Get(2)
	case 4:
		utl.ValidateCliArgument(c, 1, "by-name")
		utl.ValidateCliArgument(c, 3, "to")
		id = utl.GetUserIdByName(c.Args().Get(1))
		newName = c.Args().Get(3)
	default:
		utl.Abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/users/%s", id.String())

	var req somaproto.ProtoRequestUser
	req.User.UserName = newName

	_ = utl.PatchRequestWithBody(req, url.String())
}

func cmdUserActivate(c *cli.Context) {
	url := getApiUrl()
	id := utl.UserIdByUuidOrName(c)
	url.Path = fmt.Sprintf("/users/%s", id.String())

	var req somaproto.ProtoRequestUser
	req.User.IsActive = true

	_ = utl.PatchRequestWithBody(req, url.String())
}

func cmdUserDeactivate(c *cli.Context) {
	url := getApiUrl()
	id := utl.UserIdByUuidOrName(c)
	url.Path = fmt.Sprintf("/users/%s", id.String())

	var req somaproto.ProtoRequestUser
	req.User.IsActive = false

	_ = utl.PatchRequestWithBody(req, url.String())
}

func cmdUserList(c *cli.Context) {
	_ = utl.GetRequest("/users")
}

func cmdUserShow(c *cli.Context) {
	id := utl.UserIdByUuidOrName(c)
	path := fmt.Sprintf("/users/%s", id.String())

	_ = utl.GetRequest(path)
}

func cmdUserPasswordUpdate(c *cli.Context) {
	id := utl.UserIdByUuidOrName(c)
	path := fmt.Sprintf("/users/%s/password", id.String())
	pass := utl.GetNewPassword()

	var req somaproto.ProtoRequestUser
	req.Credentials.Password = pass

	_ = utl.PutRequestWithBody(req, path)
}

func cmdUserPasswordReset(c *cli.Context) {
	id := utl.UserIdByUuidOrName(c)
	path := fmt.Sprintf("/users/%s/password", id.String())

	var req somaproto.ProtoRequestUser
	req.Credentials.Reset = true

	_ = utl.PutRequestWithBody(req, path)
}

func cmdUserPasswordForce(c *cli.Context) {
	id := utl.UserIdByUuidOrName(c)
	path := fmt.Sprintf("/users/%s/password", id.String())
	pass := utl.GetNewPassword()

	var req somaproto.ProtoRequestUser
	req.Credentials.Force = true
	req.Credentials.Password = pass

	_ = utl.PutRequestWithBody(req, path)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
