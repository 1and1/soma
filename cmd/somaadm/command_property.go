package main

import (
	"fmt"
	"os"

	resty "gopkg.in/resty.v0"

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
	return adm.Perform(`postbody`, path, `command`, req, c)
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

	return adm.Perform(`postbody`, `/property/system/`, `command`, req, c)
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

	return adm.Perform(`postbody`, `/property/native/`, `command`, req, c)
}

func cmdPropertyServiceCreate(c *cli.Context) error {
	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	// fetch list of possible service attributes
	if resp, err = adm.GetReq(`/attributes/`); err != nil {
		return err
	}
	if err = adm.DecodedResponse(resp, res); err != nil {
		return err
	}
	if res.Attributes == nil || len(*res.Attributes) == 0 {
		return fmt.Errorf(`server returned no attributes for parsing`)
	}

	// sort attributes based on their cardinality so we can use them
	// for command line parsing
	multiple := []string{}
	unique := []string{}
	for _, attr := range *res.Attributes {
		switch attr.Cardinality {
		case "once":
			unique = append(unique, attr.Name)
		case "multi":
			multiple = append(multiple, attr.Name)
		default:
			return fmt.Errorf("Unknown attribute cardinality: %s",
				attr.Cardinality)
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
		return fmt.Errorf(
			"cmdPropertyServiceCreate called from unknown action %s",
			c.Command.Name)
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
	if err := adm.ValidateRuneCount(
		req.Property.Service.Name, 128); err != nil {
		return err
	}
	req.Property.Service.Attributes = make(
		[]proto.ServiceAttribute, 0, 16)
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
			req.Property.Service.Attributes = append(
				req.Property.Service.Attributes,
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
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// in main and not the cmpl lib because the full runtime is required
// to provide the completion options. This means we need access to
// globals that do not fit the function signature
func bashCompSvcCreate(c *cli.Context) {
	var (
		err  error
		resp *resty.Response
		res  *proto.Result
	)
	// fetch list of possible service attributes from SOMA
	if resp, err = adm.GetReq(`/attributes/`); err != nil {
		adm.Abort(err.Error())
	}
	if err = adm.DecodedResponse(resp, res); err != nil {
		adm.Abort(err.Error())
	}
	if res.Attributes == nil || len(*res.Attributes) == 0 {
		adm.Abort(`server returned no attributes for parsing`)
	}

	// sort attributes based on their cardinality so we can use them
	// for command line parsing
	multiple := []string{}
	unique := []string{}
	for _, attr := range *res.Attributes {
		switch attr.Cardinality {
		case "once":
			unique = append(unique, attr.Name)
		case "multi":
			multiple = append(multiple, attr.Name)
		default:
			adm.Abort(fmt.Sprintf("Unknown attribute cardinality: %s",
				attr.Cardinality))
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

	propId, err := adm.LookupCustomPropertyId(
		c.Args().First(), repoId)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/property/custom/%s/%s", repoId, propId)

	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdPropertySystemDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/property/system/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdPropertyNativeDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/property/native/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
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
	propId, err := adm.LookupServicePropertyId(
		c.Args().First(), teamId)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdPropertyTemplateDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	propId, err := adm.LookupTemplatePropertyId(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/property/service/global/%s", propId)
	return adm.Perform(`delete`, path, `command`, nil, c)
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
	propId, err := adm.LookupCustomPropertyId(c.Args().First(), repoId)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/property/custom/%s/%s", repoId,
		propId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdPropertySystemShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/property/system/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdPropertyNativeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	path := fmt.Sprintf("/property/native/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
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
	propId, err := adm.LookupServicePropertyId(
		c.Args().First(), teamId)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/property/service/team/%s/%s", teamId, propId)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdPropertyTemplateShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	propId, err := adm.LookupTemplatePropertyId(
		c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/property/service/global/%s", propId)
	return adm.Perform(`get`, path, `show`, nil, c)
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
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdPropertySystemList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/property/system/`, `list`, nil, c)
}

func cmdPropertyNativeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/property/native/`, `list`, nil, c)
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
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdPropertyTemplateList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/property/service/global/`, `list`, nil, c)
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
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
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
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	// deprecation warning
	switch oType {
	case `repository`, `bucket`, `node`:
		if _, ok := opts[`in`]; ok {
			fmt.Fprintf(
				os.Stderr,
				"Hint: Keyword `in` is DEPRECATED for %s objects,"+
					" since they are global objects. Ignoring.",
				oType,
			)
		}
	}

	var (
		objectId, object, repoId, bucketId string
		config                             *proto.NodeConfig
		req                                proto.Request
		err                                error
	)
	// id lookup
	switch oType {
	case `node`:
		if objectId, err = adm.LookupNodeId(opts[`to`][0]); err != nil {
			return err
		}
		if config, err = adm.LookupNodeConfig(objectId); err != nil {
			return err
		}
		repoId = config.RepositoryId
		bucketId = config.BucketId
	case `cluster`:
		bucketId, err = adm.LookupBucketId(opts["in"][0])
		if err != nil {
			return err
		}
		if objectId, err = adm.LookupClusterId(opts[`to`][0],
			bucketId); err != nil {
			return err
		}
		if repoId, err = adm.LookupRepoByBucket(bucketId); err != nil {
			return err
		}
	case `group`:
		bucketId, err = adm.LookupBucketId(opts["in"][0])
		if err != nil {
			return err
		}
		if objectId, err = adm.LookupGroupId(opts[`to`][0],
			bucketId); err != nil {
			return err
		}
		if repoId, err = adm.LookupRepoByBucket(bucketId); err != nil {
			return err
		}
	case `bucket`:
		bucketId, err = adm.LookupBucketId(opts["to"][0])
		if err != nil {
			return err
		}
		objectId = bucketId
		if repoId, err = adm.LookupRepoByBucket(bucketId); err != nil {
			return err
		}
	case `repository`:
		repoId, err = adm.LookupRepoId(opts[`to`][0])
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
		if err = adm.ValidateBool(opts[`childrenonly`][0],
			&prop.ChildrenOnly); err != nil {
			return err
		}
	} else {
		prop.ChildrenOnly = false
	}
	if _, ok := opts[`inheritance`]; ok {
		if err = adm.ValidateBool(opts[`inheritance`][0],
			&prop.Inheritance); err != nil {
			return err
		}
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
			if teamId, err = adm.LookupTeamByRepo(repoId); err != nil {
				return err
			}
		default:
			if teamId, err = adm.LookupTeamByBucket(
				bucketId); err != nil {
				return err
			}
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
		prop.Oncall.Name, prop.Oncall.Number, err = adm.LookupOncallDetails(
			oncallId,
		)
		if err != nil {
			return err
		}
	case `custom`:
		customId, err := adm.LookupCustomPropertyId(
			c.Args().First(), repoId)
		if err != nil {
			return err
		}

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
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
