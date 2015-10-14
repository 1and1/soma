package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/resty.v0"
	"log"
	"net/url"
)

func cmdDatacentersAdd(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/datacenters"

	a := c.Args()
	datacenter := a.First()
	if datacenter == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: add datacenter [%s]", datacenter)

	var req somaproto.ProtoRequestDatacenter
	req.Datacenter = datacenter

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

func cmdDatacentersRemove(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	datacenter := a.First()
	if datacenter == "" {
		log.Fatal("Syntax error")
	}
	log.Printf("Command: remove datacenter [%s]", datacenter)
	url.Path = fmt.Sprintf("/datacenters/%s", datacenter)

	resp, err := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).
		R().
		Delete(url.String())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s\n", resp.Status())
}

func cmdDatacentersRename(c *cli.Context) {
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
	log.Printf("Command: rename datacenter [%s] to [%s]", a.Get(0), a.Get(2))

	var req somaproto.ProtoRequestDatacenter
	req.Datacenter = a.Get(2)
	url.Path = fmt.Sprintf("/datacenters/%s", a.Get(0))

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

func cmdDatacentersList(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}
	url.Path = "/datacenters"

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

func cmdDatacentersShow(c *cli.Context) {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		log.Fatal(err)
	}

	a := c.Args()
	if len(a) != 1 {
		log.Fatal("Syntax error")
	}
	datacenter := a.First()
	if datacenter == "" {
		log.Fatal("Syntax error")
	}
	url.Path = fmt.Sprintf("/datacenters/%s", datacenter)

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
