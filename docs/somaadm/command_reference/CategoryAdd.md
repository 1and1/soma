# somaadm category add

This command is used to add a new permission category
to the system. Categories group permission sections
with the same scope, ie. if the actions inside the
section are for example global or per monitoringsystem.

The list of categories is defined by what is used by the
server's code. All categories have to be created.

# SYNOPSIS

```
somaadm category add ${category}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
category | string | Name of the category | | no

# PERMISSIONS

Since this command is used during bootstrap of the permission
system, it is generally advisable to have `omnipotence`.

# EXAMPLES

```
./somaadm category add global
./somaadm category add repository
./somaadm category add monitoring
```
