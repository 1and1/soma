package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/1and1/soma/internal/adm"
	"github.com/1and1/soma/internal/cmpl"
	"github.com/1and1/soma/lib/proto"
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
								Name:         "service",
								Usage:        "Create a new per-team service property",
								Action:       runtime(cmdPropertyServiceCreate),
								BashComplete: comptime(bashCompSvcCreate),
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
								Name:         "custom",
								Usage:        "Create a new per-repo custom property",
								Action:       runtime(cmdPropertyCustomCreate),
								BashComplete: cmpl.Repository,
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
								Name:         "service",
								Usage:        "Delete a team service property",
								Action:       runtime(cmdPropertyServiceDelete),
								BashComplete: cmpl.Team,
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
								Name:         "custom",
								Usage:        "Delete a repository custom property",
								Action:       runtime(cmdPropertyCustomDelete),
								BashComplete: cmpl.Repository,
							},
							{
								Name:   "template",
								Usage:  "Delete a global service property template",
								Action: runtime(cmdPropertyTemplateDelete),
							},
						},
					}, // end property delete
					{
						Name:  "show",
						Usage: "SUBCOMMANDS for property show",
						Subcommands: []cli.Command{
							{
								Name:         "service",
								Usage:        "Show a service property",
								Action:       runtime(cmdPropertyServiceShow),
								BashComplete: cmpl.Team,
							},
							{
								Name:         "custom",
								Usage:        "Show a custom property",
								Action:       runtime(cmdPropertyCustomShow),
								BashComplete: cmpl.Repository,
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
								Name:         "service",
								Usage:        "List service properties",
								Action:       runtime(cmdPropertyServiceList),
								BashComplete: cmpl.Team,
							},
							{
								Name:         "custom",
								Usage:        "List custom properties",
								Action:       runtime(cmdPropertyCustomList),
								BashComplete: cmpl.Repository,
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
	multiple := []string{}
	unique := []string{"repository"}
	required := []string{"repository"}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}
	repoId, err := adm.LookupRepoId(opts[`repository`][0])
	if err != nil {
		return err
	}

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "custom"

	req.Property.Custom = &proto.PropertyCustom{}
	req.Property.Custom.Name = c.Args().First()
	req.Property.Custom.RepositoryId = repoId

	path := fmt.Sprintf("/property/custom/%s/", repoId)

	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertySystemCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "system"

	req.Property.System = &proto.PropertySystem{}
	req.Property.System.Name = c.Args().First()

	if resp, err := adm.PostReqBody(req, "/property/system/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyNativeCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Type = "native"

	req.Property.Native = &proto.PropertyNative{}
	req.Property.Native.Name = c.Args().First()

	if resp, err := adm.PostReqBody(req, "/property/native/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyServiceCreate(c *cli.Context) error {
	// fetch list of possible service attributes from SOMA
	attrResponse := utl.GetRequest(Client, "/attributes/")
	attrs := proto.Result{}
	err := json.Unmarshal(attrResponse.Body(), &attrs)
	if err != nil {
		adm.Abort("Failed to unmarshal Service Attribute data")
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
			adm.Abort()
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
		adm.Abort(
			fmt.Sprintf("cmdPropertyServiceCreate called from unknown action %s",
				c.Command.Name),
		)
	}

	// parse command line
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}
	// team lookup only for service
	var teamId string
	if c.Command.Name == "service" {
		var err error
		teamId, err = adm.LookupTeamId(opts["team"][0])
		if err != nil {
			return err
		}
	}

	// construct request body
	req := proto.Request{}
	req.Property = &proto.Property{}
	req.Property.Service = &proto.PropertyService{}
	req.Property.Service.Name = c.Args().First()
	if err := adm.ValidateRuneCount(req.Property.Service.Name, 128); err != nil {
		return err
	}
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
			if err := adm.ValidateRuneCount(oName, 128); err != nil {
				return err
			}
			if err := adm.ValidateRuneCount(oVal, 128); err != nil {
				return err
			}
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
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

// in main and not the cmpl lib because the full runtime is required
// to provide the completion options. This means we need access to
// globals that do not fit the function signature
func bashCompSvcCreate(c *cli.Context) {
	// fetch list of possible service attributes from SOMA
	attrResponse := utl.GetRequest(Client, "/attributes/")
	attrs := proto.Result{}
	err := json.Unmarshal(attrResponse.Body(), &attrs)
	if err != nil {
		adm.Abort("Failed to unmarshal Service Attribute data")
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
			adm.Abort()
		}
	}
	cmpl.GenericMulti(c, unique, multiple)
}

/* DELETE
 */
func cmdPropertyCustomDelete(c *cli.Context) error {
	multiple := []string{}
	unique := []string{"repository"}
	required := []string{"repository"}

	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail()); err != nil {
		return err
	}

	repoId, err := adm.LookupRepoId(opts[`repository`][0])
	if err != nil {
		return err
	}

	propId := utl.TryGetCustomPropertyByUUIDOrName(Client, c.Args().First(),
		repoId)
	path := fmt.Sprintf("/property/custom/%s/%s", repoId, propId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertySystemDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	path := fmt.Sprintf("/property/system/%s", c.Args().First())

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyNativeDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	path := fmt.Sprintf("/property/native/%s", c.Args().First())

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyServiceDelete(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`team`},
		[]string{`team`},
		c.Args().Tail()); err != nil {
		return err
	}
	teamId, err := adm.LookupTeamId(opts[`team`][0])
	if err != nil {
		return err
	}
	propId := utl.TryGetServicePropertyByUUIDOrName(Client,
		c.Args().First(), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyTemplateDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	propId := utl.TryGetTemplatePropertyByUUIDOrName(Client, c.Args().Get(0))
	path := fmt.Sprintf("/property/service/global/%s", propId)

	if resp, err := adm.DeleteReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

/* SHOW
 */
func cmdPropertyCustomShow(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`repository`},
		[]string{`repository`},
		c.Args().Tail()); err != nil {
		return err
	}
	repoId, err := adm.LookupRepoId(opts[`repository`][0])
	if err != nil {
		return err
	}
	propId := utl.TryGetCustomPropertyByUUIDOrName(Client, c.Args().First(),
		repoId)
	path := fmt.Sprintf("/property/custom/%s/%s", repoId,
		propId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertySystemShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	path := fmt.Sprintf("/property/system/%s", c.Args().First())
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyNativeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	path := fmt.Sprintf("/property/native/%s", c.Args().First())
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyServiceShow(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`team`},
		[]string{`team`},
		c.Args().Tail()); err != nil {
		return err
	}
	teamId, err := adm.LookupTeamId(opts[`team`][0])
	if err != nil {
		return err
	}
	propId := utl.TryGetServicePropertyByUUIDOrName(Client,
		c.Args().First(), teamId)
	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyTemplateShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	propId := utl.TryGetTemplatePropertyByUUIDOrName(Client,
		c.Args().First())
	path := fmt.Sprintf("/property/service/global/%s", propId)
	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

/* LIST
 */
func cmdPropertyCustomList(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`in`},
		[]string{`in`},
		adm.AllArguments(c)); err != nil {
		return err
	}
	repoId, err := adm.LookupRepoId(opts[`repository`][0])
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/property/custom/%s/", repoId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertySystemList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	if resp, err := adm.GetReq("/property/system/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyNativeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}
	if resp, err := adm.GetReq("/property/native/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyServiceList(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`team`},
		[]string{`team`},
		adm.AllArguments(c)); err != nil {
		return err
	}
	teamId, err := adm.LookupTeamId(opts[`team`][0])
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/property/service/team/%s/", teamId)

	if resp, err := adm.GetReq(path); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

