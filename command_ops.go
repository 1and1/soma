package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/resty.v0"

	"github.com/codegangsta/cli"
)

func registerOps(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:   "ops",
				Usage:  "SUBCOMMANDS for SOMA administration",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:   "bootstrap",
						Usage:  "Bootstrap authenticate to a new installation",
						Action: boottime(cmdOpsBootstrap),
					},
					// -> tree suspend/resume/start/stop
					// -> settings loglevel/opendoor/...
					// -> metrics?
					// -> termui?
				},
			},
		}...,
	)
	return &app
}

func cmdOpsBootstrap(c *cli.Context) error {
	// no command line arguments
	utl.ValidateCliArgumentCount(c, 0)
	var (
		err             error
		happy           bool
		password, token string
		kex, peer       *auth.Kex
		resp            *resty.Response
		tCred           *auth.Token
	)
	jBytes := &[]byte{}
	cipher := &[]byte{}
	plain := &[]byte{}

	fmt.Println(`
Welcome to SOMA!

This dialogue will guide you to set up the system's root account of
your new instance.

As first step, enter the root password you want to set.
`)

password_read:
	password = adm.ReadVerified(`password`)

	if happy, err = adm.EvaluatePassword(3, password, `root`, `soma`); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !happy {
		password = ""
		goto password_read
	}

	fmt.Println(`
Very good. Now enter the bootstrap token printed by somadbctl at the
end of the schema installation.
`)

	token = adm.ReadVerified(`password`)

	fmt.Println(`
Alright. Let's sully that pristine database. Here we go!  `)

	fmt.Printf("\nGenerating keypair: ")
	kex = auth.NewKex()
	fmt.Println(adm.GREEN+adm.SUCCESS+adm.CLEAR, ` OK`)

	fmt.Printf(`Initiating key exchange: `)
	if resp, err = Client.R().SetBody(kex).Post(`/authenticate/`); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}
	if resp.StatusCode() != 200 {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, resp.Status())
		os.Exit(1)
	}

	peer = &auth.Kex{}
	if err = json.Unmarshal(resp.Body(), peer); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}

	kex.SetPeerKey(peer.PublicKey())
	if err = kex.SetRequestUUID(peer.Request.String()); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}
	fmt.Println(adm.GREEN+adm.SUCCESS+adm.CLEAR, ` OK`)

	tCred = &auth.Token{
		UserName: `root`,
		Password: password,
		Token:    token,
	}
	fmt.Printf(`Sending bootstrap request: `)
	if *jBytes, err = json.Marshal(tCred); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}
	kex.SetTimeUTC()
	if err = kex.EncryptAndEncode(jBytes, cipher); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}

	if resp, err = Client.R().
		SetHeader(`Content-Type`, `application/octet-stream`).
		SetBody(*cipher).
		Put(fmt.Sprintf(
			"/authenticate/bootstrap/%s", kex.Request.String())); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}
	if resp.StatusCode() != 200 {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR+` FAILED`, resp.Status())
		os.Exit(1)
	}

	b := resp.Body()
	if err = kex.DecodeAndDecrypt(&b, plain); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}
	if err = json.Unmarshal(*plain, tCred); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}
	fmt.Println(adm.GREEN+adm.SUCCESS+adm.CLEAR, ` OK`)

	fmt.Printf(`Validating received token: `)
	if resp, err = Client.R().
		SetBasicAuth(`root`, tCred.Token).
		Get(`/authenticate/validate/`); err != nil || resp.StatusCode() != 204 {

		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR+` FAILED`, resp.String())
		os.Exit(1)
	}

	defer store.Close()
	fmt.Printf(`Writing token to local cache: `)
	if err = store.SaveToken(
		`root`,
		tCred.ValidFrom,
		tCred.ExpiresAt,
		tCred.Token,
	); err != nil {
		fmt.Fprintln(os.Stderr, adm.RED+adm.FAILURE+adm.CLEAR, err)
		os.Exit(1)
	}
	fmt.Println(adm.GREEN+adm.SUCCESS+adm.CLEAR, ` OK`)

	fmt.Println(`
All done. Thank you for flying with SOMA.
Suggested next steps:
	- create system_admin permission
	- create your team
	- create your user
	- grant system_admin to your user
	- activate your user
	- switch to using your user instead of root
	`)
	// TODO: output disclaimer text

	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
