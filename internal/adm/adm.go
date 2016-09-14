package adm

import (
	"github.com/1and1/soma/internal/db"

	"gopkg.in/resty.v0"
)

var (
	client  *resty.Client
	cache   *db.DB
	async   bool
	jobSave bool
)

func ConfigureClient(c *resty.Client) {
	client = c
}

func ConfigureCache(c *db.DB) {
	cache = c
}

func ActivateAsyncWait(b bool) {
	async = b
}

func AutomaticJobSave(b bool) {
	jobSave = b
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
