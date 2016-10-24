# somaadm ops bootstrap

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

# SYNOPSIS
```
somaadm ops bootstrap
```

# PERMISSIONS
This command runs before permissions are created. It
requires the DB installation issued root token and
grants access to the root account.
The root account has the omnipotence permission and
no restrictions.

# EXAMPLES
```
somaadm init
somaadm ops bootstrap
```