func cmdPropertyTemplateList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	if resp, err := adm.GetReq("/property/service/global/"); err != nil {
		return err
	} else {
		fmt.Println(resp)
	}
	return nil
}

/* ADD
 */
func cmdPropertyAdd(c *cli.Context, pType, oType string) error {
	switch oType {
	case `node`, `bucket`, `repository`, `group`, `cluster`:
		switch pType {
		case `system`, `custom`, `service`, `oncall`:
		default:
			return fmt.Errorf("Unknown property type: %s", pType)
		}
	default:
		return fmt.Errorf("Unknown object type: %s", oType)
	}

	// argument parsing
	multiple := []string{}
	required := []string{`to`, `view`}
	unique := []string{`to`, `in`, `view`, `inheritance`, `childrenonly`}

	switch pType {
	case `system`:
		utl.CheckStringIsSystemProperty(Client, c.Args().First())
		fallthrough
	case `custom`:
		required = append(required, `value`)
		unique = append(unique, `value`)
	}
	switch oType {
	case `group`, `cluster`:
		required = append(required, `in`)
	}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(opts, multiple, unique, required,
		c.Args().Tail()); err != nil {
		return err
	}

	// deprecation warning
	switch oType {
	case `repository`, `bucket`, `node`:
		if _, ok := opts[`in`]; ok {
			fmt.Fprintf(
				os.Stderr,
				"Hint: Keyword `in` is DEPRECATED for %s objects, since they are global objects. Ignoring.",
				oType,
			)
		}
	}

	var (
		objectId, object, repoId, bucketId string
		config                             *proto.NodeConfig
		req                                proto.Request
	)
	// id lookup
	switch oType {
	case `node`:
		objectId = utl.TryGetNodeByUUIDOrName(Client, opts[`to`][0])
		config = utl.GetNodeConfigById(Client, objectId)
		repoId = config.RepositoryId
		bucketId = config.BucketId
	case `cluster`:
		bucketId, err := adm.LookupBucketId(opts["in"][0])
		if err != nil {
			return err
		}
		objectId = utl.TryGetClusterByUUIDOrName(Client, opts[`to`][0],
			bucketId)
		repoId = utl.GetRepositoryIdForBucket(Client, bucketId)
	case `group`:
		bucketId, err := adm.LookupBucketId(opts["in"][0])
		if err != nil {
			return err
		}
		objectId = utl.TryGetGroupByUUIDOrName(Client, opts[`to`][0],
			bucketId)
		repoId = utl.GetRepositoryIdForBucket(Client, bucketId)
	case `bucket`:
		bucketId, err := adm.LookupBucketId(opts["to"][0])
		if err != nil {
			return err
		}
		objectId = bucketId
		repoId = utl.GetRepositoryIdForBucket(Client, bucketId)
	case `repository`:
		repoId, err := adm.LookupRepoId(opts[`to`][0])
		if err != nil {
			return err
		}
		objectId = repoId
	}

	// property assembly
	prop := proto.Property{
		Type: pType,
		View: opts[`view`][0],
	}
	// property assembly, optional arguments
	if _, ok := opts[`childrenonly`]; ok {
		prop.ChildrenOnly = utl.GetValidatedBool(opts[`childrenonly`][0])
	} else {
		prop.ChildrenOnly = false
	}
	if _, ok := opts[`inheritance`]; ok {
		prop.Inheritance = utl.GetValidatedBool(opts[`inheritance`][0])
	} else {
		prop.Inheritance = true
	}
	switch pType {
	case `system`:
		prop.System = &proto.PropertySystem{
			Name:  c.Args().First(),
			Value: opts[`value`][0],
		}
	case `service`:
		var teamId string
		switch oType {
		case `repository`:
			teamId = utl.GetTeamIdByRepositoryId(Client, repoId)
		default:
			teamId = utl.TeamIdForBucket(Client, bucketId)
		}
		// no reason to fill out the attributes, client-provided
		// attributes are discarded by the server
		prop.Service = &proto.PropertyService{
			Name:       c.Args().First(),
			TeamId:     teamId,
			Attributes: []proto.ServiceAttribute{},
		}
	case `oncall`:
		oncallId, err := adm.LookupOncallId(c.Args().First())
		if err != nil {
			return err
		}
		prop.Oncall = &proto.PropertyOncall{
			Id: oncallId,
		}
		prop.Oncall.Name, prop.Oncall.Number = utl.GetOncallDetailsById(
			Client,
			oncallId,
		)
	case `custom`:
		customId := utl.TryGetCustomPropertyByUUIDOrName(
			Client, c.Args().First(), repoId)
		prop.Custom = &proto.PropertyCustom{
			Id:           customId,
			Name:         c.Args().First(),
			RepositoryId: repoId,
			Value:        opts[`value`][0],
		}
	}

	// request assembly
	switch oType {
	case `node`:
		req = proto.NewNodeRequest()
		req.Node.Id = objectId
		req.Node.Config = config
		req.Node.Properties = &[]proto.Property{prop}
	case `cluster`:
		req = proto.NewClusterRequest()
		req.Cluster.Id = objectId
		req.Cluster.BucketId = bucketId
		req.Cluster.Properties = &[]proto.Property{prop}
	case `group`:
		req = proto.NewGroupRequest()
		req.Group.Id = objectId
		req.Group.BucketId = bucketId
		req.Group.Properties = &[]proto.Property{prop}
	case `bucket`:
		req = proto.NewBucketRequest()
		req.Bucket.Id = objectId
		req.Bucket.Properties = &[]proto.Property{prop}
	case `repository`:
		req = proto.NewRepositoryRequest()
		req.Repository.Id = repoId
		req.Repository.Properties = &[]proto.Property{prop}
	}

	// request dispatch
	switch oType {
	case `repository`:
		object = oType
	default:
		object = oType + `s`
	}
	path := fmt.Sprintf("/%s/%s/property/%s/", object, objectId, pType)
	if resp, err := adm.PostReqBody(req, path); err != nil {
		return err
	} else {
		return adm.FormatOut(c, resp, ``)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
