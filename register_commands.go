package main

import (
	"github.com/codegangsta/cli"
)

func registerCommands(app cli.App) *cli.App {

	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "Initialize local client files",
			Action: cmdClientInit,
		}, // end init
		// views
		{
			Name:   "views",
			Usage:  "SUBCOMMANDS for views",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Register a new view",
					Action: cmdViewsAdd,
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing view",
					Action: cmdViewsRemove,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing view",
					Action: cmdViewsRename,
				},
				{
					Name:   "list",
					Usage:  "List all registered views",
					Action: cmdViewsList,
				},
				{
					Name:   "show",
					Usage:  "Show information about a specific view",
					Action: cmdViewsShow,
				},
			},
		}, // end views
		// environments
		{
			Name:   "environments",
			Usage:  "SUBCOMMANDS for environments",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Register a new view",
					Action: cmdEnvironmentsAdd,
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing unused environment",
					Action: cmdEnvironmentsRemove,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing environment",
					Action: cmdEnvironmentsRename,
				},
				{
					Name:   "list",
					Usage:  "List all available environments",
					Action: cmdEnvironmentsList,
				},
				{
					Name:   "show",
					Usage:  "Show information about a specific environment",
					Action: cmdEnvironmentsShow,
				},
			},
		}, // end environments
		// types
		{
			Name:   "types",
			Usage:  "SUBCOMMANDS for object types",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Add a new object type",
					Action: cmdObjectTypesAdd,
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing object type",
					Action: cmdObjectTypesRemove,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing object type",
					Action: cmdObjectTypesRename,
				},
				{
					Name:   "list",
					Usage:  "List all object types",
					Action: cmdObjectTypesList,
				},
				{
					Name:   "show",
					Usage:  "Show information about a specific object type",
					Action: cmdObjectTypesShow,
				},
			},
		}, // end types
		// states
		{
			Name:   "states",
			Usage:  "SUBCOMMANDS for states",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Add a new object state",
					Action: cmdObjectStatesAdd,
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing object state",
					Action: cmdObjectStatesRemove,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing object state",
					Action: cmdObjectStatesRename,
				},
				{
					Name:   "list",
					Usage:  "List all object states",
					Action: cmdObjectStatesList,
				},
				{
					Name:   "show",
					Usage:  "Show information about an object states",
					Action: cmdObjectStatesShow,
				},
			},
		}, // end states
		// datacenters
		{
			Name:   "datacenters",
			Usage:  "SUBCOMMANDS for datacenters",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Register a new datacenter",
					Action: cmdDatacentersAdd,
				},
				{
					Name:   "remove",
					Usage:  "Remove an existing datacenter",
					Action: cmdDatacentersRemove,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing datacenter",
					Action: cmdDatacentersRename,
				},
				{
					Name:   "list",
					Usage:  "List all datacenters",
					Action: cmdDatacentersList,
				},
				{
					Name:   "show",
					Usage:  "Show information about a specific datacenter",
					Action: cmdDatacentersShow,
				},
				{
					Name:   "groupadd",
					Usage:  "Add a datacenter to a datacenter group",
					Action: cmdDatacentersAddToGroup,
				},
				{
					Name:   "groupdel",
					Usage:  "Remove a datacenter from a datacenter group",
					Action: cmdDatacentersRemoveFromGroup,
				},
				{
					Name:   "grouplist",
					Usage:  "List all datacenter groups",
					Action: cmdDatacentersListGroups,
				},
				{
					Name:   "groupshow",
					Usage:  "Show information about a datacenter group",
					Action: cmdDatacentersShowGroup,
				},
			},
		}, // end datacenters
		// servers
		{
			Name:   "servers",
			Usage:  "SUBCOMMANDS for servers",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:        "create",
					Usage:       "Create a new physical server",
					Description: help.CmdServerCreate,
					Action:      cmdServerCreate,
				},
				{
					Name:   "delete",
					Usage:  "Mark an existing physical server as deleted",
					Action: cmdServerMarkAsDeleted,
				},
				{
					Name:   "purge",
					Usage:  "Remove all unreferenced servers marked as deleted",
					Action: cmdServerPurgeDeleted,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "all, a",
							Usage: "Purge all deleted servers",
						},
					},
				},
				{
					Name:   "update",
					Usage:  "Full update of server attributes (replace, not merge)",
					Action: cmdServerUpdate,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing server",
					Action: cmdServerRename,
				},
				{
					Name:   "online",
					Usage:  "Set an existing server to online",
					Action: cmdServerOnline,
				},
				{
					Name:   "offline",
					Usage:  "Set an existing server to offline",
					Action: cmdServerOffline,
				},
				{
					Name:   "move",
					Usage:  "Change a server's registered location",
					Action: cmdServerMove,
				},
				{
					Name:   "list",
					Usage:  "List all servers, see full description for possible filters",
					Action: cmdServerList,
				},
				{
					Name:   "show",
					Usage:  "Show details about a specific server",
					Action: cmdServerShow,
				},
				{
					Name:   "sync",
					Usage:  "Request a data sync for a server",
					Action: cmdServerSyncRequest,
				},
			},
		}, // end servers
		// permissions
		{
			Name:   "permissions",
			Usage:  "SUBCOMMANDS for permissions",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:  "type",
					Usage: "SUBCOMMANDS for permission types",
					Subcommands: []cli.Command{
						{
							Name:   "add",
							Usage:  "Register a new permission type",
							Action: cmdPermissionTypeAdd,
						},
						{
							Name:   "remove",
							Usage:  "Remove an existing permission type",
							Action: cmdPermissionTypeDel,
						},
						{
							Name:   "rename",
							Usage:  "Rename an existing permission type",
							Action: cmdPermissionTypeRename,
						},
						{
							Name:   "list",
							Usage:  "List all permission types",
							Action: cmdPermissionTypeList,
						},
						{
							Name:   "show",
							Usage:  "Show details for a permission type",
							Action: cmdPermissionTypeShow,
						},
					}, // end permissions type
				},
				{
					Name:   "add",
					Usage:  "Register a new permission",
					Action: cmdPermissionAdd,
				},
				{
					Name:   "remove",
					Usage:  "Remove a permission",
					Action: cmdPermissionDel,
				},
				{
					Name:   "list",
					Usage:  "List all permissions",
					Action: cmdPermissionList,
				},
				{
					Name:  "show",
					Usage: "SUBCOMMANDS for permission show",
					Subcommands: []cli.Command{
						{
							Name:   "user",
							Usage:  "Show permissions of a user",
							Action: cmdPermissionShowUser,
						},
						{
							Name:   "team",
							Usage:  "Show permissions of a team",
							Action: cmdPermissionShowTeam,
						},
						{
							Name:   "tool",
							Usage:  "Show permissions of a tool account",
							Action: cmdPermissionShowTool,
						},
						{
							Name:   "permission",
							Usage:  "Show details about a permission",
							Action: cmdPermissionShowPermission,
						},
					},
				}, // end permissions show
				{
					Name:   "audit",
					Usage:  "Show all limited permissions associated with a repository",
					Action: cmdPermissionAudit,
				},
				{
					Name:  "grant",
					Usage: "SUBCOMMANDS for permission grant",
					Subcommands: []cli.Command{
						{
							Name:   "enable",
							Usage:  "Enable a useraccount to receive GRANT permissions",
							Action: cmdPermissionGrantEnable,
						},
						{
							Name:   "global",
							Usage:  "Grant a global permission",
							Action: cmdPermissionGrantGlobal,
						},
						{
							Name:   "limited",
							Usage:  "Grant a limited permission",
							Action: cmdPermissionGrantLimited,
						},
						{
							Name:   "system",
							Usage:  "Grant a system permission",
							Action: cmdPermissionGrantSystem,
						},
					},
				}, // end permissions grant
			},
		}, // end permissions
		// teams
		{
			Name:   "teams",
			Usage:  "SUBCOMMANDS for teams",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Register a new team",
					Action: cmdTeamAdd,
				},
				{
					Name:   "remove",
					Usage:  "Delete an existing team",
					Action: cmdTeamDel,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing team",
					Action: cmdTeamRename,
				},
				{
					Name:   "migrate",
					Usage:  "Migrate users between teams",
					Action: cmdTeamMigrate,
				},
				{
					Name:   "list",
					Usage:  "List all teams",
					Action: cmdTeamList,
				},
				{
					Name:   "show",
					Usage:  "Show information about a team",
					Action: cmdTeamShow,
				},
			},
		}, // end teams
		// oncall
		{
			Name:   "oncall",
			Usage:  "SUBCOMMANDS for oncall duty teams",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Register a new oncall duty team",
					Action: cmdOnCallAdd,
				},
				{
					Name:   "remove",
					Usage:  "Delete an existing oncall duty team",
					Action: cmdOnCallDel,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing oncall duty team",
					Action: cmdOnCallRename,
				},
				{
					Name:   "update",
					Usage:  "Update phone number of an existing oncall duty team",
					Action: cmdOnCallUpdate,
				},
				{
					Name:   "list",
					Usage:  "List all registered oncall duty teams",
					Action: cmdOnCallList,
				},
				{
					Name:   "show",
					Usage:  "Show information about a specific oncall duty team",
					Action: cmdOnCallShow,
				},
				{
					Name:  "member",
					Usage: "SUBCOMMANDS to manipulate oncall duty members",
					Subcommands: []cli.Command{
						{
							Name:   "add",
							Usage:  "Add a user to an oncall duty team",
							Action: cmdOnCallMemberAdd,
						},
						{
							Name:   "remove",
							Usage:  "Remove a member from an oncall duty team",
							Action: cmdOnCallMemberDel,
						},
						{
							Name:   "list",
							Usage:  "List the users of an oncall duty team",
							Action: cmdOnCallMemberList,
						},
					},
				},
			},
		}, // end oncall
		// users
		{
			Name:   "users",
			Usage:  "SUBCOMMANDS for users",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Create a new user",
					Action: cmdUserAdd,
				},
				{
					Name:   "delete",
					Usage:  "Mark a user as deleted",
					Action: cmdUserMarkDeleted,
				},
				{
					Name:   "purge",
					Usage:  "Purge a user marked as deleted",
					Action: cmdUserPurgeDeleted,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "all, a",
							Usage: "Purge all deleted users",
						},
					},
				},
				{
					Name:   "restore",
					Usage:  "Restore a user marked as deleted",
					Action: cmdUserRestoreDeleted,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "all, a",
							Usage: "Restore all deleted users",
						},
					},
				},
				{
					Name:   "update",
					Usage:  "Set/change user information",
					Action: cmdUserUpdate,
				},
				{
					Name:   "rename",
					Usage:  "Change a user's username",
					Action: cmdUserRename,
				},
				{
					Name:   "activate",
					Usage:  "Activate a deativated user",
					Action: cmdUserActivate,
				},
				{
					Name:   "deactivate",
					Usage:  "Deactivate a user account",
					Action: cmdUserDeactivate,
				},
				{
					Name:  "password",
					Usage: "SUBCOMMANDS for user passwords",
					Subcommands: []cli.Command{
						{
							Name:   "update",
							Usage:  "Update the password of one's own user account",
							Action: cmdUserPasswordUpdate,
						},
						{
							Name:   "reset",
							Usage:  "Trigger a password reset for a user",
							Action: cmdUserPasswordReset,
						},
						{
							Name:   "force",
							Usage:  "Forcefully set the password of a user",
							Action: cmdUserPasswordForce,
						},
					},
				}, // end users password
				{
					Name:   "list",
					Usage:  "List all registered users",
					Action: cmdUserList,
				},
				{
					Name:   "show",
					Usage:  "Show information about a specific user",
					Action: cmdUserShow,
				},
			},
		}, // end users
		// nodes
		{
			Name:   "nodes",
			Usage:  "SUBCOMMANDS for nodes",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Register a new node",
					Action: cmdNodeAdd,
				},
				{
					Name:   "delete",
					Usage:  "Mark a node as deleted",
					Action: cmdNodeDel,
				},
				{
					Name:   "purge",
					Usage:  "Purge a node marked as deleted",
					Action: cmdNodePurge,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "all, a",
							Usage: "Purge all deleted nodes",
						},
					},
				},
				{
					Name:   "restore",
					Usage:  "Restore a node marked as deleted",
					Action: cmdNodeRestore,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "all, a",
							Usage: "Restore all deleted nodes",
						},
					},
				},
				{
					Name:   "rename",
					Usage:  "Rename a node",
					Action: cmdNodeRename,
				},
				{
					Name:   "repossess",
					Usage:  "Repossess a node to a different team",
					Action: cmdNodeRepo,
				},
				{
					Name:   "relocate",
					Usage:  "Relocate a node to a different server",
					Action: cmdNodeMove,
				},
				{
					Name:   "online",
					Usage:  "Set a nodes to online",
					Action: cmdNodeOnline,
				},
				{
					Name:   "offline",
					Usage:  "Set a node to offline",
					Action: cmdNodeOffline,
				},
				{
					Name:   "assign",
					Usage:  "Assign a node to configuration bucket",
					Action: cmdNodeAssign,
				},
				{
					Name:   "list",
					Usage:  "List all nodes",
					Action: cmdNodeList,
				},
				{
					Name:   "show",
					Usage:  "Show details about a node",
					Action: cmdNodeShow,
				},
				{
					Name:  "property",
					Usage: "SUBCOMMANDS for node properties",
					Subcommands: []cli.Command{
						{
							Name:   "add",
							Usage:  "Assign a property to a node",
							Action: cmdNodePropertyAdd,
						},
						{
							Name:   "get",
							Usage:  "Get the value of a node's specific property",
							Action: cmdNodePropertyGet,
						},
						{
							Name:   "delete",
							Usage:  "Delete a property from a node",
							Action: cmdNodePropertyDel,
						},
						{
							Name:   "list",
							Usage:  "List a nodes' local properties",
							Action: cmdNodePropertyList,
							Flags: []cli.Flag{
								cli.BoolFlag{
									Name:  "all, a",
									Usage: "List a nodes full properties (incl. inherited)",
								},
							},
						},
						{
							Name:   "show",
							Usage:  "Show details about a nodes properties",
							Action: cmdNodePropertyShow,
						},
					},
				}, // end nodes property
			},
		}, // end nodes
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
							Name:   "custom",
							Usage:  "Create a new per-repo custom property",
							Action: cmdPropertyCustomCreate,
						},
						{
							Name:   "template",
							Usage:  "Create a new global service template",
							Action: cmdPropertyTemplateCreate,
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
				/* XXX NOT IMPLEMENTED YET
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
							Name:   "template",
							Usage:  "Show a service property template",
							Action: cmdPropertyTemplateShow,
						},
					},
				}, // end property show
				*/
				/* XXX NOT IMPLEMENTED YET
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
							Name:   "template",
							Usage:  "List service property templates",
							Action: cmdPropertyTemplateList,
						},
					},
				}, // end property list
				*/
			},
		}, // end property
		// repository
		{
			Name:   "repository",
			Usage:  "SUBCOMMANDS for repository",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Create a new repository",
					Action: cmdRepositoryCreate,
				},
				{
					Name:   "delete",
					Usage:  "Mark an existing repository as deleted",
					Action: cmdRepositoryDelete,
				},
				{
					Name:   "restore",
					Usage:  "Restore a repository marked as deleted",
					Action: cmdRepositoryRestore,
				},
				{
					Name:   "purge",
					Usage:  "Remove an unreferenced deleted repository",
					Action: cmdRepositoryPurge,
				},
				{
					Name:   "clear",
					Usage:  "Clear all check instances for this repository",
					Action: cmdRepositoryClear,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing repository",
					Action: cmdRepositoryRename,
				},
				{
					Name:   "repossess",
					Usage:  "Change the owner of a repository",
					Action: cmdRepositoryRepossess,
				},
				/*
					{
						Name:   "clone",
						Usage:  "Create a clone of an existing repository",
						Action: cmdRepositoryClone,
					},
				*/
				{
					Name:   "activate",
					Usage:  "Activate a cloned repository",
					Action: cmdRepositoryActivate,
				},
				/*
					{
						Name:   "wipe",
						Usage:  "Clear all repository contents",
						Action: cmdRepositoryWipe,
					},
				*/
				{
					Name:   "list",
					Usage:  "List all existing repositories",
					Action: cmdRepositoryList,
				},
				{
					Name:   "show",
					Usage:  "Show information about a specific repository",
					Action: cmdRepositoryShow,
				},
			},
		}, // end repository
		// buckets
		{
			Name:   "buckets",
			Usage:  "SUBCOMMANDS for buckets",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Create a new bucket inside a repository",
					Action: cmdBucketCreate,
				},
				{
					Name:   "delete",
					Usage:  "Mark an existing bucket as deleted",
					Action: cmdBucketDelete,
				},
				{
					Name:   "restore",
					Usage:  "Restore a bucket marked as deleted",
					Action: cmdBucketRestore,
				},
				{
					Name:   "purge",
					Usage:  "Remove a deleted bucket",
					Action: cmdBucketPurge,
				},
				{
					Name:   "freeze",
					Usage:  "Freeze a bucket",
					Action: cmdBucketFreeze,
				},
				{
					Name:   "thaw",
					Usage:  "Thaw a frozen bucket",
					Action: cmdBucketThaw,
				},
				{
					Name:   "rename",
					Usage:  "Rename an existing bucket",
					Action: cmdBucketRename,
				},
			},
		}, // end buckets
		// clusters
		{
			Name:   "clusters",
			Usage:  "SUBCOMMANDS for clusters",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Create a new cluster",
					Action: cmdClusterCreate,
				},
				{
					Name:   "delete",
					Usage:  "Delete a cluster",
					Action: cmdClusterDelete,
				},
				{
					Name:   "rename",
					Usage:  "Rename a cluster",
					Action: cmdClusterRename,
				},
				{
					Name:   "list",
					Usage:  "List all clusters",
					Action: cmdClusterList,
				},
				{
					Name:   "show",
					Usage:  "Show details about a cluster",
					Action: cmdClusterShow,
				},
				{
					Name:  "members",
					Usage: "SUBCOMMANDS for cluster members",
					Subcommands: []cli.Command{
						{
							Name:   "add",
							Usage:  "Add a node to a cluster",
							Action: cmdClusterMemberAdd,
						},
						{
							Name:   "delete",
							Usage:  "Delete a node from a cluster",
							Action: cmdClusterMemberDelete,
						},
						{
							Name:   "list",
							Usage:  "List members of a cluster",
							Action: cmdClusterMemberList,
						},
					},
				},
			},
		}, // end clusters
		// groups
		{
			Name:   "groups",
			Usage:  "SUBCOMMANDS for groups",
			Before: runtimePreCmd,
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "Create a new group",
					Action: cmdGroupCreate,
				},
				{
					Name:   "delete",
					Usage:  "Delete a group",
					Action: cmdGroupDelete,
				},
				{
					Name:   "rename",
					Usage:  "Rename a group",
					Action: cmdGroupRename,
				},
				{
					Name:   "list",
					Usage:  "List all groups",
					Action: cmdGroupList,
				},
				{
					Name:   "show",
					Usage:  "Show details about a group",
					Action: cmdGroupShow,
				},
				{
					Name:  "members",
					Usage: "SUBCOMMANDS for members",
					Subcommands: []cli.Command{
						{
							Name:  "add",
							Usage: "SUBCOMMANDS for members add",
							Subcommands: []cli.Command{
								{
									Name:   "group",
									Usage:  "Add a group to a group",
									Action: cmdGroupMemberAddGroup,
								},
								{
									Name:   "cluster",
									Usage:  "Add a cluster to a group",
									Action: cmdGroupMemberAddCluster,
								},
								{
									Name:   "node",
									Usage:  "Add a node to a group",
									Action: cmdGroupMemberAddNode,
								},
							},
						},
						{
							Name:  "delete",
							Usage: "SUBCOMMANDS for members delete",
							Subcommands: []cli.Command{
								{
									Name:   "group",
									Usage:  "Delete a group from a group",
									Action: cmdGroupMemberDeleteGroup,
								},
								{
									Name:   "cluster",
									Usage:  "Delete a cluster from a group",
									Action: cmdGroupMemberDeleteCluster,
								},
								{
									Name:   "node",
									Usage:  "Delete a node from a group",
									Action: cmdGroupMemberDeleteNode,
								},
							},
						},
						{
							Name:   "list",
							Usage:  "List all members of a group",
							Action: cmdGroupMemberList,
						},
					},
				},
			},
		}, // end groups
	}
	return &app
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
