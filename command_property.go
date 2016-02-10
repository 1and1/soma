package main

import (
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
	req := somaproto.ProtoRequestProperty{}
	req.Custom = &somaproto.ProtoPropertyCustom{}
	req.Custom.Property = c.Args().First()
	req.Custom.Repository = repoId
	path := fmt.Sprintf("/property/custom/%s/", repoId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdPropertySystemCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	req := somaproto.ProtoRequestProperty{}
	req.System = &somaproto.ProtoPropertySystem{}
	req.System.Property = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/property/system/")
	fmt.Println(resp)
}

func cmdPropertyNativeCreate(c *cli.Context) {
	utl.ValidateCliArgumentCount(c, 1)
	req := somaproto.ProtoRequestProperty{}
	req.Native = &somaproto.ProtoPropertyNative{}
	req.Native.Property = c.Args().First()

	resp := utl.PostRequestWithBody(req, "/property/native/")
	fmt.Println(resp)
}

func cmdPropertyServiceCreate(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{"transport", "application", "port", "file",
		"directory", "socket"}
	unique := []string{"tls", "provider", "team", "comm", "args",
		"uid"}
	required := []string{"team"}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	teamId := utl.TryGetTeamByUUIDOrName(opts["team"][0])

	req := somaproto.ProtoRequestProperty{}
	req.Service = &somaproto.ProtoPropertyService{}
	req.Service.Property = c.Args().First()
	req.Service.Team = teamId
	for key, arr := range opts {
		for _, val := range arr {
			switch key {
			case "transport":
				req.Service.Attributes.ProtoTransport = append(
					req.Service.Attributes.ProtoTransport,
					val,
				)
			case "application":
				req.Service.Attributes.ProtoApplication = append(
					req.Service.Attributes.ProtoApplication,
					val,
				)
			case "port":
				req.Service.Attributes.Port = append(
					req.Service.Attributes.Port,
					val,
				)
			case "comm":
				req.Service.Attributes.ProcessComm = append(
					req.Service.Attributes.ProcessComm,
					val,
				)
			case "args":
				req.Service.Attributes.ProcessArgs = append(
					req.Service.Attributes.ProcessArgs,
					val,
				)
			case "file":
				req.Service.Attributes.FilePath = append(
					req.Service.Attributes.FilePath,
					val,
				)
			case "directory":
				req.Service.Attributes.DirectoryPath = append(
					req.Service.Attributes.DirectoryPath,
					val,
				)
			case "socket":
				req.Service.Attributes.UnixSocketPath = append(
					req.Service.Attributes.UnixSocketPath,
					val,
				)
			case "uid":
				req.Service.Attributes.Uid = append(
					req.Service.Attributes.Uid,
					val,
				)
			case "tls":
				req.Service.Attributes.Tls = append(
					req.Service.Attributes.Tls,
					val,
				)
			case "provider":
				req.Service.Attributes.SoftwareProvider = append(
					req.Service.Attributes.SoftwareProvider,
					val,
				)
			default:
				utl.Abort("Error assigning service attributes")
			}
		}
	}
	path := fmt.Sprintf("/property/service/team/%s/", teamId)

	resp := utl.PostRequestWithBody(req, path)
	fmt.Println(resp)
}

func cmdPropertyTemplateCreate(c *cli.Context) {
	utl.ValidateCliMinArgumentCount(c, 5)
	multiple := []string{"transport", "application", "port", "file",
		"directory", "socket"}
	unique := []string{"tls", "provider", "team", "comm", "args",
		"uid"}
	required := []string{}

	opts := utl.ParseVariadicArguments(
		multiple,
		unique,
		required,
		c.Args().Tail())

	req := somaproto.ProtoRequestProperty{}
	req.Service = &somaproto.ProtoPropertyService{}
	req.Service.Property = c.Args().First()
	for key, arr := range opts {
		for _, val := range arr {
			switch key {
			case "transport":
				req.Service.Attributes.ProtoTransport = append(
					req.Service.Attributes.ProtoTransport,
					val,
				)
			case "application":
				req.Service.Attributes.ProtoApplication = append(
					req.Service.Attributes.ProtoApplication,
					val,
				)
			case "port":
				req.Service.Attributes.Port = append(
					req.Service.Attributes.Port,
					val,
				)
			case "comm":
				req.Service.Attributes.ProcessComm = append(
					req.Service.Attributes.ProcessComm,
					val,
				)
			case "args":
				req.Service.Attributes.ProcessArgs = append(
					req.Service.Attributes.ProcessArgs,
					val,
				)
			case "file":
				req.Service.Attributes.FilePath = append(
					req.Service.Attributes.FilePath,
					val,
				)
			case "directory":
				req.Service.Attributes.DirectoryPath = append(
					req.Service.Attributes.DirectoryPath,
					val,
				)
			case "socket":
				req.Service.Attributes.UnixSocketPath = append(
					req.Service.Attributes.UnixSocketPath,
					val,
				)
			case "uid":
				req.Service.Attributes.Uid = append(
					req.Service.Attributes.Uid,
					val,
				)
			case "tls":
				req.Service.Attributes.Tls = append(
					req.Service.Attributes.Tls,
					val,
				)
			case "provider":
				req.Service.Attributes.SoftwareProvider = append(
					req.Service.Attributes.SoftwareProvider,
					val,
				)
			default:
				utl.Abort("Error assigning service attributes")
			}
		}
	}

	resp := utl.PostRequestWithBody(req, "/property/service/global/")
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
