package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/internal/help"
	"github.com/1and1/soma/lib/auth"
	"github.com/1and1/soma/lib/proto"
	"github.com/codegangsta/cli"
)

func registerUsers(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// users
			{
				Name:  "users",
				Usage: "SUBCOMMANDS for users",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new user",
						Action:       runtime(cmdUserAdd),
						BashComplete: cmpl.UserAdd,
					},
					{
						Name:   "delete",
						Usage:  "Mark a user as deleted",
						Action: runtime(cmdUserMarkDeleted),
					},
					{
						Name:   "purge",
						Usage:  "Purge a user marked as deleted",
						Action: runtime(cmdUserPurgeDeleted),
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "all, a",
								Usage: "Purge all deleted users",
							},
						},
					},
					{
						Name:         "update",
						Usage:        "Set/change user information",
						Action:       runtime(cmdUserUpdate),
						BashComplete: cmpl.UserUpdate,
					},
					/*
						{
							Name:   "restore",
							Usage:  "Restore a user marked as deleted",
							Action: cmdUserRestoreDeleted,
							Flags: []cli.Flag{
								cli.BoolFlag{
									Name:  "all, a",
									Usage: "Restore all deleted users",
								},
							},
						},
						{
							Name:   "rename",
							Usage:  "Change a user's username",
							Action: cmdUserRename,
						},
					*/
					{
						Name:   "activate",
						Usage:  "Activate a deativated user",
						Action: cmdUserActivate,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "force, f",
								Usage: "Apply administrative force to the activation",
							},
						},
					},
					/*
						{
							Name:   "deactivate",
							Usage:  "Deactivate a user account",
							Action: cmdUserDeactivate,
						},
					*/
					{
						Name:  `password`,
						Usage: "SUBCOMMANDS for user passwords",
						Subcommands: []cli.Command{
							{
								Name:        `update`,
								Usage:       `Update the password of one's own user account`,
								Action:      boottime(cmdUserPasswordUpdate),
								Description: help.Text(`UsersPasswordUpdate`),
								Flags: []cli.Flag{
									cli.BoolFlag{
										Name:  `reset, r`,
										Usage: `Reset the password via activation credentials`,
									},
								},
							},
							/*
								{
									Name:   `reset`,
									Usage:  `Trigger a password reset for a user`,
									Action: cmdUserPasswordReset,
								},
								{
									Name:   `force`,
									Usage:  `Forcefully set the password of a user`,
									Action: cmdUserPasswordForce,
								},
							*/
						},
					}, // end users password
					{
						Name:   "list",
						Usage:  "List all registered users",
						Action: runtime(cmdUserList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific user",
						Action: runtime(cmdUserShow),
					},
					{
						Name:   "synclist",
						Usage:  "List all registered users suitable for sync",
						Action: runtime(cmdUserSync),
					},
				},
			}, // end users
		}...,
	)
	return &app
}

func cmdUserAdd(c *cli.Context) error {
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

	req := proto.Request{}
	req.User = &proto.User{}
	req.User.UserName = c.Args().First()
	req.User.FirstName = opts["firstname"][0]
	req.User.LastName = opts["lastname"][0]
	req.User.TeamId = utl.TryGetTeamByUUIDOrName(Client, opts["team"][0])
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

	resp := utl.PostRequestWithBody(Client, req, "/users/")
	fmt.Println(resp)
	return nil
}

func cmdUserUpdate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 15)
	multiple := []string{}
	unique := []string{`username`, "firstname", "lastname", "employeenr",
		"mailaddr", "team", `deleted`}
	required := []string{`username`, "firstname", "lastname", "employeenr",
		"mailaddr", "team", `deleted`}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	// validate
	utl.ValidateStringAsEmployeeNumber(opts["employeenr"][0])
	utl.ValidateStringAsMailAddress(opts["mailaddr"][0])
	if !utl.IsUUID(c.Args().First()) {
		return fmt.Errorf(`users update requiress UUID as first argument`)
	}

	req := proto.NewUserRequest()
	req.User.Id = c.Args().First()
	req.User.UserName = opts[`username`][0]
	req.User.FirstName = opts["firstname"][0]
	req.User.LastName = opts["lastname"][0]
	req.User.TeamId = utl.TryGetTeamByUUIDOrName(Client, opts["team"][0])
	req.User.MailAddress = opts["mailaddr"][0]
	req.User.EmployeeNumber = opts["employeenr"][0]
	req.User.IsDeleted = utl.GetValidatedBool(opts[`deleted`][0])

	path := fmt.Sprintf("/users/%s", req.User.Id)
	resp := utl.PutRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdUserMarkDeleted(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	userId := utl.TryGetUserByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/users/%s", userId)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdUserPurgeDeleted(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	userId := utl.TryGetUserByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/users/%s", userId)

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	resp := utl.DeleteRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
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

	_ = utl.PatchRequestWithBody(Client, req, url.String())
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

	_ = utl.PatchRequestWithBody(Client, req, url.String())
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

	_ = utl.PatchRequestWithBody(Client, req, url.String())
}
*/

func cmdUserActivate(c *cli.Context) error {
	// administrative use, full runtime is available
	if c.GlobalIsSet(`admin`) {
		utl.ValidateCliArgumentCount(c, 1)
		return runtime(cmdUserActivateAdmin)(c)
	}
	// user trying to activate the account for the first
	// time, reduced runtime
	utl.ValidateCliArgumentCount(c, 0)
	return boottime(cmdUserActivateUser)(c)
}

