package main

import (
	"fmt"
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
	unique := []string{
		"firstname",
		"lastname",
		"employeenr",
		"mailaddr",
		"team",
		"system"}
	required := []string{
		"firstname",
		"lastname",
		"employeenr",
		"mailaddr",
		"team"}
	var err error

	opts := map[string][]string{}
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	// validate
	if err = adm.ValidateEmployeeNumber(
		opts[`employeenr`][0]); err != nil {
		return err
	}
	if err = adm.ValidateMailAddress(
		opts[`mailaddr`][0]); err != nil {
		return err
	}

	req := proto.NewUserRequest()
	req.User.UserName = c.Args().First()
	req.User.FirstName = opts["firstname"][0]
	req.User.LastName = opts["lastname"][0]
	req.User.MailAddress = opts["mailaddr"][0]
	req.User.EmployeeNumber = opts["employeenr"][0]
	req.User.IsDeleted = false
	if req.User.TeamId, err = adm.LookupTeamId(
		opts["team"][0]); err != nil {
		return err
	}

	// optional arguments
	if _, ok := opts["system"]; ok {
		if err := adm.ValidateBool(opts["system"][0],
			&req.User.IsSystem); err != nil {
			return fmt.Errorf("Syntax error, system argument not"+
				" boolean: %s, %s", opts["active"][0], err.Error())
		}
	} else {
		req.User.IsSystem = false
	}

	return adm.Perform(`postbody`, `/users/`, `command`, req, c)
}

func cmdUserUpdate(c *cli.Context) error {
	unique := []string{
		`username`,
		"firstname",
		"lastname",
		"employeenr",
		"mailaddr",
		"team",
		`deleted`}
	required := []string{
		`username`,
		"firstname",
		"lastname",
		"employeenr",
		"mailaddr",
		"team",
		`deleted`}

	opts := map[string][]string{}
	var err error
	if err = adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	// validate
	if err = adm.ValidateEmployeeNumber(
		opts[`employeenr`][0]); err != nil {
		return err
	}
	if err = adm.ValidateMailAddress(
		opts[`mailaddr`][0]); err != nil {
		return err
	}
	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf(`users update requiress UUID as` +
			` first argument`)
	}

	req := proto.NewUserRequest()
	req.User.Id = c.Args().First()
	req.User.UserName = opts[`username`][0]
	req.User.FirstName = opts["firstname"][0]
	req.User.LastName = opts["lastname"][0]
	req.User.MailAddress = opts["mailaddr"][0]
	req.User.EmployeeNumber = opts["employeenr"][0]
	if req.User.TeamId, err = adm.LookupTeamId(
		opts[`team`][0]); err != nil {
		return err
	}
	if err = adm.ValidateBool(opts[`deleted`][0],
		&req.User.IsDeleted); err != nil {
		return fmt.Errorf("Syntax error, deleted argument not"+
			" boolean: %s, %s", opts[`deleted`][0], err.Error())
	}

	path := fmt.Sprintf("/users/%s", req.User.Id)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

func cmdUserMarkDeleted(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var (
		err    error
		userId string
	)
	if userId, err = adm.LookupUserId(c.Args().First()); err != nil {
		return err
	}

	path := fmt.Sprintf("/users/%s", userId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdUserPurgeDeleted(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var (
		err    error
		userId string
	)
	if userId, err = adm.LookupUserId(c.Args().First()); err != nil {
		return err
	}

	req := proto.Request{
		Flags: &proto.Flags{
			Purge: true,
		},
	}

	path := fmt.Sprintf("/users/%s", userId)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

func cmdUserActivate(c *cli.Context) error {
	// administrative use, full runtime is available
	if c.GlobalIsSet(`admin`) {
		if err := adm.VerifySingleArgument(c); err != nil {
			return err
		}
		return runtime(cmdUserActivateAdmin)(c)
	}
	// user trying to activate the account for the first
	// time, reduced runtime
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
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
		fmt.Printf("Starting with activation of account '%s' in"+
			" 2 seconds.\n", Cfg.Auth.User)
		fmt.Printf(`Use --user flag to activate a different account.`)
		time.Sleep(2 * time.Second)
	}
	if strings.Contains(Cfg.Auth.User, `:`) {
		return fmt.Errorf(`Usernames must not contain : character.`)
	}

	fmt.Printf("\nPlease provide the password you want to use.\n")
password_read:
	password = adm.ReadVerified(`password`)

	if happy, err = adm.EvaluatePassword(3, password,
		Cfg.Auth.User, `soma`); err != nil {
		return err
	} else if !happy {
		goto password_read
	}

	fmt.Printf("\nTo confirm that this is your account, an" +
		" additional credential is required" +
		" this once.\n")

	switch Cfg.Activation {
	case `ldap`:
		fmt.Println(`Please provide your LDAP password to"+
		" establish ownership.`)
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
	if err = adm.ValidateToken(Client, Cfg.Auth.User,
		cred.Token); err != nil {
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

func cmdUserList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/users/`, `list`, nil, c)
}

func cmdUserSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/users/`, `list`, nil, c)
}

func cmdUserShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	var (
		err error
		id  string
	)
	if id, err = adm.LookupUserId(c.Args().First()); err != nil {
		return err
	}

	path := fmt.Sprintf("/users/%s", id)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdUserPasswordUpdate(c *cli.Context) error {
	var (
		err               error
		password, passKey string
		happy             bool
		cred              *auth.Token
	)

	if Cfg.Auth.User == `` {
		fmt.Println(`Please specify for which  account the"+
		" password should be changed.`)
		if Cfg.Auth.User, err = adm.Read(`user`); err != nil {
			return err
		}
	} else {
		fmt.Printf("Starting with password update of account '%s'"+
			" in 2 seconds.\n", Cfg.Auth.User)
		fmt.Printf(`Use --user flag to switch account account.`)
		time.Sleep(2 * time.Second)
	}
	if strings.Contains(Cfg.Auth.User, `:`) {
		return fmt.Errorf(`Usernames must not contain : character.`)
	}

	fmt.Printf("\nPlease provide the new password you want to set.\n")
password_read:
	password = adm.ReadVerified(`password`)

	if happy, err = adm.EvaluatePassword(3, password,
		Cfg.Auth.User, `soma`); err != nil {
		return err
	} else if !happy {
		goto password_read
	}

	if c.Bool(`reset`) {
		fmt.Printf("\nTo confirm that you are allowed to reset" +
			" this account, an additional" +
			"credential is required.\n")

		switch Cfg.Activation {
		case `ldap`:
			fmt.Println(`Please provide your LDAP password to` +
				`establish ownership.`)
			passKey = adm.ReadVerified(`password`)
		case `mailtoken`:
			fmt.Println(`Please provide the token you received` +
				`via email.`)
			passKey = adm.ReadVerified(`token`)
		default:
			return fmt.Errorf(`Unknown activation mode`)
		}
	} else {
		fmt.Printf("\nPlease provide your currently active/old" +
			" password.\n")
		passKey = adm.ReadVerified(`password`)
	}

	if cred, err = adm.ChangeAccountPassword(
		Client,
		c.Bool(`reset`),
		&auth.Token{
			UserName: Cfg.Auth.User,
			Password: password,
			Token:    passKey,
		},
	); err != nil {
		return err
	}

	// validate received token
	if err = adm.ValidateToken(
		Client,
		Cfg.Auth.User,
		cred.Token,
	); err != nil {
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
