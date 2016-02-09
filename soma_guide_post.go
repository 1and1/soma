package main

import (
	"database/sql"

	"github.com/satori/go.uuid"

)

type guidePost struct {
	input chan treeRequest
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
