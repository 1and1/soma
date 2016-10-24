# somaadm users password update

This command can be used by a user to reset or change the password
of the user's own account.

It is an interactive command, prompting the user for required
input.

In the case of a password change, the current password must be
provided. Otherwise either the LDAP password or mailtoken, similar
to during activation.

# SYNOPSIS

```
somaadm users password update [-r|--reset]
```

# ARGUMENT TYPES

# PERMISSIONS

This command requires no permissions (other than valid credentials).

# EXAMPLES

```
./somaadm -u ${USER} users password update

./somaadm users password update --reset
```
