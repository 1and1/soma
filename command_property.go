package main

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
)

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

	req := somaproto.PropertyRequest{}
	req.PropertyType = "custom"

	req.Custom = &somaproto.TreePropertyCustom{}
	req.Custom.Name = c.Args().First()
	req.Custom.RepositoryId = repoId

	path := fmt.Sprintf("/property/custom/%s/", repoId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdPropertySystemCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := somaproto.PropertyRequest{}
	req.PropertyType = "system"

	req.System = &somaproto.TreePropertySystem{}
	req.System.Name = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/property/system/")
	fmt.Println(resp)
}

func cmdPropertyNativeCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)

	req := somaproto.PropertyRequest{}
	req.PropertyType = "native"

	req.Native = &somaproto.TreePropertyNative{}
	req.Native.Name = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/property/native/")
	fmt.Println(resp)
}

func cmdPropertyServiceCreate(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 5)

	// fetch list of possible service attributes from SOMA
	attrResponse := utl.GetRequest("/attributes/")
	attrs := somaproto.AttributeResult{}
	err := json.Unmarshal(attrResponse.Body(), &attrs)
	if err != nil {
		utl.Abort("Failed to unmarshal Service Attribute data")
	}

	// sort attributes based on their cardinality so we can use them
	// for command line parsing
	multiple := []string{}
	unique := []string{}
	for _, attr := range attrs.Attributes {
		switch attr.Cardinality {
		case "once":
			unique = append(unique, attr.Attribute)
		case "multi":
			multiple = append(multiple, attr.Attribute)
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
	req := somaproto.PropertyRequest{}
	req.Service = &somaproto.TreePropertyService{}
	req.Service.Name = c.Args().First()
	req.Service.Attributes = make([]somaproto.TreeServiceAttribute, 0, 16)
	if c.Command.Name == "service" {
		req.PropertyType = `service`
		req.Service.TeamId = teamId
	} else {
		req.PropertyType = `template`
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
			req.Service.Attributes = append(req.Service.Attributes,
				somaproto.TreeServiceAttribute{
					Attribute: oName,
					Value:     oVal,
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
}

func cmdPropertySystemDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/system/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdPropertyNativeDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/native/%s", c.Args().First())

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdPropertyServiceDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")
	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetServicePropertyByUUIDOrName(c.Args().Get(0), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
}

func cmdPropertyTemplateDelete(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetTemplatePropertyByUUIDOrName(c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId)

	resp := utl.DeleteRequest(path)
	fmt.Println(resp)
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

}

func cmdPropertySystemShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/system/%s", c.Args().First())
	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdPropertyNativeShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	path := fmt.Sprintf("/property/native/%s", c.Args().First())
	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdPropertyServiceShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "team")
	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(2))
	propId := utl.TryGetServicePropertyByUUIDOrName(c.Args().Get(0), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdPropertyTemplateShow(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	propId := utl.TryGetTemplatePropertyByUUIDOrName(c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId)
	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdPropertyCustomList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 3)
	utl.ValidateCliArgument(c, 2, "repository")
	repoId := utl.TryGetRepositoryByUUIDOrName(c.Args().Get(2))

	path := fmt.Sprintf("/property/custom/%s/", repoId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdPropertySystemList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest("/property/system/")
	fmt.Println(resp)
}

func cmdPropertyNativeList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)
	resp := utl.GetRequest("/property/native/")
	fmt.Println(resp)
}

func cmdPropertyServiceList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 2)
	utl.ValidateCliArgument(c, 1, "team")
	teamId := utl.TryGetTeamByUUIDOrName(c.Args().Get(1))

	path := fmt.Sprintf("/property/service/team/%s/", teamId)

	resp := utl.GetRequest(path)
	fmt.Println(resp)
}

func cmdPropertyTemplateList(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 0)

	resp := utl.GetRequest("/property/service/global/")
	fmt.Println(resp)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
