# somaadm ops repository rebuild

The repository rebuild command rebuilds the dynamic objects
inside a repository. All user created objects like groups,
properties and check configurations are preserved.

This is achieved by:

- stopping the repository
- marking all checks and check instances as deleted
- restarting the repository in rebuild mode which persists the
  created checks and check instances during startup into the database
- restarting the repository in normal mode

Rebuilds can run at two different levels:

- `instances`, only rebuilds check instances
- `checks`, rebuilds checks and check instances

# SYNOPSIS

```
somaadm ops repository rebuild ${repository} \
   level ${level}

```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | -------- 
repository | string | Name or UUID of the repository to rebuild | | no
level | string | Level to perform the rebuild at | | no

# PERMISSIONS

This command requires one of the following permissions:

* system\_all

# EXAMPLES

```
./somaadm ops repository rebuild common level checks
```
