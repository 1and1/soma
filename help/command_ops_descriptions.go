package help

const CmdOpsBootstrap string = `
   This command is used to bootstrap the root account of
   a new SOMA installation.
   
   It is an interactive command that will prompt you for
   the desired root user password, as well as the root
   token that was printed during installation of the
   database schema.

   After activating the root account and setting the
   password, a ticket will be granted which will be validated
   and stored in the local cache.

   For this reason, the local client files must have been
   initialized before.

SYNOPSIS:
   somaadm ops bootstrap

PERMISSIONS:
   This command runs before permissions are created. It
   requires the DB installation issued root token and
   grants access to the root account.
   The root account has the omnipotence permission and
   no restrictions.

EXAMPLES:
   ./somaadm init
   ./somaadm ops bootstrap`

const CmdOpsDumpToken string = `
   This commands dumps the currently active password token
   from the local client-cache.

   Note that this command will actually request a new token
   if there is none.

SYNOPSIS:
   somaadm ops dumptoken

PERMISSIONS:
   This command works on the local client cache and requires
   no permissions.`

const CmdOpsRepoStop string = `
   This command stops a specific repository on the server. The
   repository has to be either in state 'ready' or 'broken',
   ie. fully operational or faulted.
   
   A running repository that is still loading can not be stopped.

SYNOPSIS:
   somaadm ops repository stop ${repository}

ARGUMENT TYPES:
   repository   string   Name or UUID of the repository to stop

PERMISSIONS:
   This command requires one of the following permissions:
      - system_all

EXAMPLES:
   ./somaadm ops repository stop fump
   
   ./somaadm ops repository stop 21a4eda3-fafd-4c1b-9dd2-a29ef96ac916`

const CmdOpsRepoRebuild string = `
   The repository rebuild command rebuilds the dynamic objects
   inside a repository. All user created objects like groups,
   properties and check configurations are preserved.
   
   This is achieved by:
   - stopping the repository
   - marking all checks and check instances as deleted
   - restarting the repository in rebuild mode which persists the
     created checks and check instances into the database
   - restarting the repository in normal mode
   
   Rebuilds can run at two different levels:
   - instances, only rebuilds check instances
   - checks, rebuilds checks and check instances

SYNOPSIS:
   somaadm ops repository rebuild ${repository} \
      level ${level}

ARGUMENT TYPES:
   repository   string   Name or UUID of the repository to rebuild
   level        string   Level to perform the rebuild at

PERMISSIONS:
   This command requires one of the following permissions:
      - system_all

EXAMPLES:
   ./somaadm ops repository rebuild fump level checks
   
   ./somaadm ops repository rebuild \
      21a4eda3-fafd-4c1b-9dd2-a29ef96ac916 \
      level instances`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
