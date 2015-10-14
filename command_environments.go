package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
	"log"
	"net/url"
)

func cmdEnvironmentsAdd(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/environments"

	a := c.Args()
	environment := a.First()
	if environment == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: add environment [%s]", environment)

	var req somaproto.ProtoRequestEnvironment
	req.Environment = environment

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Post(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
}

func cmdEnvironmentsRemove(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	environment := a.First()
	if environment == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: remove environment [%s]", environment)
	url.Path = fmt.Sprintf("/environments/%s", environment)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
}

func cmdEnvironmentsRename(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	// we expected exactly 3 arguments
	if len(a) != 3 {
		log.Fatal("Syntax error")
	}
	// second arg must be `to`
	if a.Get(1) != "to" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: rename environment [%s] to [%s]", a.Get(0), a.Get(2))

	var req somaproto.ProtoRequestEnvironment
	req.Environment = a.Get(2)
	url.Path = fmt.Sprintf("/environments/%s", a.Get(0))

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		SetBody(req).
		Put(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
}

func cmdEnvironmentsList(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/environments"

	a := c.Args()
	if len(a) != 0 {
		log.Fatal("Syntax error")
	}

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
}

func cmdEnvironmentsShow(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	if len(a) != 1 {
		log.Fatal("Syntax error")
	}
	environment := a.First()
	if environment == "" {
		log.Fatal("Syntax error")
	}
	url.Path = fmt.Sprintf("/environments/%s", environment)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
