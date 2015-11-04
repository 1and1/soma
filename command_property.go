package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

// GLOBAL SERVICE PROPERTY TEMPLATES
func cmdPropertyTemplateCreate(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetTemplatePropertyByUUIDOrName(c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId.String())

	_ = utl.DeleteRequest(path)
}

func cmdPropertyTemplateEdit(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateShow(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateList(c *cli.Context) {
	utl.NotImplemented()
}

// GLOBAL SYSTEM PROPERTIES
func cmdPropertySystemCreate(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertySystemDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetSystemPropertyByUUIDOrName(c.Args().Get(0))
	path := fmt.Sprintf("/property/system/%s", propId.String())

	_ = utl.DeleteRequest(path)
}

func cmdPropertySystemRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertySystemShow(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertySystemList(c *cli.Context) {
	utl.NotImplemented()
}

// PER-TEAM SERVICE PROPERTIES
func cmdPropertyServiceCreate(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyServiceDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")
	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetServicePropertyByUUIDOrName(c.Args().Get(0), teamId.String())
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId.String(), propId.String())

	_ = utl.DeleteRequest(path)
}

func cmdPropertyServiceEdit(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyServiceRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyServiceList(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyServiceShow(c *cli.Context) {
	utl.NotImplemented()
}

// PER-REPO CUSTOM PROPERTIES
func cmdPropertyCustomCreate(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyCustomDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetCustomPropertyByUUIDOrName(c.Args().Get(0), repoId.String())
	path := fmt.Sprintf("/property/custom/%s/%s", repoId.String(), propId.String())

	_ = utl.DeleteRequest(path)
}

func cmdPropertyCustomRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyCustomShow(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyCustomList(c *cli.Context) {
	utl.NotImplemented()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
