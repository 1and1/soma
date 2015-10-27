package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
	"gopkg.in/resty.v0"
)

func cmdUserAdd(c *cli.Context) {
	url := getApiUrl()
	url.Path = "/users"

	keySlice := []string{"firstname", "lastname", "employeenr",
		"mailaddr", "team", "active", "system"}
	reqSlice := []string{"firstname", "lastname", "employeenr",
		"mailaddr", "team"}
	var err error

	switch getCliArgumentCount(c) {
	case 11, 13, 15:
		break
	default:
		abort("Syntax error, unexpected argument count")
	}

	argSlice := c.Args().Tail()
	options, opts := parseVariableArguments(keySlice, reqSlice, argSlice)

	// validate
	validateStringAsEmployeeNumber(options["employeenr"])
	validateStringAsMailAddress(options["mailaddr"])

	var req somaproto.ProtoRequestUser
	req.User.UserName = c.Args().First()
	req.User.FirstName = options["firstname"]
	req.User.LastName = options["LastName"]
	req.User.Team = options["team"]
	req.User.MailAddress = options["mailaddr"]
	req.User.EmployeeNumber = options["employeenr"]
	req.User.IsDeleted = false

	// optional arguments
	if sliceContainsString("active", opts) {
		req.User.IsActive, err = strconv.ParseBool(options["active"])
		abortOnError(err, "Syntax error, active argument not boolean")
	} else {
		req.User.IsActive = true
	}

	if sliceContainsString("system", opts) {
		req.User.IsSystem, err = strconv.ParseBool(options["system"])
		abortOnError(err, "Syntax error, system argument not boolean")
	} else {
		req.User.IsSystem = false
	}

	_ = postRequestWithBody(req, url)
}

func cmdUserMarkDeleted(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	switch getCliArgumentCount(c) {
	case 1:
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Syntax error, argument not a uuid")
	case 2:
		validateCliArgument(c, 1, "by-name")
		id = getUserIdByName(c.Args().Get(1))
	default:
		abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/users/%s", id.String())

	_ = deleteRequest(url)
}

func cmdUserPurgeDeleted(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	if c.Bool("all") {
		validateCliArgumentCount(c, 0)
		url.Path = fmt.Sprintf("/users")
	} else {
		switch getCliArgumentCount(c) {
		case 1:
			id, err = uuid.FromString(c.Args().First())
			abortOnError(err, "Syntax error, argument not a uuid")
		case 2:
			validateCliArgument(c, 1, "by-name")
			id = getUserIdByName(c.Args().Get(1))
		default:
			abort("Syntax error, unexpected argument count")
		}
		url.Path = fmt.Sprintf("/users/%s", id.String)
	}

	var req somaproto.ProtoRequestUser
	req.Purge = true

	_ = deleteRequestWithBody(req, url)
}

func cmdUserRestoreDeleted(c *cli.Context) {
	url := getApiUrl()
	var (
		id  uuid.UUID
		err error
	)

	if c.Bool("all") {
		validateCliArgumentCount(c, 0)
		url.Path = fmt.Sprintf("/users")
	} else {
		switch getCliArgumentCount(c) {
		case 1:
			id, err = uuid.FromString(c.Args().First())
			abortOnError(err, "Syntax error, argument not a uuid")
		case 2:
			validateCliArgument(c, 1, "by-name")
			id = getUserIdByName(c.Args().Get(1))
		default:
			abort("Syntax error, unexpected argument count")
		}
		url.Path = fmt.Sprintf("/users/%s", id.String())
	}

	var req somaproto.ProtoRequestUser
	req.Restore = true

	_ = patchRequestWithBody(req, url)
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

	switch getCliArgumentCount(c) {
	case 1, 3, 5, 7, 9, 11:
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Syntax error, argument not a uuid")
		argSlice = c.Args().Tail()
	case 2, 4, 6, 8, 10, 12:
		validateCliArgument(c, 1, "by-name")
		id = getUserIdByName(c.Args().Tail()[0])
		argSlice = c.Args().Tail()[1:]
	default:
		abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/users/%s", id.String())

	options, opts := parseVariableArguments(keySlice, reqSlice, argSlice)
	var req somaproto.ProtoRequestUser

	for _, v := range opts {
		switch v {
		case "firstname":
			req.User.FirstName = options["firstname"]
		case "lastname":
			req.User.LastName = options["lastname"]
		case "employeenr":
			validateStringAsEmployeeNumber(options["employeenr"])
			req.User.EmployeeNumber = options["employeenr"]
		case "mailaddr":
			validateStringAsMailAddress(options["mailaddr"])
			req.User.MailAddress = options["mailaddr"]
		case "team":
			req.User.Team = options["team"]
		}
	}

	_ = patchRequestWithBody(req, url)
}

func cmdUserRename(c *cli.Context) {
	url := getApiUrl()
	var (
		id      uuid.UUID
		err     error
		newName string
	)

	switch getCliArgumentCount(c) {
	case 3:
		validateCliArgument(c, 2, "to")
		id, err = uuid.FromString(c.Args().First())
		abortOnError(err, "Syntax error, argument not a uuid")
		newName = c.Args().Get(2)
	case 4:
		validateCliArgument(c, 1, "by-name")
		validateCliArgument(c, 3, "to")
		id = getUserIdByName(c.Args().Get(1))
		newName = c.Args().Get(3)
	default:
		abort("Syntax error, unexpected argument count")
	}
	url.Path = fmt.Sprintf("/users/%s", id.String())

	var req somaproto.ProtoRequestUser
	req.User.UserName = newName

	_ = patchRequestWithBody(req, url)
}

func cmdUserActivate(c *cli.Context) {
	abort("Not implemented")
}

func cmdUserDeactivate(c *cli.Context) {
	abort("Not implemented")
}

func cmdUserList(c *cli.Context) {
	abort("Not implemented")
}

func cmdUserShow(c *cli.Context) {
	abort("Not implemented")
}

func cmdUserPasswordUpdate(c *cli.Context) {
	abort("Not implemented")
}

func cmdUserPasswordReset(c *cli.Context) {
	abort("Not implemented")
}

func cmdUserPasswordForce(c *cli.Context) {
	abort("Not implemented")
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
