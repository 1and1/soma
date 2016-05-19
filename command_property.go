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
				Name:  "property",
				Usage: "SUBCOMMANDS for property",
				Subcommands: []cli.Command{
					{
						Name:  "create",
						Usage: "SUBCOMMANDS for property create",
						Subcommands: []cli.Command{
							{
								Name:   "service",
								Usage:  "Create a new per-team service property",
								Action: runtime(cmdPropertyServiceCreate),
							},
							{
								Name:   "system",
								Usage:  "Create a new global system property",
								Action: runtime(cmdPropertySystemCreate),
							},
							{
								Name:   "native",
								Usage:  "Create a new global native property",
								Action: runtime(cmdPropertyNativeCreate),
							},
							{
								Name:   "custom",
								Usage:  "Create a new per-repo custom property",
								Action: runtime(cmdPropertyCustomCreate),
							},
							{
								Name:   "template",
								Usage:  "Create a new global service template",
								Action: runtime(cmdPropertyServiceCreate),
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
								Action: runtime(cmdPropertyServiceDelete),
							},
							{
								Name:   "system",
								Usage:  "Delete a system property",
								Action: runtime(cmdPropertySystemDelete),
							},
							{
								Name:   "native",
								Usage:  "Delete a native property",
								Action: runtime(cmdPropertyNativeDelete),
							},
							{
								Name:   "custom",
								Usage:  "Delete a repository custom property",
								Action: runtime(cmdPropertyCustomDelete),
							},
							{
								Name:   "template",
								Usage:  "Delete a global service property template",
								Action: runtime(cmdPropertyTemplateDelete),
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
								Action: runtime(cmdPropertyServiceShow),
							},
							{
								Name:   "custom",
								Usage:  "Show a custom property",
								Action: runtime(cmdPropertyCustomShow),
							},
							{
								Name:   "system",
								Usage:  "Show a system property",
								Action: runtime(cmdPropertySystemShow),
							},
							{
								Name:   "native",
								Usage:  "Show a native property",
								Action: runtime(cmdPropertyNativeShow),
							},
							{
								Name:   "template",
								Usage:  "Show a service property template",
								Action: runtime(cmdPropertyTemplateShow),
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
								Action: runtime(cmdPropertyServiceList),
							},
							{
								Name:   "custom",
								Usage:  "List custom properties",
								Action: runtime(cmdPropertyCustomList),
							},
							{
								Name:   "system",
								Usage:  "List system properties",
								Action: runtime(cmdPropertySystemList),
							},
							{
								Name:   "native",
								Usage:  "List native properties",
								Action: runtime(cmdPropertyNativeList),
							},
							{
								Name:   "template",
								Usage:  "List service property templates",
								Action: runtime(cmdPropertyTemplateList),
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
func cmdPropertyCustomCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"repository"}
	required := []string{"repository"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts["repository"][0])

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "custom"

	req.Property.Custom = &proto.PropertyCustom{}
	req.Property.Custom.Name = c.Args().First()
	req.Property.Custom.RepositoryId = repoId

	path := fmt.Sprintf("/property/custom/%s/", repoId)

	resp := utl.PostRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertySystemCreate(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "system"

	req.Property.System = &proto.PropertySystem{}
	req.Property.System.Name = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, "/property/system/")
	fmt.Println(resp)
	return nil
}

func cmdPropertyNativeCreate(c *cli.Context) error {

	utl.ValidateCliArgumentCount(c, 1)

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "native"

	req.Property.Native = &proto.PropertyNative{}
	req.Property.Native.Name = c.Args().First()

	resp := utl.PostRequestWithBody(Client, req, "/property/native/")
	fmt.Println(resp)
	return nil
}

func cmdPropertyServiceCreate(c *cli.Context) error {
	utl.ValidateCliMinArgumentCount(c, 5)

	// fetch list of possible service attributes from SOMA
	attrResponse := utl.GetRequest(Client, "/attributes/")
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
		teamId = utl.TryGetTeamByUUIDOrName(Client, opts["team"][0])
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
	resp := utl.PostRequestWithBody(Client, req, path)
	fmt.Println(resp)
	return nil
}

/* DELETE
 */
func cmdPropertyCustomDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	multiple := []string{}
	unique := []string{"repository"}
	required := []string{"repository"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	repoId := utl.TryGetRepositoryByUUIDOrName(Client, opts["repository"][0])

	propId := utl.TryGetCustomPropertyByUUIDOrName(Client, c.Args().First(),
		repoId)
	path := fmt.Sprintf("/property/custom/%s/%s", repoId, propId)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertySystemDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/system/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertyNativeDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/native/%s", c.Args().First())

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertyServiceDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")
	teamId := utl.TryGetTeamByUUIDOrName(Client, c.Args().Get(2))
	propId := utl.TryGetServicePropertyByUUIDOrName(Client, c.Args().Get(0), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertyTemplateDelete(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetTemplatePropertyByUUIDOrName(Client, c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId)

	resp := utl.DeleteRequest(Client, path)
	fmt.Println(resp)
	return nil
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

func cmdPropertyCustomShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().Get(2))
	propId := utl.TryGetCustomPropertyByUUIDOrName(Client, c.Args().Get(0),
		repoId)
	path := fmt.Sprintf("/property/custom/%s/", repoId,
		propId)
	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertySystemShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/system/%s", c.Args().First())
	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertyNativeShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/native/%s", c.Args().First())
	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertyServiceShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")
	teamId := utl.TryGetTeamByUUIDOrName(Client, c.Args().Get(2))
	propId := utl.TryGetServicePropertyByUUIDOrName(Client, c.Args().Get(0), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertyTemplateShow(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetTemplatePropertyByUUIDOrName(Client, c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId)
	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertyCustomList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(Client, c.Args().Get(2))

	path := fmt.Sprintf("/property/custom/%s/", repoId)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertySystemList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest(Client, "/property/system/")
	fmt.Println(resp)
	return nil
}

func cmdPropertyNativeList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest(Client, "/property/native/")
	fmt.Println(resp)
	return nil
}

func cmdPropertyServiceList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 2)
	utl.ValidateCliArgument(c, 1, "team")
	teamId := utl.TryGetTeamByUUIDOrName(Client, c.Args().Get(1))

	path := fmt.Sprintf("/property/service/team/%s/", teamId)

	resp := utl.GetRequest(Client, path)
	fmt.Println(resp)
	return nil
}

func cmdPropertyTemplateList(c *cli.Context) error {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest(Client, "/property/service/global/")
	fmt.Println(resp)
	return nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
