package main

import (
	"fmt"
	"strconv"

	"github.com/codegangsta/cli"
)

func cmdUserAdd(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 11)
	multiple := []string{}
	unique := []string{"firstname", "lastname", "employeenr",
		"mailaddr", "team", "active", "system"}
	required := []string{"firstname", "lastname", "employeenr",
		"mailaddr", "team"}
	var err error

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	// validate
	utl.ValidateStringAsEmployeeNumber(opts["employeenr"][0])
	utl.ValidateStringAsMailAddress(opts["mailaddr"][0])

	req := somaproto.ProtoRequestUser{}
	req.User = &somaproto.ProtoUser{}
	req.User.UserName = c.Args().First()
	req.User.FirstName = opts["firstname"][0]
	req.User.LastName = opts["lastname"][0]
	req.User.Team = utl.TryGetTeamByUUIDOrName(opts["team"][0])
	req.User.MailAddress = opts["mailaddr"][0]
	req.User.EmployeeNumber = opts["employeenr"][0]
	req.User.IsDeleted = false

	// optional arguments
	if _, ok := opts["active"]; ok {
		req.User.IsActive, err = strconv.ParseBool(opts["active"][0])
		utl.AbortOnError(err, "Syntax error, active argument not boolean")
	} else {
		req.User.IsActive = true
	}

	if _, ok := opts["system"]; ok {
		req.User.IsSystem, err = strconv.ParseBool(opts["system"][0])
		utl.AbortOnError(err, "Syntax error, system argument not boolean")
	} else {
		req.User.IsSystem = false
	}

	resp := utl.PostRequestWithBody(req, "/users/")
	fmt.Println(resp)
}

func cmdUserMarkDeleted(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	userId := utl.TryGetUserByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/users/%s", userId)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdUserPurgeDeleted(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	userId := utl.TryGetUserByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/users/%s", userId)

	req := somaproto.ProtoRequestUser{}
	req.Purge = true

	resp := utl.DeleteRequestWithBody(req, path)
	fmt.Println(resp)
}

/*
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
*/

func cmdUserList(c *cli.Context) {
	resp := utl.GetRequest("/users/")
	fmt.Println(resp)
}

func cmdUserShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetUserByUUIDOrName(c.Args().First())
	path := fmt.Sprintf("/users/%s", id)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

/*
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
*/

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
