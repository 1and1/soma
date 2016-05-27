package main

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
)

func registerProperty(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// property
			{
				Name:   "property",
				Usage:  "SUBCOMMANDS for property",
				Before: runtimePreCmd,
				Subcommands: []cli.Command{
					{
						Name:  "create",
						Usage: "SUBCOMMANDS for property create",
						Subcommands: []cli.Command{
							{
								Name:   "service",
								Usage:  "Create a new per-team service property",
								Action: cmdPropertyServiceCreate,
							},
							{
								Name:   "system",
								Usage:  "Create a new global system property",
								Action: cmdPropertySystemCreate,
							},
							{
								Name:   "native",
								Usage:  "Create a new global native property",
								Action: cmdPropertyNativeCreate,
							},
							{
								Name:   "custom",
								Usage:  "Create a new per-repo custom property",
								Action: cmdPropertyCustomCreate,
							},
							{
								Name:   "template",
								Usage:  "Create a new global service template",
								Action: cmdPropertyServiceCreate,
							},
						},
					}, // end property create
					{
						Name:  "delete",
						Usage: "SUBCOMMANDS for property delete",
						Subcommands: []cli.Command{
							{
								Name:   "service",
								Usage:  "Delete a team service property",
								Action: cmdPropertyServiceDelete,
							},
							{
								Name:   "system",
								Usage:  "Delete a system property",
								Action: cmdPropertySystemDelete,
							},
							{
								Name:   "native",
								Usage:  "Delete a native property",
								Action: cmdPropertyNativeDelete,
							},
							{
								Name:   "custom",
								Usage:  "Delete a repository custom property",
								Action: cmdPropertyCustomDelete,
							},
							{
								Name:   "template",
								Usage:  "Delete a global service property template",
								Action: cmdPropertyTemplateDelete,
							},
						},
					}, // end property delete
					/* XXX NOT IMPLEMENTED YET
					{
						Name:  "edit",
						Usage: "SUBCOMMANDS for property edit",
						Subcommands: []cli.Command{
							{
								Name:   "service",
								Usage:  "Edit a service property",
								Action: cmdPropertyServiceEdit,
							},
							{
								Name:   "template",
								Usage:  "Edit a service property template",
								Action: cmdPropertyTemplateEdit,
							},
						},
					}, // end property edit
					*/
					/* XXX NOT IMPLEMENTED YET
					{
						Name:  "rename",
						Usage: "SUBCOMMANDS for property rename",
						Subcommands: []cli.Command{
							{
								Name:   "service",
								Usage:  "Rename a service property",
								Action: cmdPropertyServiceRename,
							},
							{
								Name:   "custom",
								Usage:  "Rename a custom property",
								Action: cmdPropertyCustomRename,
							},
							{
								Name:   "system",
								Usage:  "Rename a system property",
								Action: cmdPropertySystemRename,
							},
							{
								Name:   "template",
								Usage:  "Rename a service property template",
								Action: cmdPropertyTemplateRename,
							},
						},
					}, // end property rename
					*/
					{
						Name:  "show",
						Usage: "SUBCOMMANDS for property show",
						Subcommands: []cli.Command{
							{
								Name:   "service",
								Usage:  "Show a service property",
								Action: cmdPropertyServiceShow,
							},
							{
								Name:   "custom",
								Usage:  "Show a custom property",
								Action: cmdPropertyCustomShow,
							},
							{
								Name:   "system",
								Usage:  "Show a system property",
								Action: cmdPropertySystemShow,
							},
							{
								Name:   "native",
								Usage:  "Show a native property",
								Action: cmdPropertyNativeShow,
							},
							{
								Name:   "template",
								Usage:  "Show a service property template",
								Action: cmdPropertyTemplateShow,
							},
						},
					}, // end property show
					{
						Name:  "list",
						Usage: "SUBCOMMANDS for property list",
						Subcommands: []cli.Command{
							{
								Name:   "service",
								Usage:  "List service properties",
								Action: cmdPropertyServiceList,
							},
							{
								Name:   "custom",
								Usage:  "List custom properties",
								Action: cmdPropertyCustomList,
							},
							{
								Name:   "system",
								Usage:  "List system properties",
								Action: cmdPropertySystemList,
							},
							{
								Name:   "native",
								Usage:  "List native properties",
								Action: cmdPropertyNativeList,
							},
							{
								Name:   "template",
								Usage:  "List service property templates",
								Action: cmdPropertyTemplateList,
							},
						},
					}, // end property list
				},
			}, // end property
		}...,
	)
	return &app
}

/* CREATE
 */
func cmdPropertyCustomCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"repository"}
	required := []string{"repository"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())
	repoId := utl.TryGetRepositoryByUUIDOrName(opts["repository"][0])

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "custom"

	req.Property.Custom = &proto.PropertyCustom{}
	req.Property.Custom.Name = c.Args().First()
	req.Property.Custom.RepositoryId = repoId

	path := fmt.Sprintf("/property/custom/%s/", repoId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertySystemCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "system"

	req.Property.System = &proto.PropertySystem{}
	req.Property.System.Name = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/property/system/")
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyNativeCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "native"

	req.Property.Native = &proto.PropertyNative{}
	req.Property.Native.Name = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/property/native/")
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyServiceCreate(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 5)

	// fetch list of possible service attributes from SOMA
	attrResponse := utl.GetRequest("/attributes/")
	attrs := proto.Result{}
	err := json.Unmarshal(attrResponse.Body(), &attrs)
	if err != nil {
		utl.Abort("Failed to unmarshal Service Attribute data")
	}

	// sort attributes based on their cardinality so we can use them
	// for command line parsing
	multiple := []string{}
	unique := []string{}
	for _, attr := range *attrs.Attributes {
		switch attr.Cardinality {
		case "once":
			unique = append(unique, attr.Name)
		case "multi":
			multiple = append(multiple, attr.Name)
		default:
			utl.Abort()
		}
	}
	required := []string{}

	switch c.Command.Name {
	case "service":
		// services are per team; add this as required as well
		required = append(required, "team")
		unique = append(unique, "team")
	case "template":
	default:
		utl.Abort(
			fmt.Sprintf("cmdPropertyServiceCreate called from unknown action %s",
				c.Command.Name),
		)
	}

	// parse command line
	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())
	// team lookup only for service
	var teamId string
	if c.Command.Name == "service" {
		teamId = utl.TryGetTeamByUUIDOrName(opts["team"][0])
	}

	// construct request body
	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Service = &proto.PropertyService{}
	req.Property.Service.Name = c.Args().First()
	req.Property.Service.Attributes = make([]proto.ServiceAttribute, 0, 16)
	if c.Command.Name == "service" {
		req.Property.Type = `service`
		req.Property.Service.TeamId = teamId
	} else {
		req.Property.Type = `template`
	}

	// fill attributes into request body
attrConversionLoop:
	for oName, _ := range opts {
		// the team that registers this service is not a service
		// attribute
		if c.Command.Name == `service` && oName == `team` {
			continue attrConversionLoop
		}
		for _, oVal := range opts[oName] {
			req.Property.Service.Attributes = append(req.Property.Service.Attributes,
				proto.ServiceAttribute{
					Name:  oName,
					Value: oVal,
				},
			)
		}
	}

	// send request
	var path string
	switch c.Command.Name {
	case `service`:
		path = fmt.Sprintf("/property/service/team/%s/", teamId)
	case `template`:
		path = `/property/service/global/`
	}
	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

/* DELETE
 */
func cmdPropertyCustomDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"repository"}
	required := []string{"repository"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	repoId := utl.TryGetRepositoryByUUIDOrName(opts["repository"][0])

	propId := utl.TryGetCustomPropertyByUUIDOrName(c.Args().First(),
		repoId)
	path := fmt.Sprintf("/property/custom/%s/%s", repoId, propId)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertySystemDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/system/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyNativeDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/native/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyServiceDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")
	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetServicePropertyByUUIDOrName(c.Args().Get(0), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyTemplateDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetTemplatePropertyByUUIDOrName(c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

/*
func cmdPropertyServiceEdit(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateEdit(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyCustomRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertySystemRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyServiceRename(c *cli.Context) {
	utl.NotImplemented()
}

func cmdPropertyTemplateRename(c *cli.Context) {
	utl.NotImplemented()
}
*/

func cmdPropertyCustomShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetCustomPropertyByUUIDOrName(c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/property/custom/%s/", repoId,
		propId)
	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)

}

func cmdPropertySystemShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/system/%s", c.Args().First())
	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyNativeShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/native/%s", c.Args().First())
	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyServiceShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")
	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetServicePropertyByUUIDOrName(c.Args().Get(0), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyTemplateShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetTemplatePropertyByUUIDOrName(c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId)
	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyCustomList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))

	path := fmt.Sprintf("/property/custom/%s/", repoId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertySystemList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest("/property/system/")
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyNativeList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest("/property/native/")
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyServiceList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 2)
	utl.ValidateCliArgument(c, 1, "team")
	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(1))

	path := fmt.Sprintf("/property/service/team/%s/", teamId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

func cmdPropertyTemplateList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/property/service/global/")
	fmt.Println(resp)
	utl.AsyncWait(Cfg.AsyncWait, resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
