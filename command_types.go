package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
	"log"
	"net/url"
)

func cmdObjectTypesAdd(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/objtypes"

	a := c.Args()
	objectType := a.First()
	if objectType == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: add objectType [%s]", objectType)

	var req somaproto.ProtoRequestObjectType
	req.Type = objectType

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

func cmdObjectTypesRemove(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	objectType := a.First()
	if objectType == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: remove objectType [%s]", objectType)
	url.Path = fmt.Sprintf("/objtypes/%s", objectType)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
}

func cmdObjectTypesRename(c *cli.Context) {
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
	log.Printf("Command: rename objectType [%s] to [%s]", a.Get(0), a.Get(2))

	var req somaproto.ProtoRequestObjectType
	req.Type = a.Get(2)
	url.Path = fmt.Sprintf("/objtypes/%s", a.Get(0))

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

func cmdObjectTypesList(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/objtypes"

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

func cmdObjectTypesShow(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	if len(a) != 1 {
		log.Fatal("Syntax error")
	}
	objectType := a.First()
	if objectType == "" {
		log.Fatal("Syntax error")
	}
	url.Path = fmt.Sprintf("/objtypes/%s", objectType)

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