func cmdUserActivateUser(c *cli.Context) error {
	var err error
	var password string
	var passKey string
	var happy bool
	var cred *auth.Token

	if Cfg.Auth.User == "" {
		fmt.Println(`Please specify which account to activate.`)
		if Cfg.Auth.User, err = adm.Read(`user`); err != nil {
			return err
		}
	} else {
		fmt.Printf("Starting with activation of account '%s' in 2 seconds.\n", Cfg.Auth.User)
		fmt.Printf(`Use --user flag to activate a different account.`)
		time.Sleep(2 * time.Second)
	}
	if strings.Contains(Cfg.Auth.User, `:`) {
		return fmt.Errorf(`Usernames must not contain : character.`)
	}

	fmt.Printf("\nPlease provide the password you want to use.\n")
password_read:
	password = adm.ReadVerified(`password`)

	if happy, err = adm.EvaluatePassword(3, password, Cfg.Auth.User, `soma`); err != nil {
		return err
	} else if !happy {
		password = ""
		goto password_read
	}

	fmt.Printf("\nTo confirm that this is your account, an additional credential is required" +
		" this once.\n")

	switch Cfg.Activation {
	case `ldap`:
		fmt.Println(`Please provide your LDAP password to establish ownership.`)
		passKey = adm.ReadVerified(`password`)
	case `mailtoken`:
		fmt.Println(`Please provide the token you received via email.`)
		passKey = adm.ReadVerified(`token`)
	default:
		return fmt.Errorf(`Unknown activation mode`)
	}

	if cred, err = adm.ActivateAccount(Client, &auth.Token{
		UserName: Cfg.Auth.User,
		Password: password,
		Token:    passKey,
	}); err != nil {
		return err
	}

	// validate received token
	if err = adm.ValidateToken(Client, Cfg.Auth.User, cred.Token); err != nil {
		return err
	}
	// save received token
	if err = store.SaveToken(
		Cfg.Auth.User,
		cred.ValidFrom,
		cred.ExpiresAt,
		cred.Token,
	); err != nil {
		return err
	}
	return nil
}

func cmdUserActivateAdmin(c *cli.Context) error {
	return nil
}

/*
func cmdUserDeactivate(c *cli.Context) {
	url := getApiUrl()
	id := utl.UserIdByUuidOrName(c)
	url.Path = fmt.Sprintf("/users/%s", id.String())

	var req somaproto.ProtoRequestUser
	req.User.IsActive = false

	_ = utl.PatchRequestWithBody(Client, req, url.String())
}
*/

func cmdUserList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest(Client, "/users/")
	fmt.Println(resp)
	return nil
}

func cmdUserSync(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest(Client, `/sync/users/`)
	fmt.Println(resp)
	return nil
}

func cmdUserShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	id := utl.TryGetUserByUUIDOrName(Client, c.Args().First())
	path := fmt.Sprintf("/users/%s", id)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

/*
func cmdUserPasswordReset(c *cli.Context) {
	id := utl.UserIdByUuidOrName(c)
	path := fmt.Sprintf("/users/%s/password", id.String())

	var req somaproto.ProtoRequestUser
	req.Credentials.Reset = true

	_ = utl.PutRequestWithBody(Client, req, path)
}

func cmdUserPasswordForce(c *cli.Context) {
	id := utl.UserIdByUuidOrName(c)
	path := fmt.Sprintf("/users/%s/password", id.String())
	pass := utl.GetNewPassword()

	var req somaproto.ProtoRequestUser
	req.Credentials.Force = true
	req.Credentials.Password = pass

	_ = utl.PutRequestWithBody(Client, req, path)
}
*/

func cmdUserPasswordUpdate(c *cli.Context) error {
	var (
		err               error
		password, passKey string
		happy             bool
		cred              *auth.Token
	)

	if Cfg.Auth.User == `` {
		fmt.Println(`Please specify for which  account the password should be changed.`)
		if Cfg.Auth.User, err = adm.Read(`user`); err != nil {
			return err
		}
	} else {
		fmt.Printf("Starting with password update of account '%s' in 2 seconds.\n", Cfg.Auth.User)
		fmt.Printf(`Use --user flag to switch account account.`)
		time.Sleep(2 * time.Second)
	}
	if strings.Contains(Cfg.Auth.User, `:`) {
		return fmt.Errorf(`Usernames must not contain : character.`)
	}

	fmt.Printf("\nPlease provide the new password you want to set.\n")
password_read:
	password = adm.ReadVerified(`password`)

	if happy, err = adm.EvaluatePassword(3, password, Cfg.Auth.User, `soma`); err != nil {
		return err
	} else if !happy {
		password = ``
		goto password_read
	}

	if c.Bool(`reset`) {
		fmt.Printf("\nTo confirm that you are allowed to reset this account, an additional" +
			"credential is required.\n")

		switch Cfg.Activation {
		case `ldap`:
			fmt.Println(`Please provide your LDAP password to establish ownership.`)
			passKey = adm.ReadVerified(`password`)
		case `mailtoken`:
			fmt.Println(`Please provide the token you received via email.`)
			passKey = adm.ReadVerified(`token`)
		default:
			return fmt.Errorf(`Unknown activation mode`)
		}
	} else {
		fmt.Printf("\nPlease provide your currently active/old password.\n")
		passKey = adm.ReadVerified(`password`)
	}

	if cred, err = adm.ChangeAccountPassword(Client, c.Bool(`reset`), &auth.Token{
		UserName: Cfg.Auth.User,
		Password: password,
		Token:    passKey,
	}); err != nil {
		return err
	}

	// validate received token
	if err = adm.ValidateToken(Client, Cfg.Auth.User, cred.Token); err != nil {
		return err
	}
	// save received token
	if err = store.SaveToken(
		Cfg.Auth.User,
		cred.ValidFrom,
		cred.ExpiresAt,
		cred.Token,
	); err != nil {
		return err
	}
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
