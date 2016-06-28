/*
 * Copyright (c) 2016, 1&1 Internet SE
 * Written by Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved.
 */

package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"gopkg.in/resty.v0"
)

func Failed(id string) {
	soma, _ := url.Parse(Eye.Soma.url.String())
	soma.Path = fmt.Sprintf("/deployments/id/%s/failed", id)
	client := resty.New().SetTimeout(500 * time.Millisecond)
	log.Printf("Sending fail feedback for %s\n", id)
	go client.R().Patch(soma.String())
}

func Success(id string) {
	soma, _ := url.Parse(Eye.Soma.url.String())
	soma.Path = fmt.Sprintf("/deployments/id/%s/success", id)
	client := resty.New().SetTimeout(500 * time.Millisecond)
	log.Printf("Sending success feedback for %s\n", id)
	go client.R().Patch(soma.String())
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
