package help

const CmdUserAdd string = ``

const CmdUserMarkDeleted string = ``

const CmdUserPurgeDeleted string = ``

const CmdUserUpdate string = ``

const CmdUserActivate string = ``

const CmdUserPasswordUpdate string = `
   This command can be used by a user to reset or change the password
   of the user's own account.

   It is an interactive command, prompting the user for required
   input.

   In the case of a password change, the current password must be
   provided. Otherwise either the LDAP password or mailtoken, similar
   to during activation.

SYNOPSIS:
   somaadm users password update [-r|--reset]

PERMISSIONS:
   This command requires no permissions (other than valid credentials).

EXAMPLES:
   ./somaadm -u myname users password update

   ./somaadm users password update --reset
`

const CmdUserList string = ``

const CmdUserShow string = ``

const CmdUserSync string = ``

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
